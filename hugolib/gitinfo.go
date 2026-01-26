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
	"io"
	"path/filepath"
	"strings"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/cache/filecache"
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

	// GitHub client for fetching per-file commit info.
	ghClient *gitHubClient

	// Module path -> Module for quick lookup.
	modulesByPath map[string]modules.Module

	logger loggers.Logger
}

func (g *gitInfo) forPage(p page.Page) *source.GitInfo {
	if p.File() == nil {
		return nil
	}

	mod := p.File().FileInfo().Meta().Module
	if mod.IsGoMod() {
		mm := mod.(modules.Module)
		if isGitHubModule(mm) {
			return g.gitInfoForGitHubFile(p, mm)
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

func (g *gitInfo) gitInfoForGitHubFile(p page.Page, mod modules.Module) *source.GitInfo {
	if g.ghClient == nil {
		return nil
	}

	origin := mod.Origin()
	owner, repo := parseGitHubURL(origin.URL)
	if owner == "" || repo == "" {
		return nil
	}

	// Get the file path relative to the module root (repo-relative path).
	filename := filepath.ToSlash(p.File().Filename())
	modDir := filepath.ToSlash(mod.Dir())
	filePath := strings.TrimPrefix(filename, modDir)
	filePath = strings.TrimPrefix(filePath, "/")

	// Fetch the commit info for this specific file.
	// Use origin.Ref (e.g., "refs/tags/v3.0.1") for the GraphQL query.
	commit, err := g.ghClient.fetchFileCommit(owner, repo, origin.Ref, filePath)
	if err != nil {
		g.logger.Warnf("Failed to fetch GitHub commit for %s/%s:%s: %v", owner, repo, filePath, err)
		return nil
	}
	if commit == nil {
		return nil
	}

	return commit.toGitInfo()
}

type gitInfoConfig struct {
	Deps         *deps.Deps
	Modules      modules.Modules
	GitInfoCache *filecache.Cache
	Logger       loggers.Logger
}

func newGitInfo(cfg gitInfoConfig) (*gitInfo, error) {
	g := &gitInfo{
		modulesByPath: make(map[string]modules.Module),
		logger:        cfg.Logger,
	}

	// Build a map of module paths to modules.
	for _, mod := range cfg.Modules {
		g.modulesByPath[mod.Path()] = mod
	}

	var hasGitHubModules bool
	var projectIsGitHubModule bool

	for _, mod := range cfg.Modules {
		if isGitHubModule(mod) {
			hasGitHubModules = true
			projectIsGitHubModule = mod.Owner() == nil // Project is always the first, so no need to look further.
			break
		}
	}

	if hasGitHubModules {
		g.ghClient = newGitHubClient(cfg.GitInfoCache)
	}

	if projectIsGitHubModule {
		return g, nil
	}

	// Load local git repo info.
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

	gitRepo, err := gitmap.Map(opts)
	if err != nil {
		// Don't fail if local git repo is not available,
		// module-based git info may still work.
		if hasGitHubModules {
			cfg.Logger.Warnf("Failed to read local Git log: %v", err)
			return g, nil
		}
		return nil, err
	}

	g.contentDir = gitRepo.TopLevelAbsPath
	g.repo = gitRepo

	return g, nil
}
