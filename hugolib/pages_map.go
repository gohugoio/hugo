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

	"github.com/gohugoio/hugo/common/maps"

	radix "github.com/armon/go-radix"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/resources/page"
)

func newPagesMap(s *Site) *pagesMap {
	return &pagesMap{
		r: radix.New(),
		s: s,
	}
}

type pagesMap struct {
	r *radix.Tree
	s *Site
}

func (m *pagesMap) Get(key string) *pagesMapBucket {
	key = m.cleanKey(key)
	v, found := m.r.Get(key)
	if !found {
		return nil
	}

	return v.(*pagesMapBucket)
}

func (m *pagesMap) getKey(p *pageState) string {
	if !p.File().IsZero() {
		return m.cleanKey(p.File().Dir())
	}
	return m.cleanKey(p.SectionsPath())
}

func (m *pagesMap) getOrCreateHome() *pageState {
	var home *pageState
	b, found := m.r.Get("/")
	if !found {
		home = m.s.newPage(page.KindHome)
		m.addBucketFor("/", home, nil)
	} else {
		home = b.(*pagesMapBucket).owner
	}

	return home
}

func (m *pagesMap) initPageMeta(p *pageState, bucket *pagesMapBucket) error {
	var err error
	p.metaInit.Do(func() {
		if p.metaInitFn != nil {
			err = p.metaInitFn(bucket)
		}
	})
	return err
}

func (m *pagesMap) initPageMetaFor(prefix string, bucket *pagesMapBucket) error {
	parentBucket := m.parentBucket(prefix)

	m.mergeCascades(bucket, parentBucket)

	if err := m.initPageMeta(bucket.owner, bucket); err != nil {
		return err
	}

	if !bucket.view {
		for _, p := range bucket.pages {
			ps := p.(*pageState)
			if err := m.initPageMeta(ps, bucket); err != nil {
				return err
			}

			for _, p := range ps.resources.ByType(pageResourceType) {
				if err := m.initPageMeta(p.(*pageState), bucket); err != nil {
					return err
				}
			}
		}

		// Now that the metadata is initialized (with dates, draft set etc.)
		// we can remove the pages that we for some reason should not include
		// in this build.
		tmp := bucket.pages[:0]
		for _, x := range bucket.pages {
			if m.s.shouldBuild(x) {
				if x.(*pageState).m.headless {
					bucket.headlessPages = append(bucket.headlessPages, x)
				} else {
					tmp = append(tmp, x)
				}

			}
		}
		bucket.pages = tmp
	}

	return nil
}

func (m *pagesMap) createSectionIfNotExists(section string) {
	key := m.cleanKey(section)
	_, found := m.r.Get(key)
	if !found {
		kind := m.s.kindFromSectionPath(section)
		p := m.s.newPage(kind, section)
		m.addBucketFor(key, p, nil)
	}
}

func (m *pagesMap) addBucket(p *pageState) {
	key := m.getKey(p)

	m.addBucketFor(key, p, nil)
}

func (m *pagesMap) addBucketFor(key string, p *pageState, meta map[string]interface{}) *pagesMapBucket {
	var isView bool
	switch p.Kind() {
	case page.KindTaxonomy, page.KindTaxonomyTerm:
		isView = true
	}

	disabled := !m.s.isEnabled(p.Kind())

	var cascade map[string]interface{}
	if p.bucket != nil {
		cascade = p.bucket.cascade
	}

	bucket := &pagesMapBucket{
		owner:    p,
		view:     isView,
		cascade:  cascade,
		meta:     meta,
		disabled: disabled,
	}

	p.bucket = bucket

	m.r.Insert(key, bucket)

	return bucket
}

func (m *pagesMap) addPage(p *pageState) {
	if !p.IsPage() {
		m.addBucket(p)
		return
	}

	if !m.s.isEnabled(page.KindPage) {
		return
	}

	key := m.getKey(p)

	var bucket *pagesMapBucket

	_, v, found := m.r.LongestPrefix(key)
	if !found {
		panic(fmt.Sprintf("[BUG] bucket with key %q not found", key))
	}

	bucket = v.(*pagesMapBucket)
	bucket.pages = append(bucket.pages, p)
}

