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
	"path/filepath"
	"strings"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/config"
)

type gitInfo struct {
	contentDir string
	repo       *gitmap.GitRepo
}

func (g *gitInfo) forPage(p *Page) (*gitmap.GitInfo, bool) {
	if g == nil {
		return nil, false
	}

	name := strings.TrimPrefix(filepath.ToSlash(p.Filename()), g.contentDir)
	name = strings.TrimPrefix(name, "/")

	return g.repo.Files[name], true
}

func newGitInfo(cfg config.Provider) (*gitInfo, error) {
	workingDir := cfg.GetString("workingDir")

	gitRepo, err := gitmap.Map(workingDir, "")
	if err != nil {
		return nil, err
	}

	return &gitInfo{contentDir: gitRepo.TopLevelAbsPath, repo: gitRepo}, nil
}
