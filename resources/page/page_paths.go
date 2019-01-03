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

package page

import (
	"path/filepath"

	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
)

// TargetPathDescriptor describes how a file path for a given resource
// should look like on the file system. The same descriptor is then later used to
// create both the permalinks and the relative links, paginator URLs etc.
//
// The big motivating behind this is to have only one source of truth for URLs,
// and by that also get rid of most of the fragile string parsing/encoding etc.
//
//
type TargetPathDescriptor struct {
	PathSpec *helpers.PathSpec

	Type output.Format
	Kind string

	Sections []string

	// For regular content pages this is either
	// 1) the Slug, if set,
	// 2) the file base name (TranslationBaseName).
	BaseName string

	// Source directory.
	Dir string

	// Language prefix, set if multilingual and if page should be placed in its
	// language subdir.
	LangPrefix string

	// Whether this is a multihost multilingual setup.
	IsMultihost bool

	// URL from front matter if set. Will override any Slug etc.
	URL string

	// Used to create paginator links.
	Addends string

	// The expanded permalink if defined for the section, ready to use.
	ExpandedPermalink string

	// Some types cannot have uglyURLs, even if globally enabled, RSS being one example.
	UglyURLs bool
}

func CreateTargetPath(d TargetPathDescriptor) string {
	if d.Type.Name == "" {
		panic("CreateTargetPath: missing type")
	}

	pagePath := helpers.FilePathSeparator

	// The top level index files, i.e. the home page etc., needs
	// the index base even when uglyURLs is enabled.
	needsBase := true

	isUgly := d.UglyURLs && !d.Type.NoUgly

	if d.ExpandedPermalink == "" && d.BaseName != "" && d.BaseName == d.Type.BaseName {
		isUgly = true
	}

	if d.Kind != KindPage && d.URL == "" && len(d.Sections) > 0 {
		if d.ExpandedPermalink != "" {
			pagePath = filepath.Join(pagePath, d.ExpandedPermalink)
		} else {
			pagePath = filepath.Join(d.Sections...)
		}
		needsBase = false
	}

	if d.Type.Path != "" {
		pagePath = filepath.Join(pagePath, d.Type.Path)
	}

	if d.Kind != KindHome && d.URL != "" {
		if d.IsMultihost && d.LangPrefix != "" && !strings.HasPrefix(d.URL, "/"+d.LangPrefix) {
			pagePath = filepath.Join(d.LangPrefix, pagePath, d.URL)
		} else {
			pagePath = filepath.Join(pagePath, d.URL)
		}

		if d.Addends != "" {
			pagePath = filepath.Join(pagePath, d.Addends)
		}

		if strings.HasSuffix(d.URL, "/") || !strings.Contains(d.URL, ".") {
			pagePath = filepath.Join(pagePath, d.Type.BaseName+d.Type.MediaType.FullSuffix())
		}

	} else if d.Kind == KindPage {

		if d.ExpandedPermalink != "" {
			pagePath = filepath.Join(pagePath, d.ExpandedPermalink)

		} else {
			if d.Dir != "" {
				pagePath = filepath.Join(pagePath, d.Dir)
			}
			if d.BaseName != "" {
				pagePath = filepath.Join(pagePath, d.BaseName)
			}
		}

		if d.Addends != "" {
			pagePath = filepath.Join(pagePath, d.Addends)
		}

		if isUgly {
			pagePath += d.Type.MediaType.FullSuffix()
		} else {
			pagePath = filepath.Join(pagePath, d.Type.BaseName+d.Type.MediaType.FullSuffix())
		}

		if d.LangPrefix != "" {
			pagePath = filepath.Join(d.LangPrefix, pagePath)
		}
	} else {
		if d.Addends != "" {
			pagePath = filepath.Join(pagePath, d.Addends)
		}

		needsBase = needsBase && d.Addends == ""

		// No permalink expansion etc. for node type pages (for now)
		base := ""

		if needsBase || !isUgly {
			base = helpers.FilePathSeparator + d.Type.BaseName
		}

		pagePath += base + d.Type.MediaType.FullSuffix()

		if d.LangPrefix != "" {
			pagePath = filepath.Join(d.LangPrefix, pagePath)
		}
	}

	pagePath = filepath.Join(helpers.FilePathSeparator, pagePath)

	// Note: MakePathSanitized will lower case the path if
	// disablePathToLower isn't set.
	return d.PathSpec.MakePathSanitized(pagePath)
}