func (m *pagesMap) assemblePageMeta() error {
	var walkErr error
	m.r.Walk(func(s string, v interface{}) bool {
		bucket := v.(*pagesMapBucket)

		if err := m.initPageMetaFor(s, bucket); err != nil {
			walkErr = err
			return true
		}
		return false
	})

	return walkErr
}

func (m *pagesMap) assembleTaxonomies(s *Site) error {
	s.Taxonomies = make(TaxonomyList)

	type bucketKey struct {
		plural  string
		termKey string
	}

	// Temporary cache.
	taxonomyBuckets := make(map[bucketKey]*pagesMapBucket)

	for singular, plural := range s.siteCfg.taxonomiesConfig {
		s.Taxonomies[plural] = make(Taxonomy)
		bkey := bucketKey{
			plural: plural,
		}

		bucket := m.Get(plural)

		if bucket == nil {
			// Create the page and bucket
			n := s.newPage(page.KindTaxonomyTerm, plural)

			key := m.cleanKey(plural)
			bucket = m.addBucketFor(key, n, nil)
			if err := m.initPageMetaFor(key, bucket); err != nil {
				return err
			}
		}

		if bucket.meta == nil {
			bucket.meta = map[string]interface{}{
				"singular": singular,
				"plural":   plural,
			}
		}

		// Add it to the temporary cache.
		taxonomyBuckets[bkey] = bucket

		// Taxonomy entries used in page front matter will be picked up later,
		// but there may be some yet to be used.
		pluralPrefix := m.cleanKey(plural) + "/"
		m.r.WalkPrefix(pluralPrefix, func(k string, v interface{}) bool {
			tb := v.(*pagesMapBucket)
			termKey := strings.TrimPrefix(k, pluralPrefix)
			if tb.meta == nil {
				tb.meta = map[string]interface{}{
					"singular": singular,
					"plural":   plural,
					"term":     tb.owner.Title(),
					"termKey":  termKey,
				}
			}

			bucket.pages = append(bucket.pages, tb.owner)
			bkey.termKey = termKey
			taxonomyBuckets[bkey] = tb

			return false
		})

	}

	addTaxonomy := func(singular, plural, term string, weight int, p page.Page) error {
		bkey := bucketKey{
			plural: plural,
		}

		termKey := s.getTaxonomyKey(term)

		b1 := taxonomyBuckets[bkey]

		var b2 *pagesMapBucket
		bkey.termKey = termKey
		b, found := taxonomyBuckets[bkey]
		if found {
			b2 = b
		} else {

			// Create the page and bucket
			n := s.newTaxonomyPage(term, plural, termKey)
			meta := map[string]interface{}{
				"singular": singular,
				"plural":   plural,
				"term":     term,
				"termKey":  termKey,
			}

			key := m.cleanKey(path.Join(plural, termKey))
			b2 = m.addBucketFor(key, n, meta)
			if err := m.initPageMetaFor(key, b2); err != nil {
				return err
			}
			b1.pages = append(b1.pages, b2.owner)
			taxonomyBuckets[bkey] = b2

		}

		w := page.NewWeightedPage(weight, p, b2.owner)

		s.Taxonomies[plural].add(termKey, w)

		b1.owner.m.Dates.UpdateDateAndLastmodIfAfter(p)
		b2.owner.m.Dates.UpdateDateAndLastmodIfAfter(p)

		return nil
	}

	m.r.Walk(func(k string, v interface{}) bool {
		b := v.(*pagesMapBucket)
		if b.view {
			return false
		}

		for singular, plural := range s.siteCfg.taxonomiesConfig {
			for _, p := range b.pages {

				vals := getParam(p, plural, false)

				w := getParamToLower(p, plural+"_weight")
				weight, err := cast.ToIntE(w)
				if err != nil {
					m.s.Log.ERROR.Printf("Unable to convert taxonomy weight %#v to int for %q", w, p.Path())
					// weight will equal zero, so let the flow continue
				}

				if vals != nil {
					if v, ok := vals.([]string); ok {
						for _, idx := range v {
							if err := addTaxonomy(singular, plural, idx, weight, p); err != nil {
								m.s.Log.ERROR.Printf("Failed to add taxonomy %q for %q: %s", plural, p.Path(), err)
							}
						}
					} else if v, ok := vals.(string); ok {
						if err := addTaxonomy(singular, plural, v, weight, p); err != nil {
							m.s.Log.ERROR.Printf("Failed to add taxonomy %q for %q: %s", plural, p.Path(), err)
						}
					} else {
						m.s.Log.ERROR.Printf("Invalid %s in %q\n", plural, p.Path())
					}
				}

			}
		}
		return false
	})

	for _, plural := range s.siteCfg.taxonomiesConfig {
		for k := range s.Taxonomies[plural] {
			s.Taxonomies[plural][k].Sort()
		}
	}

	return nil
}

