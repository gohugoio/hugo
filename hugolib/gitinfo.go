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
	"github.com/spf13/hugo/helpers"
)

func (h *HugoSites) assembleGitInfo() {
	if !h.Cfg.GetBool("enableGitInfo") {
		return
	}

	var (
		workingDir = h.Cfg.GetString("workingDir")
		contentDir = h.Cfg.GetString("contentDir")
	)

	gitRepo, err := gitmap.Map(workingDir, "")
	if err != nil {
		h.Log.ERROR.Printf("Got error reading Git log: %s", err)
		return
	}

	gitMap := gitRepo.Files
	repoPath := filepath.FromSlash(gitRepo.TopLevelAbsPath)

	// The Hugo site may be placed in a sub folder in the Git repo,
	// one example being the Hugo docs.
	// We have to find the root folder to the Hugo site below the Git root.
	contentRoot := strings.TrimPrefix(workingDir, repoPath)
	contentRoot = strings.TrimPrefix(contentRoot, helpers.FilePathSeparator)

	s := h.Sites[0]

	for _, p := range s.AllPages {
		if p.Path() == "" {
			// Home page etc. with no content file.
			continue
		}
		// Git normalizes file paths on this form:
		filename := path.Join(filepath.ToSlash(contentRoot), contentDir, filepath.ToSlash(p.Path()))
		g, ok := gitMap[filename]
		if !ok {
			h.Log.WARN.Printf("Failed to find GitInfo for %q", filename)
			return
		}

		p.GitInfo = g
		p.Lastmod = p.GitInfo.AuthorDate
	}

}
