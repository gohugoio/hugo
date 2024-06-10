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
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// The NPM package.json "template" file.
	FilenamePackageHugoJSON = "package.hugo.json"
	// The NPM package file.
	FilenamePackageJSON = "package.json"

	FilenameHugoStatsJSON = "hugo_stats.json"
)

func IsGoTmplExt(ext string) bool {
	return ext == "gotmpl"
}

// Supported data file extensions for _content.* files.
func IsContentDataExt(ext string) bool {
	return IsGoTmplExt(ext)
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
	FolderJSConfig  = "_jsconfig" // Mounted below /assets with postcss.config.js etc.

	NameContentData = "_content"
)

var (
	JsConfigFolderMountPrefix = filepath.Join(ComponentFolderAssets, FolderJSConfig)

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