func (m *pagesMap) cleanKey(key string) string {
	key = filepath.ToSlash(strings.ToLower(key))
	key = strings.Trim(key, "/")
	return "/" + key
}

func (m *pagesMap) mergeCascades(b1, b2 *pagesMapBucket) {
	if b1.cascade == nil {
		b1.cascade = make(maps.Params)
	}
	if b2 != nil && b2.cascade != nil {
		for k, v := range b2.cascade {
			if _, found := b1.cascade[k]; !found {
				b1.cascade[k] = v
			}
		}
	}
}

func (m *pagesMap) parentBucket(prefix string) *pagesMapBucket {
	if prefix == "/" {
		return nil
	}
	_, parentv, found := m.r.LongestPrefix(path.Dir(prefix))
	if !found {
		panic(fmt.Sprintf("[BUG] parent bucket not found for %q", prefix))
	}
	return parentv.(*pagesMapBucket)

}

func (m *pagesMap) withEveryPage(f func(p *pageState)) {
	m.r.Walk(func(k string, v interface{}) bool {
		b := v.(*pagesMapBucket)
		f(b.owner)
		if !b.view {
			for _, p := range b.pages {
				f(p.(*pageState))
			}
		}

		return false
	})
}

type pagesMapBucket struct {
	// Set if the pages in this bucket is also present in another bucket.
	view bool

	// Some additional metatadata attached to this node.
	meta map[string]interface{}

	// Cascading front matter.
	cascade map[string]interface{}

	owner *pageState // The branch node

	// When disableKinds is enabled for this node.
	disabled bool

	// Used to navigate the sections tree
	parent         *pagesMapBucket
	bucketSections []*pagesMapBucket

	pagesInit     sync.Once
	pages         page.Pages
	headlessPages page.Pages

	pagesAndSectionsInit sync.Once
	pagesAndSections     page.Pages

	sectionsInit sync.Once
	sections     page.Pages
}

func (b *pagesMapBucket) isEmpty() bool {
	return len(b.pages) == 0 && len(b.bucketSections) == 0
}

func (b *pagesMapBucket) getPages() page.Pages {
	b.pagesInit.Do(func() {
		page.SortByDefault(b.pages)
	})
	return b.pages
}

func (b *pagesMapBucket) getPagesAndSections() page.Pages {
	b.pagesAndSectionsInit.Do(func() {
		var pas page.Pages
		pas = append(pas, b.getPages()...)
		for _, p := range b.bucketSections {
			pas = append(pas, p.owner)
		}
		b.pagesAndSections = pas
		page.SortByDefault(b.pagesAndSections)
	})
	return b.pagesAndSections
}

func (b *pagesMapBucket) getSections() page.Pages {
	b.sectionsInit.Do(func() {
		for _, p := range b.bucketSections {
			b.sections = append(b.sections, p.owner)
		}
		page.SortByDefault(b.sections)
	})

	return b.sections
}
