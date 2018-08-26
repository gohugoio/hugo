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
	"github.com/gohugoio/hugo/config"
	"github.com/tsuyoshiwada/go-gitlog"
	"path/filepath"
	"strings"
	"time"
)

// RevNumber alias for `-n <number>`
type RevFile struct {
	FileName string
}

// Args ...
func (rev *RevFile) Args() []string {
	return []string{"-p", rev.FileName}
}

type gitInfo struct {
	contentDir string
	gitlog     gitlog.GitLog
}

type gitPageInfo struct {
	AbbreviatedHash string
	AuthorName      string
	AuthorEmail     string
	AuthorDate      time.Time
	Hash            string
	Subject         string
	Commits         []*gitlog.Commit
}

func (g *gitInfo) forPage(p *Page) (*gitPageInfo, bool) {
	if g == nil {
		return nil, false
	}

	name := strings.TrimPrefix(filepath.ToSlash(p.Filename()), g.contentDir)
	name = strings.TrimPrefix(name, "/")
	logs, err := g.gitlog.Log(&RevFile{FileName: name}, nil)
	if err != nil {
		return nil, false
	}
	if len(logs) == 0 {
		return nil, false
	}
	return &gitPageInfo{
		AbbreviatedHash: logs[0].Hash.Short,
		AuthorName:      logs[0].Author.Name,
		AuthorEmail:     logs[0].Author.Email,
		AuthorDate:      logs[0].Author.Date,
		Hash:            logs[0].Hash.Long,
		Subject:         logs[0].Subject,
		Commits:         logs,
	}, true
}

func newGitInfo(cfg config.Provider) (*gitInfo, error) {
	workingDir := cfg.GetString("workingDir")
	git := gitlog.New(&gitlog.Config{
		Path: workingDir,
	})
	return &gitInfo{contentDir: workingDir, gitlog: git}, nil
}
