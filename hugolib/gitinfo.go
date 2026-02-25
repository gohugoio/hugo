// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/source"
)

type gitInfo struct {
	contentDir string
	repo       *gitmap.GitRepo

	// Per-module git repos keyed by module dir.
	moduleRepoLoaders map[string]func() (*gitmap.GitRepo, error)

	logger loggers.Logger
}

func (g *gitInfo) forPage(p page.Page) *source.GitInfo {
	if p.File() == nil {
		return nil
	}

	mod := p.File().FileInfo().Meta().Module
	if mod != nil && mod.IsGoMod() {
		mm := mod.(modules.Module)
		if loadRepo, ok := g.moduleRepoLoaders[mm.Dir()]; ok {
			repo, err := loadRepo()
			if err != nil {
				g.logger.Warnf("Failed to load git repo for module %s: %v", mm.Path(), err)
				return nil
			}
			return g.gitInfoFromModuleRepo(p, mm, repo)
		}
	}

	// Fall back to local git info.
	if g.repo == nil {
		return nil
	}
	name := strings.TrimPrefix(filepath.ToSlash(p.File().Filename()), g.contentDir)
	name = strings.TrimPrefix(name, "/")
	gi, found := g.repo.Files[name]
	if !found {
		return nil
	}
	return gi
}

func (g *gitInfo) gitInfoFromModuleRepo(p page.Page, mod modules.Module, repo *gitmap.GitRepo) *source.GitInfo {
	filename := filepath.ToSlash(p.File().Filename())
	modDir := filepath.ToSlash(mod.Dir())
	filePath := strings.TrimPrefix(filename, modDir)
	filePath = strings.TrimPrefix(filePath, "/")

	gi, found := repo.Files[filePath]
	if !found {
		return nil
	}
	return gi
}

// isGitModule returns true if the module has a git VCS origin.
func isGitModule(m modules.Module) bool {
	origin := m.Origin()
	return origin.VCS == "git" && origin.URL != ""
}

type gitInfoConfig struct {
	Deps         *deps.Deps
	Modules      modules.Modules
	GitInfoCache *filecache.Cache
	Logger       loggers.Logger
}

func newGitInfo(cfg gitInfoConfig) (*gitInfo, error) {
	g := &gitInfo{
		moduleRepoLoaders: make(map[string]func() (*gitmap.GitRepo, error)),
		logger:            cfg.Logger,
	}

	var hasGitModules bool
	var projectIsGitModule bool

	for _, mod := range cfg.Modules {
		if !isGitModule(mod) {
			continue
		}
		hasGitModules = true
		if mod.Owner() == nil && !mod.Origin().IsZero() {
			// TODO(bep) I'm not sure if this will ever be true. Needs some investigation,
			// as I'm sure there will be cases where this could be useful.
			projectIsGitModule = true
		}
	}

	if hasGitModules {
		if err := g.loadModuleRepos(cfg); err != nil {
			return nil, err
		}
	}

	if projectIsGitModule {
		return g, nil
	}

	// Load local git repo info.
	gitRepo, err := mapLocalRepo(cfg)
	if err != nil {
		if hasGitModules {
			// We added GitInfo module support in Hugo v0.157.0 and need some real world experience,
			// but for now, don't fail if the local repo is not a Git repo, but there are Git modules.
			// I'm not sure we should even warn, but let's do that for now.
			cfg.Logger.Warnf("Failed to read local Git log: %v", err)
			return g, nil
		}
		return nil, err
	}

	g.contentDir = gitRepo.TopLevelAbsPath
	g.repo = gitRepo

	return g, nil
}

func (g *gitInfo) loadModuleRepos(cfg gitInfoConfig) error {
	for _, mod := range cfg.Modules {
		if !isGitModule(mod) {
			continue
		}

		loadGitInfo := sync.OnceValues(func() (*gitmap.GitRepo, error) {
			origin := mod.Origin()

			cloneDir, err := ensureClone(cfg, origin)
			if err != nil {
				return nil, fmt.Errorf("failed to clone %s: %v", origin.URL, err)
			}

			repo, err := mapModuleRepo(cfg, cloneDir, origin.Hash)
			if err != nil {
				return nil, fmt.Errorf("failed to map git repo for module %s: %v", mod.Path(), err)
			}
			return repo, nil
		})

		g.moduleRepoLoaders[mod.Dir()] = loadGitInfo
	}

	return nil
}

// ensureClone ensures a blobless clone of the module's origin repo exists in the cache.
func ensureClone(cfg gitInfoConfig, origin modules.ModuleOrigin) (string, error) {
	key := "repo_" + hashing.XxHashFromStringHexEncoded(origin.URL)
	info, err := cfg.GitInfoCache.GetOrCreateInfo(key, func(id string) error {
		cloneDir := cfg.GitInfoCache.AbsFilenameFromID(id)
		if err := os.MkdirAll(cloneDir, 0o777); err != nil {
			return err
		}

		var stderr bytes.Buffer
		args := []any{
			"clone",
			"--filter=blob:none",
			"--no-checkout",
			origin.URL,
			cloneDir,
			hexec.WithStdout(io.Discard),
			hexec.WithStderr(&stderr),
		}

		cfg.Logger.Infof("Cloning gitinfo for repo %s into cache", origin.URL)

		cmd, err := cfg.Deps.ExecHelper.New("git", args...)
		if err != nil {
			return fmt.Errorf("git clone: %w", err)
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone %s: %w: %s", origin.URL, err, stderr.String())
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return cfg.GitInfoCache.AbsFilenameFromID(info.Name), nil
}

func mapModuleRepo(cfg gitInfoConfig, repoDir, revision string) (*gitmap.GitRepo, error) {
	opts := gitmap.Options{
		Repository: repoDir,
		Revision:   revision,
		GetGitCommandFunc: func(stdout, stderr io.Writer, args ...string) (gitmap.Runner, error) {
			var argsv []any
			for _, arg := range args {
				argsv = append(argsv, arg)
			}
			argsv = append(
				argsv,
				hexec.WithStdout(stdout),
				hexec.WithStderr(stderr),
			)
			return cfg.Deps.ExecHelper.New("git", argsv...)
		},
	}
	return gitmap.Map(opts)
}

func mapLocalRepo(cfg gitInfoConfig) (*gitmap.GitRepo, error) {
	opts := gitmap.Options{
		Repository: cfg.Deps.Conf.BaseConfig().WorkingDir,
		GetGitCommandFunc: func(stdout, stderr io.Writer, args ...string) (gitmap.Runner, error) {
			var argsv []any
			for _, arg := range args {
				argsv = append(argsv, arg)
			}
			argsv = append(
				argsv,
				hexec.WithStdout(stdout),
				hexec.WithStderr(stderr),
			)
			return cfg.Deps.ExecHelper.New("git", argsv...)
		},
	}
	return gitmap.Map(opts)
}
