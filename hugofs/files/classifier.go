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

package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/spf13/afero"
)

var (
	// This should be the only list of valid extensions for content files.
	contentFileExtensions = []string{
		"html", "htm",
		"mdown", "markdown", "md",
		"asciidoc", "adoc", "ad",
		"rest", "rst",
		"mmark",
		"org",
		"pandoc", "pdc"}

	contentFileExtensionsSet map[string]bool

	htmlFileExtensions = []string{
		"html", "htm"}

	htmlFileExtensionsSet map[string]bool
)

func init() {
	contentFileExtensionsSet = make(map[string]bool)
	for _, ext := range contentFileExtensions {
		contentFileExtensionsSet[ext] = true
	}
	htmlFileExtensionsSet = make(map[string]bool)
	for _, ext := range htmlFileExtensions {
		htmlFileExtensionsSet[ext] = true
	}
}

func IsContentFile(filename string) bool {
	return contentFileExtensionsSet[strings.TrimPrefix(filepath.Ext(filename), ".")]
}

func IsHTMLFile(filename string) bool {
	return htmlFileExtensionsSet[strings.TrimPrefix(filepath.Ext(filename), ".")]
}

func IsContentExt(ext string) bool {
	return contentFileExtensionsSet[ext]
}

type ContentClass string

const (
	ContentClassLeaf    ContentClass = "leaf"
	ContentClassBranch  ContentClass = "branch"
	ContentClassFile    ContentClass = "zfile" // Sort below
	ContentClassContent ContentClass = "zcontent"
)

func (c ContentClass) IsBundle() bool {
	return c == ContentClassLeaf || c == ContentClassBranch
}

func ClassifyContentFile(filename string, open func() (afero.File, error)) ContentClass {
	if !IsContentFile(filename) {
		return ContentClassFile
	}

	if IsHTMLFile(filename) {
		// We need to look inside the file. If the first non-whitespace
		// character is a "<", then we treat it as a regular file.
		// Eearlier we created pages for these files, but that had all sorts
		// of troubles, and isn't what it says in the documentation.
		// See https://github.com/gohugoio/hugo/issues/7030
		if open == nil {
			panic(fmt.Sprintf("no file opener provided for %q", filename))
		}

		f, err := open()
		if err != nil {
			return ContentClassFile
		}
		ishtml := isHTMLContent(f)
		f.Close()
		if ishtml {
			return ContentClassFile
		}

	}

	if strings.HasPrefix(filename, "_index.") {
		return ContentClassBranch
	}

	if strings.HasPrefix(filename, "index.") {
		return ContentClassLeaf
	}

	return ContentClassContent
}

var htmlComment = []rune{'<', '!', '-', '-'}

func isHTMLContent(r io.Reader) bool {
	br := bufio.NewReader(r)
	i := 0
	for {
		c, _, err := br.ReadRune()
		if err != nil {
			break
		}

		if i > 0 {
			if i >= len(htmlComment) {
				return false
			}

			if c != htmlComment[i] {
				return true
			}

			i++
			continue
		}

		if !unicode.IsSpace(c) {
			if i == 0 && c != '<' {
				return false
			}
			i++
		}
	}
	return true
}

const (
	ComponentFolderArchetypes = "archetypes"
	ComponentFolderStatic     = "static"
	ComponentFolderLayouts    = "layouts"
	ComponentFolderContent    = "content"
	ComponentFolderData       = "data"
	ComponentFolderAssets     = "assets"
	ComponentFolderI18n       = "i18n"

	FolderResources = "resources"
)

var (
	ComponentFolders = []string{
		ComponentFolderArchetypes,
		ComponentFolderStatic,
		ComponentFolderLayouts,
		ComponentFolderContent,
		ComponentFolderData,
		ComponentFolderAssets,
		ComponentFolderI18n,
	}

	componentFoldersSet = make(map[string]bool)
)

func init() {
	sort.Strings(ComponentFolders)
	for _, f := range ComponentFolders {
		componentFoldersSet[f] = true
	}
}

// ResolveComponentFolder returns "content" from "content/blog/foo.md" etc.
func ResolveComponentFolder(filename string) string {
	filename = strings.TrimPrefix(filename, string(os.PathSeparator))
	for _, cf := range ComponentFolders {
		if strings.HasPrefix(filename, cf) {
			return cf
		}
	}

	return ""
}

func IsComponentFolder(name string) bool {
	return componentFoldersSet[name]
}
