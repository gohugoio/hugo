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
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/resources/page"
)

// PageCollections contains the page collections for a site.
type PageCollections struct {
	pageMap *pageMap

	// Lazy initialized page collections
	pages           *lazyPagesFactory
	regularPages    *lazyPagesFactory
	allPages        *lazyPagesFactory
	allRegularPages *lazyPagesFactory
}

// Pages returns all pages.
// This is for the current language only.
func (c *PageCollections) Pages() page.Pages {
	return c.pages.get()
}

// RegularPages returns all the regular pages.
// This is for the current language only.
func (c *PageCollections) RegularPages() page.Pages {
	return c.regularPages.get()
}

// AllPages returns all pages for all languages.
func (c *PageCollections) AllPages() page.Pages {
	return c.allPages.get()
}

// AllPages returns all regular pages for all languages.
func (c *PageCollections) AllRegularPages() page.Pages {
	return c.allRegularPages.get()
}

type lazyPagesFactory struct {
	pages page.Pages

	init    sync.Once
	factory page.PagesFactory
}

func (l *lazyPagesFactory) get() page.Pages {
	l.init.Do(func() {
		l.pages = l.factory()
	})
	return l.pages
}

func newLazyPagesFactory(factory page.PagesFactory) *lazyPagesFactory {
	return &lazyPagesFactory{factory: factory}
}

func newPageCollections(m *pageMap) *PageCollections {
	if m == nil {
		panic("must provide a pageMap")
	}

	c := &PageCollections{pageMap: m}

	c.pages = newLazyPagesFactory(func() page.Pages {
		return m.createListAllPages()
	})

	c.regularPages = newLazyPagesFactory(func() page.Pages {
		return c.findPagesByKindIn(page.KindPage, c.pages.get())
	})

	return c
}

// This is an adapter func for the old API with Kind as first argument.
// This is invoked when you do .Site.GetPage. We drop the Kind and fails
// if there are more than 2 arguments, which would be ambigous.
func (c *PageCollections) getPageOldVersion(ref ...string) (page.Page, error) {
	var refs []string
	for _, r := range ref {
		// A common construct in the wild is
		// .Site.GetPage "home" "" or
		// .Site.GetPage "home" "/"
		if r != "" && r != "/" {
			refs = append(refs, r)
		}
	}

	var key string

	if len(refs) > 2 {
		// This was allowed in Hugo <= 0.44, but we cannot support this with the
		// new API. This should be the most unusual case.
		return nil, fmt.Errorf(`too many arguments to .Site.GetPage: %v. Use lookups on the form {{ .Site.GetPage "/posts/mypage-md" }}`, ref)
	}

	if len(refs) == 0 || refs[0] == page.KindHome {
		key = "/"
	} else if len(refs) == 1 {
		if len(ref) == 2 && refs[0] == page.KindSection {
			// This is an old style reference to the "Home Page section".
			// Typically fetched via {{ .Site.GetPage "section" .Section }}
			// See https://github.com/gohugoio/hugo/issues/4989
			key = "/"
		} else {
			key = refs[0]
		}
	} else {
		key = refs[1]
	}

	key = filepath.ToSlash(key)
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	return c.getPageNew(nil, key)
}

// 	Only used in tests.
func (c *PageCollections) getPage(typ string, sections ...string) page.Page {
	refs := append([]string{typ}, path.Join(sections...))
	p, _ := c.getPageOldVersion(refs...)
	return p
}

// getPageRef resolves a Page from ref/relRef, with a slightly more comprehensive
// search path than getPageNew.
func (c *PageCollections) getPageRef(context page.Page, ref string) (page.Page, error) {
	n, err := c.getContentNode(context, true, ref)
	if err != nil || n == nil || n.p == nil {
		return nil, err
	}
	return n.p, nil
}

func (c *PageCollections) getPageNew(context page.Page, ref string) (page.Page, error) {
	n, err := c.getContentNode(context, false, ref)
	if err != nil || n == nil || n.p == nil {
		return nil, err
	}
	return n.p, nil
}

