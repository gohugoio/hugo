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

	// Module path -> GitInfo for modules with GitHub origins.
	moduleGitInfo map[string]*source.GitInfo

	// Module path -> Module for quick lookup.
	modulesByPath map[string]modules.Module
}

func (g *gitInfo) forPage(p page.Page) *source.GitInfo {
	if p.File() == nil {
		return nil
	}

	// Check if this file is from a module.
	modulePath := p.File().FileInfo().Meta().Module
	if modulePath != "" {
		// Look up the module.
		mod, ok := g.modulesByPath[modulePath]
		if ok {
			origin := mod.Origin()
			if isGitHubModule(origin) {
				// Return the cached GitInfo for this module.
				if gi, found := g.moduleGitInfo[modulePath]; found {
					return gi
				}
			}
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

type gitInfoConfig struct {
	Deps         *deps.Deps
	Modules      modules.Modules
	GitInfoCache *filecache.Cache
	Logger       loggers.Logger
}

func newGitInfo(cfg gitInfoConfig) (*gitInfo, error) {
	g := &gitInfo{
		moduleGitInfo: make(map[string]*source.GitInfo),
		modulesByPath: make(map[string]modules.Module),
	}

	// Build a map of module paths to modules.
	for _, mod := range cfg.Modules {
		g.modulesByPath[mod.Path()] = mod
	}

	// Create GitHub client for fetching commit info.
	ghClient := newGitHubClient(cfg.GitInfoCache)

	// Group GitHub modules by repo and collect needed hashes.
	type repoInfo struct {
		owner, repo string
		needed      map[string]bool
		modules     []modules.Module
	}
	repos := make(map[string]*repoInfo) // keyed by "owner/repo"

	for _, mod := range cfg.Modules {
		origin := mod.Origin()
		if !isGitHubModule(origin) {
			continue
		}
		owner, repo := parseGitHubURL(origin.URL)
		if owner == "" || repo == "" {
			continue
		}
		rk := repoKey(owner, repo)
		ri, ok := repos[rk]
		if !ok {
			ri = &repoInfo{owner: owner, repo: repo, needed: make(map[string]bool)}
			repos[rk] = ri
		}
		ri.needed[origin.Hash] = true
		ri.modules = append(ri.modules, mod)
	}

	// Fetch commit history per repo in batches.
	for _, ri := range repos {
		// Use the first hash as the starting point for pagination.
		var startSHA string
		for h := range ri.needed {
			startSHA = h
			break
		}
		if err := ghClient.fetchRepoCommits(ri.owner, ri.repo, startSHA, ri.needed); err != nil {
			cfg.Logger.Warnf("Failed to fetch GitHub commit history for %s/%s: %v", ri.owner, ri.repo, err)
			continue
		}

		for _, mod := range ri.modules {
			commit := ghClient.getCommit(ri.owner, ri.repo, mod.Origin().Hash)
			if commit == nil {
				cfg.Logger.Warnf("Commit %s not found for module %s", mod.Origin().Hash, mod.Path())
				continue
			}
			g.moduleGitInfo[mod.Path()] = commit.toGitInfo()
		}
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
		if len(g.moduleGitInfo) > 0 {
			cfg.Logger.Warnf("Failed to read local Git log: %v", err)
			return g, nil
		}
		return nil, err
	}

	g.contentDir = gitRepo.TopLevelAbsPath
	g.repo = gitRepo

	return g, nil
}
