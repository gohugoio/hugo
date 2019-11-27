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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/cache"
	"github.com/gohugoio/hugo/resources/page"
)

// Used in the page cache to mark more than one hit for a given key.
var ambiguityFlag = &pageState{}

// PageCollections contains the page collections for a site.
type PageCollections struct {
	pagesMap *pagesMap

	// Includes absolute all pages (of all types), including drafts etc.
	rawAllPages pageStatePages

	// rawAllPages plus additional pages created during the build process.
	workAllPages pageStatePages

	// Includes headless bundles, i.e. bundles that produce no output for its content page.
	headlessPages pageStatePages

	// Lazy initialized page collections
	pages           *lazyPagesFactory
	regularPages    *lazyPagesFactory
	allPages        *lazyPagesFactory
	allRegularPages *lazyPagesFactory

	// The index for .Site.GetPage etc.
	pageIndex *cache.Lazy
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

// Get initializes the index if not already done so, then
// looks up the given page ref, returns nil if no value found.
func (c *PageCollections) getFromCache(ref string) (page.Page, error) {
	v, found, err := c.pageIndex.Get(ref)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	p := v.(page.Page)

	if p != ambiguityFlag {
		return p, nil
	}
	return nil, fmt.Errorf("page reference %q is ambiguous", ref)
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

func newPageCollections() *PageCollections {
	return newPageCollectionsFromPages(nil)
}

func newPageCollectionsFromPages(pages pageStatePages) *PageCollections {

	c := &PageCollections{rawAllPages: pages}

	c.pages = newLazyPagesFactory(func() page.Pages {
		pages := make(page.Pages, len(c.workAllPages))
		for i, p := range c.workAllPages {
			pages[i] = p
		}
		return pages
	})

	c.regularPages = newLazyPagesFactory(func() page.Pages {
		return c.findPagesByKindInWorkPages(page.KindPage, c.workAllPages)
	})

	c.pageIndex = cache.NewLazy(func() (map[string]interface{}, error) {
		index := make(map[string]interface{})

		add := func(ref string, p page.Page) {
			ref = strings.ToLower(ref)
			existing := index[ref]
			if existing == nil {
				index[ref] = p
			} else if existing != ambiguityFlag && existing != p {
				index[ref] = ambiguityFlag
			}
		}

		for _, pageCollection := range []pageStatePages{c.workAllPages, c.headlessPages} {
			for _, p := range pageCollection {
				if p.IsPage() {
					sourceRef := p.sourceRef()
					if sourceRef != "" {
						// index the canonical ref
						// e.g. /section/article.md
						add(sourceRef, p)
					}

					// Ref/Relref supports this potentially ambiguous lookup.
					add(p.File().LogicalName(), p)

					translationBaseName := p.File().TranslationBaseName()

					dir, _ := path.Split(sourceRef)
					dir = strings.TrimSuffix(dir, "/")

					if translationBaseName == "index" {
						add(dir, p)
						add(path.Base(dir), p)
					} else {
						add(translationBaseName, p)
					}

					// We need a way to get to the current language version.
					pathWithNoExtensions := path.Join(dir, translationBaseName)
					add(pathWithNoExtensions, p)
				} else {
					// index the canonical, unambiguous ref for any backing file
					// e.g. /section/_index.md
					sourceRef := p.sourceRef()
					if sourceRef != "" {
						add(sourceRef, p)
					}

					ref := p.SectionsPath()

					// index the canonical, unambiguous virtual ref
					// e.g. /section
					// (this may already have been indexed above)
					add("/"+ref, p)
				}
			}
		}

		return index, nil
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

// Case insensitive page lookup.
func (c *PageCollections) getPageNew(context page.Page, ref string) (page.Page, error) {
	var anError error

	ref = strings.ToLower(ref)

	// Absolute (content root relative) reference.
	if strings.HasPrefix(ref, "/") {
		p, err := c.getFromCache(ref)
		if err == nil && p != nil {
			return p, nil
		}
		if err != nil {
			anError = err
		}

	} else if context != nil {
		// Try the page-relative path.
		ppath := path.Join("/", strings.ToLower(context.SectionsPath()), ref)
		p, err := c.getFromCache(ppath)
		if err == nil && p != nil {
			return p, nil
		}
		if err != nil {
			anError = err
		}
	}

	if !strings.HasPrefix(ref, "/") {
		// Many people will have "post/foo.md" in their content files.
		p, err := c.getFromCache("/" + ref)
		if err == nil && p != nil {
			return p, nil
		}
		if err != nil {
			anError = err
		}
	}

	// Last try.
	ref = strings.TrimPrefix(ref, "/")
	p, err := c.getFromCache(ref)
	if err != nil {
		anError = err
	}

	if p == nil && anError != nil {
		return nil, wrapErr(errors.Wrap(anError, "failed to resolve ref"), context)
	}

	return p, nil
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

func (c *PageCollections) findPagesByKind(kind string) page.Pages {
	return c.findPagesByKindIn(kind, c.Pages())
}

func (c *PageCollections) findWorkPagesByKind(kind string) pageStatePages {
	var pages pageStatePages
	for _, p := range c.workAllPages {
		if p.Kind() == kind {
			pages = append(pages, p)
		}
	}
	return pages
}

func (*PageCollections) findPagesByKindInWorkPages(kind string, inPages pageStatePages) page.Pages {
	var pages page.Pages
	for _, p := range inPages {
		if p.Kind() == kind {
			pages = append(pages, p)
		}
	}
	return pages
}

func (c *PageCollections) addPage(page *pageState) {
	c.rawAllPages = append(c.rawAllPages, page)
}

func (c *PageCollections) removePageFilename(filename string) {
	if i := c.rawAllPages.findPagePosByFilename(filename); i >= 0 {
		c.clearResourceCacheForPage(c.rawAllPages[i])
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}

}

func (c *PageCollections) removePage(page *pageState) {
	if i := c.rawAllPages.findPagePos(page); i >= 0 {
		c.clearResourceCacheForPage(c.rawAllPages[i])
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}
}

func (c *PageCollections) replacePage(page *pageState) {
	// will find existing page that matches filepath and remove it
	c.removePage(page)
	c.addPage(page)
}

func (c *PageCollections) clearResourceCacheForPage(page *pageState) {
	if len(page.resources) > 0 {
		page.s.ResourceSpec.DeleteCacheByPrefix(page.targetPaths().SubResourceBaseTarget)
	}
}

func (c *PageCollections) assemblePagesMap(s *Site) error {

	c.pagesMap = newPagesMap(s)

	rootSections := make(map[string]bool)

	// Add all branch nodes first.
	for _, p := range c.rawAllPages {
		rootSections[p.Section()] = true
		if p.IsPage() {
			continue
		}
		c.pagesMap.addPage(p)
	}

	// Create missing home page and the first level sections if no
	// _index provided.
	s.home = c.pagesMap.getOrCreateHome()
	for k := range rootSections {
		c.pagesMap.createSectionIfNotExists(k)
	}

	// Attach the regular pages to their section.
	for _, p := range c.rawAllPages {
		if p.IsNode() {
			continue
		}
		c.pagesMap.addPage(p)
	}

	return nil
}

func (c *PageCollections) createWorkAllPages() error {
	c.workAllPages = make(pageStatePages, 0, len(c.rawAllPages))
	c.headlessPages = make(pageStatePages, 0)

	var (
		homeDates    *resource.Dates
		sectionDates *resource.Dates
		siteLastmod  time.Time
		siteLastDate time.Time

		sectionsParamId      = "mainSections"
		sectionsParamIdLower = strings.ToLower(sectionsParamId)
	)

	mainSections, mainSectionsFound := c.pagesMap.s.Info.Params()[sectionsParamIdLower]

	var (
		bucketsToRemove []string
		rootBuckets     []*pagesMapBucket
		walkErr         error
	)

	c.pagesMap.r.Walk(func(s string, v interface{}) bool {
		bucket := v.(*pagesMapBucket)
		parentBucket := c.pagesMap.parentBucket(s)

		if parentBucket != nil {

			if !mainSectionsFound && strings.Count(s, "/") == 1 && bucket.owner.IsSection() {
				// Root section
				rootBuckets = append(rootBuckets, bucket)
			}
		}

		if bucket.owner.IsHome() {
			if resource.IsZeroDates(bucket.owner) {
				// Calculate dates from the page tree.
				homeDates = &bucket.owner.m.Dates
			}
		}

		sectionDates = nil
		if resource.IsZeroDates(bucket.owner) {
			sectionDates = &bucket.owner.m.Dates
		}

		if parentBucket != nil {
			bucket.parent = parentBucket
			if bucket.owner.IsSection() {
				parentBucket.bucketSections = append(parentBucket.bucketSections, bucket)
			}
		}

		if bucket.isEmpty() {
			if bucket.owner.IsSection() && bucket.owner.File().IsZero() {
				// Check for any nested section.
				var hasDescendant bool
				c.pagesMap.r.WalkPrefix(s, func(ss string, v interface{}) bool {
					if s != ss {
						hasDescendant = true
						return true
					}
					return false
				})
				if !hasDescendant {
					// This is an auto-created section with, now, nothing in it.
					bucketsToRemove = append(bucketsToRemove, s)
					return false
				}
			}
		}

		if !bucket.disabled {
			c.workAllPages = append(c.workAllPages, bucket.owner)
		}

		if !bucket.view {
			for _, p := range bucket.headlessPages {
				ps := p.(*pageState)
				ps.parent = bucket.owner
				c.headlessPages = append(c.headlessPages, ps)
			}
			for _, p := range bucket.pages {
				ps := p.(*pageState)
				ps.parent = bucket.owner
				c.workAllPages = append(c.workAllPages, ps)

				if homeDates != nil {
					homeDates.UpdateDateAndLastmodIfAfter(ps)
				}

				if sectionDates != nil {
					sectionDates.UpdateDateAndLastmodIfAfter(ps)
				}

				if p.Lastmod().After(siteLastmod) {
					siteLastmod = p.Lastmod()
				}
				if p.Date().After(siteLastDate) {
					siteLastDate = p.Date()
				}
			}
		}

		return false
	})

	if walkErr != nil {
		return walkErr
	}

	c.pagesMap.s.lastmod = siteLastmod

	if !mainSectionsFound {

		// Calculare main section
		var (
			maxRootBucketWeight int
			maxRootBucket       *pagesMapBucket
		)

		for _, b := range rootBuckets {
			weight := len(b.pages) + (len(b.bucketSections) * 5)
			if weight >= maxRootBucketWeight {
				maxRootBucket = b
				maxRootBucketWeight = weight
			}
		}

		if maxRootBucket != nil {
			// Try to make this as backwards compatible as possible.
			mainSections = []string{maxRootBucket.owner.Section()}
		}
	}

	c.pagesMap.s.Info.Params()[sectionsParamId] = mainSections
	c.pagesMap.s.Info.Params()[sectionsParamIdLower] = mainSections

	for _, key := range bucketsToRemove {
		c.pagesMap.r.Delete(key)
	}

	sort.Sort(c.workAllPages)

	return nil
}
