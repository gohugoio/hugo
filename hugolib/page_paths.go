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

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
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

// createTargetPathDescriptor adapts a Page and the given output.Format into
// a targetPathDescriptor. This descriptor can then be used to create paths
// and URLs for this Page.
func (p *Page) createTargetPathDescriptor(t output.Format) (targetPathDescriptor, error) {
	if p.targetPathDescriptorPrototype == nil {
		panic(fmt.Sprintf("Must run initTargetPathDescriptor() for page %q, kind %q", p.title, p.Kind))
	}
	d := *p.targetPathDescriptorPrototype
	d.Type = t
	return d, nil
}

func (p *Page) initTargetPathDescriptor() error {
	d := &targetPathDescriptor{
		PathSpec:    p.s.PathSpec,
		Kind:        p.Kind,
		Sections:    p.sections,
		UglyURLs:    p.s.Info.uglyURLs(p),
		Dir:         filepath.ToSlash(p.Dir()),
		URL:         p.frontMatterURL,
		IsMultihost: p.s.owner.IsMultihost(),
	}

	if p.Slug != "" {
		d.BaseName = p.Slug
	} else {
		d.BaseName = p.TranslationBaseName()
	}

	if p.shouldAddLanguagePrefix() {
		d.LangPrefix = p.Lang()
	}

	// Expand only KindPage and KindTaxonomy; don't expand other Kinds of Pages
	// like KindSection or KindTaxonomyTerm because they are "shallower" and
	// the permalink configuration values are likely to be redundant, e.g.
	// naively expanding /category/:slug/ would give /category/categories/ for
	// the "categories" KindTaxonomyTerm.
	if p.Kind == KindPage || p.Kind == KindTaxonomy {
		if override, ok := p.Site.Permalinks[p.Section()]; ok {
			opath, err := override.Expand(p)
			if err != nil {
				return err
			}

			opath, _ = url.QueryUnescape(opath)
			opath = filepath.FromSlash(opath)
			d.ExpandedPermalink = opath
		}
	}

	p.targetPathDescriptorPrototype = d
	return nil

}

func (p *Page) initURLs() error {
	if len(p.outputFormats) == 0 {
		p.outputFormats = p.s.outputFormats[p.Kind]
	}
	target := filepath.ToSlash(p.createRelativeTargetPath())
	rel := p.s.PathSpec.URLizeFilename(target)

	var err error
	f := p.outputFormats[0]
	p.permalink, err = p.s.permalinkForOutputFormat(rel, f)
	if err != nil {
		return err
	}

	p.relTargetPathBase = strings.TrimPrefix(strings.TrimSuffix(target, f.MediaType.FullSuffix()), "/")
	if prefix := p.s.GetLanguagePrefix(); prefix != "" {
		// Any language code in the path will be added later.
		p.relTargetPathBase = strings.TrimPrefix(p.relTargetPathBase, prefix+"/")
	}
	p.relPermalink = p.s.PathSpec.PrependBasePath(rel, false)
	p.layoutDescriptor = p.createLayoutDescriptor()
	return nil
}

func (p *Page) initPaths() error {
	if err := p.initTargetPathDescriptor(); err != nil {
		return err
	}
	if err := p.initURLs(); err != nil {
		return err
	}
	return nil
}

// createTargetPath creates the target filename for this Page for the given
// output.Format. Some additional URL parts can also be provided, the typical
// use case being pagination.
func (p *Page) createTargetPath(t output.Format, noLangPrefix bool, addends ...string) (string, error) {
	d, err := p.createTargetPathDescriptor(t)
	if err != nil {
		return "", nil
	}

	if noLangPrefix {
		d.LangPrefix = ""
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

func (p *Page) createRelativeTargetPath() string {

	if len(p.outputFormats) == 0 {
		if p.Kind == kindUnknown {
			panic(fmt.Sprintf("Page %q has unknown kind", p.title))
		}
		panic(fmt.Sprintf("Page %q missing output format(s)", p.title))
	}

	// Choose the main output format. In most cases, this will be HTML.
	f := p.outputFormats[0]

	return p.createRelativeTargetPathForOutputFormat(f)

}

func (p *Page) createRelativePermalinkForOutputFormat(f output.Format) string {
	return p.s.PathSpec.URLizeFilename(p.createRelativeTargetPathForOutputFormat(f))
}

func (p *Page) createRelativeTargetPathForOutputFormat(f output.Format) string {
	tp, err := p.createTargetPath(f, p.s.owner.IsMultihost())

	if err != nil {
		p.s.Log.ERROR.Printf("Failed to create permalink for page %q: %s", p.FullFilePath(), err)
		return ""
	}

	// For /index.json etc. we must  use the full path.
	if f.MediaType.FullSuffix() == ".html" && filepath.Base(tp) == "index.html" {
		tp = strings.TrimSuffix(tp, f.BaseFilename())
	}

	return tp
}
