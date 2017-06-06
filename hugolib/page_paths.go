// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"fmt"
	"path/filepath"

	"net/url"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/output"
)

// targetPathDescriptor describes how a file path for a given resource
// should look like on the file system. The same descriptor is then later used to
// create both the permalinks and the relative links, paginator URLs etc.
//
// The big motivating behind this is to have only one source of truth for URLs,
// and by that also get rid of most of the fragile string parsing/encoding etc.
//
// Page.createTargetPathDescriptor is the Page adapter.
//
type targetPathDescriptor struct {
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

	// Page.URLPath.URL. Will override any Slug etc. for regular pages.
	URL string

	// Used to create paginator links.
	Addends string

	// The expanded permalink if defined for the section, ready to use.
	ExpandedPermalink string

	// Some types cannot have uglyURLs, even if globally enabled, RSS being one example.
	UglyURLs bool
}

// createTargetPathDescriptor adapts a Page and the given output.Format into
// a targetPathDescriptor. This descriptor can then be used to create paths
// and URLs for this Page.
func (p *Page) createTargetPathDescriptor(t output.Format) (targetPathDescriptor, error) {
	if p.targetPathDescriptorPrototype == nil {
		panic(fmt.Sprintf("Must run initTargetPathDescriptor() for page %q, kind %q", p.Title, p.Kind))
	}
	d := *p.targetPathDescriptorPrototype
	d.Type = t
	return d, nil
}

func (p *Page) initTargetPathDescriptor() error {

	d := &targetPathDescriptor{
		PathSpec: p.s.PathSpec,
		Kind:     p.Kind,
		Sections: p.sections,
		UglyURLs: p.s.Info.uglyURLs,
		Dir:      filepath.ToSlash(p.Source.Dir()),
		URL:      p.URLPath.URL,
	}

	if p.Slug != "" {
		d.BaseName = p.Slug
	} else {
		d.BaseName = p.TranslationBaseName()
	}

	if p.shouldAddLanguagePrefix() {
		d.LangPrefix = p.Lang()
	}

	if override, ok := p.Site.Permalinks[p.Section()]; ok {
		opath, err := override.Expand(p)
		if err != nil {
			return err
		}

		opath, _ = url.QueryUnescape(opath)
		opath = filepath.FromSlash(opath)
		d.ExpandedPermalink = opath
	}

	p.targetPathDescriptorPrototype = d
	return nil

}

// createTargetPath creates the target filename for this Page for the given
// output.Format. Some additional URL parts can also be provided, the typical
// use case being pagination.
func (p *Page) createTargetPath(t output.Format, addends ...string) (string, error) {
	d, err := p.createTargetPathDescriptor(t)
	if err != nil {
		return "", nil
	}

	if len(addends) > 0 {
		d.Addends = filepath.Join(addends...)
	}

	return createTargetPath(d), nil
}

func createTargetPath(d targetPathDescriptor) string {

	pagePath := helpers.FilePathSeparator

	// The top level index files, i.e. the home page etc., needs
	// the index base even when uglyURLs is enabled.
	needsBase := true

	isUgly := d.UglyURLs && !d.Type.NoUgly

	// If the page output format's base name is the same as the page base name,
	// we treat it as an ugly path, i.e.
	// my-blog-post-1/index.md => my-blog-post-1/index.html
	// (given the default values for that content file, i.e. no slug set etc.).
	// This introduces the behaviour from < Hugo 0.20, see issue #3396.
	if d.BaseName != "" && d.BaseName == d.Type.BaseName {
		isUgly = true
	}

	if d.Kind != KindPage && len(d.Sections) > 0 {
		pagePath = filepath.Join(d.Sections...)
		needsBase = false
	}

	if d.Type.Path != "" {
		pagePath = filepath.Join(pagePath, d.Type.Path)
	}

	if d.Kind == KindPage {
		// Always use URL if it's specified
		if d.URL != "" {
			pagePath = filepath.Join(pagePath, d.URL)
			if strings.HasSuffix(d.URL, "/") || !strings.Contains(d.URL, ".") {
				pagePath = filepath.Join(pagePath, d.Type.BaseName+"."+d.Type.MediaType.Suffix)
			}
		} else {
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
				pagePath += "." + d.Type.MediaType.Suffix
			} else {
				pagePath = filepath.Join(pagePath, d.Type.BaseName+"."+d.Type.MediaType.Suffix)
			}

			if d.LangPrefix != "" {
				pagePath = filepath.Join(d.LangPrefix, pagePath)
			}
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

		pagePath += base + "." + d.Type.MediaType.Suffix

		if d.LangPrefix != "" {
			pagePath = filepath.Join(d.LangPrefix, pagePath)
		}
	}

	pagePath = filepath.Join(helpers.FilePathSeparator, pagePath)

	// Note: MakePathSanitized will lower case the path if
	// disablePathToLower isn't set.
	return d.PathSpec.MakePathSanitized(pagePath)
}

func (p *Page) createRelativePermalink() string {

	if len(p.outputFormats) == 0 {
		panic(fmt.Sprintf("Page %q missing output format(s)", p.Title))
	}

	// Choose the main output format. In most cases, this will be HTML.
	f := p.outputFormats[0]

	return p.createRelativePermalinkForOutputFormat(f)

}

func (p *Page) createRelativePermalinkForOutputFormat(f output.Format) string {
	tp, err := p.createTargetPath(f)

	if err != nil {
		p.s.Log.ERROR.Printf("Failed to create permalink for page %q: %s", p.FullFilePath(), err)
		return ""
	}
	// For /index.json etc. we must  use the full path.
	if strings.HasSuffix(f.BaseFilename(), "html") {
		tp = strings.TrimSuffix(tp, f.BaseFilename())
	}

	return p.s.PathSpec.URLizeFilename(tp)
}

func (p *Page) TargetPath() (outfile string) {
	// Delete in Hugo 0.22
	helpers.Deprecated("Page", "TargetPath", "This method does not make sanse any more.", true)
	return ""
}