func (c *PageCollections) getSectionOrPage(ref string) (*contentNode, string) {
	var n *contentNode

	s, v, found := c.pageMap.sections.LongestPrefix(ref)

	if found {
		n = v.(*contentNode)
	}

	if found && s == ref {
		// A section
		return n, ""
	}

	m := c.pageMap
	filename := strings.TrimPrefix(strings.TrimPrefix(ref, s), "/")
	langSuffix := "." + m.s.Lang()

	// Trim both extension and any language code.
	name := helpers.PathNoExt(filename)
	name = strings.TrimSuffix(name, langSuffix)

	// These are reserved bundle names and will always be stored by their owning
	// folder name.
	name = strings.TrimSuffix(name, "/index")
	name = strings.TrimSuffix(name, "/_index")

	if !found {
		return nil, name
	}

	// Check if it's a section with filename provided.
	if !n.p.File().IsZero() && n.p.File().LogicalName() == filename {
		return n, name
	}

	return m.getPage(s, name), name

}

// For Ref/Reflink and .Site.GetPage do simple name lookups for the potentially ambigous myarticle.md and /myarticle.md,
// but not when we get ./myarticle*, section/myarticle.
func shouldDoSimpleLookup(ref string) bool {
	if ref[0] == '.' {
		return false
	}

	slashCount := strings.Count(ref, "/")

	if slashCount > 1 {
		return false
	}

	return slashCount == 0 || ref[0] == '/'
}

func (c *PageCollections) getContentNode(context page.Page, isReflink bool, ref string) (*contentNode, error) {
	ref = filepath.ToSlash(strings.ToLower(strings.TrimSpace(ref)))
	if ref == "" {
		ref = "/"
	}
	inRef := ref
	navUp := strings.HasPrefix(ref, "..")
	var doSimpleLookup bool
	if isReflink || context == nil {
		doSimpleLookup = shouldDoSimpleLookup(ref)
	}

	if context != nil && !strings.HasPrefix(ref, "/") {
		// Try the page-relative path.
		var base string
		if context.File().IsZero() {
			base = context.SectionsPath()
		} else {
			meta := context.File().FileInfo().Meta()
			base = filepath.ToSlash(filepath.Dir(meta.Path()))
			if meta.Classifier() == files.ContentClassLeaf {
				// Bundles are stored in subfolders e.g. blog/mybundle/index.md,
				// so if the user has not explicitly asked to go up,
				// look on the "blog" level.
				if !navUp {
					base = path.Dir(base)
				}
			}
		}
		ref = path.Join("/", strings.ToLower(base), ref)
	}

	if !strings.HasPrefix(ref, "/") {
		ref = "/" + ref
	}

	m := c.pageMap

	// It's either a section, a page in a section or a taxonomy node.
	// Start with the most likely:
	n, name := c.getSectionOrPage(ref)
	if n != nil {
		return n, nil
	}

	if !strings.HasPrefix(inRef, "/") {
		// Many people will have "post/foo.md" in their content files.
		if n, _ := c.getSectionOrPage("/" + inRef); n != nil {
			return n, nil
		}
	}

	// Check if it's a taxonomy node
	s, v, found := m.taxonomies.LongestPrefix(ref)
	if found {
		if !m.onSameLevel(ref, s) {
			return nil, nil
		}
		return v.(*contentNode), nil
	}

	getByName := func(s string) (*contentNode, error) {
		n := m.pageReverseIndex.Get(s)
		if n != nil {
			if n == ambigousContentNode {
				return nil, fmt.Errorf("page reference %q is ambiguous", ref)
			}
			return n, nil
		}

		return nil, nil
	}

	var module string
	if context != nil && !context.File().IsZero() {
		module = context.File().FileInfo().Meta().Module()
	}

	if module == "" && !c.pageMap.s.home.File().IsZero() {
		module = c.pageMap.s.home.File().FileInfo().Meta().Module()
	}

	if module != "" {
		n, err := getByName(module + ref)
		if err != nil {
			return nil, err
		}
		if n != nil {
			return n, nil
		}
	}

	if !doSimpleLookup {
		return nil, nil
	}

	// Ref/relref supports this potentially ambigous lookup.
	return getByName(path.Base(name))

}

func (*PageCollections) findPagesByKindIn(kind string, inPages page.Pages) page.Pages {
	var pages page.Pages
	for _, p := range inPages {
		if p.Kind() == kind {
			pages = append(pages, p)
		}
	}
	return pages
}
