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
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/source"
)

type gitInfo struct {
	contentDir string
	repo       *gitmap.GitRepo
}

func (g *gitInfo) forPage(p page.Page) source.GitInfo {
	name := strings.TrimPrefix(filepath.ToSlash(p.File().Filename()), g.contentDir)
	name = strings.TrimPrefix(name, "/")
	gi, found := g.repo.Files[name]
	if !found {
		return source.GitInfo{}
	}
	return source.NewGitInfo(*gi)
}

func newGitInfo(d *deps.Deps) (*gitInfo, error) {
	opts := gitmap.Options{
		Repository: d.Conf.BaseConfig().WorkingDir,
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
			return d.ExecHelper.New("git", argsv...)
		},
	}

	gitRepo, err := gitmap.Map(opts)
	if err != nil {
		return nil, err
	}

	return &gitInfo{contentDir: gitRepo.TopLevelAbsPath, repo: gitRepo}, nil
}
