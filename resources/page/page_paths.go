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
	"path"
	"path/filepath"

	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
)

const slash = "/"

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

	// Typically a language prefix added to file paths.
	PrefixFilePath string

	// Typically a language prefix added to links.
	PrefixLink string

	// If in multihost mode etc., every link/path needs to be prefixed, even
	// if set in URL.
	ForcePrefix bool

	// URL from front matter if set. Will override any Slug etc.
	URL string

	// Used to create paginator links.
	Addends string

	// The expanded permalink if defined for the section, ready to use.
	ExpandedPermalink string

	// Some types cannot have uglyURLs, even if globally enabled, RSS being one example.
	UglyURLs bool
}

// TODO(bep) move this type.
type TargetPaths struct {

	// Where to store the file on disk relative to the publish dir. OS slashes.
	TargetFilename string

	// The directory to write sub-resources of the above.
	SubResourceBaseTarget string

	// The base for creating links to sub-resources of the above.
	SubResourceBaseLink string

	// The relative permalink to this resources. Unix slashes.
	Link string
}

func (p TargetPaths) RelPermalink(s *helpers.PathSpec) string {
	return s.PrependBasePath(p.Link, false)
}

func (p TargetPaths) PermalinkForOutputFormat(s *helpers.PathSpec, f output.Format) string {
	var baseURL string
	var err error
	if f.Protocol != "" {
		baseURL, err = s.BaseURL.WithProtocol(f.Protocol)
		if err != nil {
			return ""
		}
	} else {
		baseURL = s.BaseURL.String()
	}

	return s.PermalinkForBaseURL(p.Link, baseURL)
}

func isHtmlIndex(s string) bool {
	return strings.HasSuffix(s, "/index.html")
}

