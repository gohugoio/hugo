// Copyright 2016-present The Hugo Authors. All rights reserved.
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
	"path"
	"path/filepath"
	"strings"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
)

type gitInfo struct {
	contentDir string
	repo       *gitmap.GitRepo
}

func (g *gitInfo) forPage(p *Page) (*gitmap.GitInfo, bool) {
	if g == nil {
		return nil, false
	}
	name := path.Join(g.contentDir, filepath.ToSlash(p.Path()))
	return g.repo.Files[name], true
}

func newGitInfo(cfg config.Provider) (*gitInfo, error) {
	var (
		workingDir = cfg.GetString("workingDir")
		contentDir = cfg.GetString("contentDir")
	)

	gitRepo, err := gitmap.Map(workingDir, "")
	if err != nil {
		return nil, err
	}

	repoPath := filepath.FromSlash(gitRepo.TopLevelAbsPath)
	// The Hugo site may be placed in a sub folder in the Git repo,
	// one example being the Hugo docs.
	// We have to find the root folder to the Hugo site below the Git root.
	contentRoot := strings.TrimPrefix(workingDir, repoPath)
	contentRoot = strings.TrimPrefix(contentRoot, helpers.FilePathSeparator)
	contentDir = path.Join(filepath.ToSlash(contentRoot), contentDir)

	return &gitInfo{contentDir: contentDir, repo: gitRepo}, nil
}