func CreateTargetPaths(d TargetPathDescriptor) (tp TargetPaths) {

	if d.Type.Name == "" {
		panic("CreateTargetPath: missing type")
	}

	// Normalize all file Windows paths to simplify what's next.
	if helpers.FilePathSeparator != slash {
		d.Dir = filepath.ToSlash(d.Dir)
		d.PrefixFilePath = filepath.ToSlash(d.PrefixFilePath)

	}

	if d.URL != "" && !strings.HasPrefix(d.URL, "/") {
		// Treat this as a context relative URL
		d.ForcePrefix = true
	}

	pagePath := slash

	var (
		pagePathDir string
		link        string
		linkDir     string
	)

	// The top level index files, i.e. the home page etc., needs
	// the index base even when uglyURLs is enabled.
	needsBase := true

	isUgly := d.UglyURLs && !d.Type.NoUgly
	baseNameSameAsType := d.BaseName != "" && d.BaseName == d.Type.BaseName

	if d.ExpandedPermalink == "" && baseNameSameAsType {
		isUgly = true
	}

	if d.Kind != KindPage && d.URL == "" && len(d.Sections) > 0 {
		if d.ExpandedPermalink != "" {
			pagePath = pjoin(pagePath, d.ExpandedPermalink)
		} else {
			pagePath = pjoin(d.Sections...)
		}
		needsBase = false
	}

	if d.Type.Path != "" {
		pagePath = pjoin(pagePath, d.Type.Path)
	}

	if d.Kind != KindHome && d.URL != "" {
		pagePath = pjoin(pagePath, d.URL)

		if d.Addends != "" {
			pagePath = pjoin(pagePath, d.Addends)
		}

		pagePathDir = pagePath
		link = pagePath
		hasDot := strings.Contains(d.URL, ".")
		hasSlash := strings.HasSuffix(d.URL, slash)

		if hasSlash || !hasDot {
			pagePath = pjoin(pagePath, d.Type.BaseName+d.Type.MediaType.FullSuffix())
		} else if hasDot {
			pagePathDir = path.Dir(pagePathDir)
		}

		if !isHtmlIndex(pagePath) {
			link = pagePath
		} else if !hasSlash {
			link += slash
		}

		linkDir = pagePathDir

		if d.ForcePrefix {

			// Prepend language prefix if not already set in URL
			if d.PrefixFilePath != "" && !strings.HasPrefix(d.URL, slash+d.PrefixFilePath) {
				pagePath = pjoin(d.PrefixFilePath, pagePath)
				pagePathDir = pjoin(d.PrefixFilePath, pagePathDir)
			}

			if d.PrefixLink != "" && !strings.HasPrefix(d.URL, slash+d.PrefixLink) {
				link = pjoin(d.PrefixLink, link)
				linkDir = pjoin(d.PrefixLink, linkDir)
			}
		}

	} else if d.Kind == KindPage {

		if d.ExpandedPermalink != "" {
			pagePath = pjoin(pagePath, d.ExpandedPermalink)

		} else {
			if d.Dir != "" {
				pagePath = pjoin(pagePath, d.Dir)
			}
			if d.BaseName != "" {
				pagePath = pjoin(pagePath, d.BaseName)
			}
		}

		if d.Addends != "" {
			pagePath = pjoin(pagePath, d.Addends)
		}

		link = pagePath

		// TODO(bep) this should not happen after the fix in https://github.com/gohugoio/hugo/issues/4870
		// but we may need some more testing before we can remove it.
		if baseNameSameAsType {
			link = strings.TrimSuffix(link, d.BaseName)
		}

		pagePathDir = link
		link = link + slash
		linkDir = pagePathDir

		if isUgly {
			pagePath = addSuffix(pagePath, d.Type.MediaType.FullSuffix())
		} else {
			pagePath = pjoin(pagePath, d.Type.BaseName+d.Type.MediaType.FullSuffix())
		}

		if !isHtmlIndex(pagePath) {
			link = pagePath
		}

		if d.PrefixFilePath != "" {
			pagePath = pjoin(d.PrefixFilePath, pagePath)
			pagePathDir = pjoin(d.PrefixFilePath, pagePathDir)
		}

		if d.PrefixLink != "" {
			link = pjoin(d.PrefixLink, link)
			linkDir = pjoin(d.PrefixLink, linkDir)
		}

	} else {
		if d.Addends != "" {
			pagePath = pjoin(pagePath, d.Addends)
		}

		needsBase = needsBase && d.Addends == ""

		// No permalink expansion etc. for node type pages (for now)
		base := ""

		if needsBase || !isUgly {
			base = d.Type.BaseName
		}

		pagePathDir = pagePath
		link = pagePath
		linkDir = pagePathDir

		if base != "" {
			pagePath = path.Join(pagePath, addSuffix(base, d.Type.MediaType.FullSuffix()))
		} else {
			pagePath = addSuffix(pagePath, d.Type.MediaType.FullSuffix())

		}

		if !isHtmlIndex(pagePath) {
			link = pagePath
		} else {
			link += slash
		}

		if d.PrefixFilePath != "" {
			pagePath = pjoin(d.PrefixFilePath, pagePath)
			pagePathDir = pjoin(d.PrefixFilePath, pagePathDir)
		}

		if d.PrefixLink != "" {
			link = pjoin(d.PrefixLink, link)
			linkDir = pjoin(d.PrefixLink, linkDir)
		}
	}

	pagePath = pjoin(slash, pagePath)
	pagePathDir = strings.TrimSuffix(path.Join(slash, pagePathDir), slash)

	hadSlash := strings.HasSuffix(link, slash)
	link = strings.Trim(link, slash)
	if hadSlash {
		link += slash
	}

	if !strings.HasPrefix(link, slash) {
		link = slash + link
	}

	linkDir = strings.TrimSuffix(path.Join(slash, linkDir), slash)

	// Note: MakePathSanitized will lower case the path if
	// disablePathToLower isn't set.
	pagePath = d.PathSpec.MakePathSanitized(pagePath)
	pagePathDir = d.PathSpec.MakePathSanitized(pagePathDir)
	link = d.PathSpec.MakePathSanitized(link)
	linkDir = d.PathSpec.MakePathSanitized(linkDir)

	tp.TargetFilename = filepath.FromSlash(pagePath)
	tp.SubResourceBaseTarget = filepath.FromSlash(pagePathDir)
	tp.SubResourceBaseLink = linkDir
	tp.Link = d.PathSpec.URLizeFilename(link)
	if tp.Link == "" {
		tp.Link = slash
	}

	return
}

func addSuffix(s, suffix string) string {
	return strings.Trim(s, slash) + suffix
}

// Like path.Join, but preserves one trailing slash if present.
func pjoin(elem ...string) string {
	hadSlash := strings.HasSuffix(elem[len(elem)-1], slash)
	joined := path.Join(elem...)
	if hadSlash && !strings.HasSuffix(joined, slash) {
		return joined + slash
	}
	return joined
}
