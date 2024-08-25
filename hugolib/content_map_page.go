// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"
)

var pagePredicates = struct {
	KindPage         predicate.P[*pageState]
	KindSection      predicate.P[*pageState]
	KindHome         predicate.P[*pageState]
	KindTerm         predicate.P[*pageState]
	ShouldListLocal  predicate.P[*pageState]
	ShouldListGlobal predicate.P[*pageState]
	ShouldListAny    predicate.P[*pageState]
	ShouldLink       predicate.P[page.Page]
}{
	KindPage: func(p *pageState) bool {
		return p.Kind() == kinds.KindPage
	},
	KindSection: func(p *pageState) bool {
		return p.Kind() == kinds.KindSection
	},
	KindHome: func(p *pageState) bool {
		return p.Kind() == kinds.KindHome
	},
	KindTerm: func(p *pageState) bool {
		return p.Kind() == kinds.KindTerm
	},
	ShouldListLocal: func(p *pageState) bool {
		return p.m.shouldList(false)
	},
	ShouldListGlobal: func(p *pageState) bool {
		return p.m.shouldList(true)
	},
	ShouldListAny: func(p *pageState) bool {
		return p.m.shouldListAny()
	},
	ShouldLink: func(p page.Page) bool {
		return !p.(*pageState).m.noLink()
	},
}

type pageMap struct {
	i int
	s *Site

	// Main storage for all pages.
	*pageTrees

	// Used for simple page lookups by name, e.g. "mypage.md" or "mypage".
	pageReverseIndex *contentTreeReverseIndex

	cachePages1            *dynacache.Partition[string, page.Pages]
	cachePages2            *dynacache.Partition[string, page.Pages]
	cacheResources         *dynacache.Partition[string, resource.Resources]
	cacheGetTerms          *dynacache.Partition[string, map[string]page.Pages]
	cacheContentRendered   *dynacache.Partition[string, *resources.StaleValue[contentSummary]]
	cacheContentPlain      *dynacache.Partition[string, *resources.StaleValue[contentPlainPlainWords]]
	contentTableOfContents *dynacache.Partition[string, *resources.StaleValue[contentTableOfContents]]

	contentDataFileSeenItems *maps.Cache[string, map[uint64]bool]

	cfg contentMapConfig
}

// pageTrees holds pages and resources in a tree structure for all sites/languages.
// Each site gets its own tree set via the Shape method.
type pageTrees struct {
	// This tree contains all Pages.
	// This include regular pages, sections, taxonomies and so on.
	// Note that all of these trees share the same key structure,
	// so you can take a leaf Page key and do a prefix search
	// with key + "/" to get all of its resources.
	treePages *doctree.NodeShiftTree[contentNodeI]

	// This tree contains Resources bundled in pages.
	treeResources *doctree.NodeShiftTree[contentNodeI]

	// All pages and resources.
	treePagesResources doctree.WalkableTrees[contentNodeI]

	// This tree contains all taxonomy entries, e.g "/tags/blue/page1"
	treeTaxonomyEntries *doctree.TreeShiftTree[*weightedContentNode]

	// Stores the state for _content.gotmpl files.
	// Mostly releveant for rebuilds.
	treePagesFromTemplateAdapters *doctree.TreeShiftTree[*pagesfromdata.PagesFromTemplate]

	// A slice of the resource trees.
	resourceTrees doctree.MutableTrees
}

// collectAndMarkStaleIdentities collects all identities from in all trees matching the given key.
// We currently re-read all page/resources for all languages that share the same path,
// so we mark all entries as stale (which will trigger cache invalidation), then
// return the first.
func (t *pageTrees) collectAndMarkStaleIdentities(p *paths.Path) []identity.Identity {
	key := p.Base()
	var ids []identity.Identity
	// We need only one identity sample per dimension.
	nCount := 0
	cb := func(n contentNodeI) bool {
		if n == nil {
			return false
		}
		n.MarkStale()
		if nCount > 0 {
			return true
		}
		nCount++
		n.ForEeachIdentity(func(id identity.Identity) bool {
			ids = append(ids, id)
			return false
		})

		return false
	}
	tree := t.treePages
	nCount = 0
	tree.ForEeachInDimension(key, doctree.DimensionLanguage.Index(),
		cb,
	)

	tree = t.treeResources
	nCount = 0
	tree.ForEeachInDimension(key, doctree.DimensionLanguage.Index(),
		cb,
	)

	if p.Component() == files.ComponentFolderContent {
		// It may also be a bundled content resource.
		key := p.ForBundleType(paths.PathTypeContentResource).Base()
		tree = t.treeResources
		nCount = 0
		tree.ForEeachInDimension(key, doctree.DimensionLanguage.Index(),
			cb,
		)

	}
	return ids
}

// collectIdentitiesSurrounding collects all identities surrounding the given key.
func (t *pageTrees) collectIdentitiesSurrounding(key string, maxSamplesPerTree int) []identity.Identity {
	ids := t.collectIdentitiesSurroundingIn(key, maxSamplesPerTree, t.treePages)
	ids = append(ids, t.collectIdentitiesSurroundingIn(key, maxSamplesPerTree, t.treeResources)...)
	return ids
}

func (t *pageTrees) collectIdentitiesSurroundingIn(key string, maxSamples int, tree *doctree.NodeShiftTree[contentNodeI]) []identity.Identity {
	var ids []identity.Identity
	section, ok := tree.LongestPrefixAll(path.Dir(key))
	if ok {
		count := 0
		prefix := section + "/"
		level := strings.Count(prefix, "/")
		tree.WalkPrefixRaw(prefix, func(s string, n contentNodeI) bool {
			if level != strings.Count(s, "/") {
				return false
			}
			n.ForEeachIdentity(func(id identity.Identity) bool {
				ids = append(ids, id)
				return false
			})
			count++
			return count > maxSamples
		})
	}

	return ids
}

func (t *pageTrees) DeletePageAndResourcesBelow(ss ...string) {
	commit1 := t.resourceTrees.Lock(true)
	defer commit1()
	commit2 := t.treePages.Lock(true)
	defer commit2()
	for _, s := range ss {
		t.resourceTrees.DeletePrefix(paths.AddTrailingSlash(s))
		t.treePages.Delete(s)
	}
}

// Shape shapes all trees in t to the given dimension.
func (t pageTrees) Shape(d, v int) *pageTrees {
	t.treePages = t.treePages.Shape(d, v)
	t.treeResources = t.treeResources.Shape(d, v)
	t.treeTaxonomyEntries = t.treeTaxonomyEntries.Shape(d, v)
	t.treePagesFromTemplateAdapters = t.treePagesFromTemplateAdapters.Shape(d, v)
	t.createMutableTrees()

	return &t
}

func (t *pageTrees) createMutableTrees() {
	t.treePagesResources = doctree.WalkableTrees[contentNodeI]{
		t.treePages,
		t.treeResources,
	}

	t.resourceTrees = doctree.MutableTrees{
		t.treeResources,
	}
}

var (
	_ resource.Identifier = pageMapQueryPagesInSection{}
	_ resource.Identifier = pageMapQueryPagesBelowPath{}
)

type pageMapQueryPagesInSection struct {
	pageMapQueryPagesBelowPath

	Recursive   bool
	IncludeSelf bool
}

func (q pageMapQueryPagesInSection) Key() string {
	return "gagesInSection" + "/" + q.pageMapQueryPagesBelowPath.Key() + "/" + strconv.FormatBool(q.Recursive) + "/" + strconv.FormatBool(q.IncludeSelf)
}

// This needs to be hashable.
type pageMapQueryPagesBelowPath struct {
	Path string

	// Additional identifier for this query.
	// Used as part of the cache key.
	KeyPart string

	// Page inclusion filter.
	// May be nil.
	Include predicate.P[*pageState]
}

func (q pageMapQueryPagesBelowPath) Key() string {
	return q.Path + "/" + q.KeyPart
}

// Apply fn to all pages in m matching the given predicate.
// fn may return true to stop the walk.
func (m *pageMap) forEachPage(include predicate.P[*pageState], fn func(p *pageState) (bool, error)) error {
	if include == nil {
		include = func(p *pageState) bool {
			return true
		}
	}
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     m.treePages,
		LockType: doctree.LockTypeRead,
		Handle: func(key string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if p, ok := n.(*pageState); ok && include(p) {
				if terminate, err := fn(p); terminate || err != nil {
					return terminate, err
				}
			}
			return false, nil
		},
	}

	return w.Walk(context.Background())
}

func (m *pageMap) forEeachPageIncludingBundledPages(include predicate.P[*pageState], fn func(p *pageState) (bool, error)) error {
	if include == nil {
		include = func(p *pageState) bool {
			return true
		}
	}

	if err := m.forEachPage(include, fn); err != nil {
		return err
	}

	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     m.treeResources,
		LockType: doctree.LockTypeRead,
		Handle: func(key string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if rs, ok := n.(*resourceSource); ok {
				if p, ok := rs.r.(*pageState); ok && include(p) {
					if terminate, err := fn(p); terminate || err != nil {
						return terminate, err
					}
				}
			}
			return false, nil
		},
	}

	return w.Walk(context.Background())
}

func (m *pageMap) getOrCreatePagesFromCache(
	cache *dynacache.Partition[string, page.Pages],
	key string, create func(string) (page.Pages, error),
) (page.Pages, error) {
	if cache == nil {
		cache = m.cachePages1
	}
	return cache.GetOrCreate(key, create)
}

func (m *pageMap) getPagesInSection(q pageMapQueryPagesInSection) page.Pages {
	cacheKey := q.Key()

	pages, err := m.getOrCreatePagesFromCache(nil, cacheKey, func(string) (page.Pages, error) {
		prefix := paths.AddTrailingSlash(q.Path)

		var (
			pas         page.Pages
			otherBranch string
		)

		include := q.Include
		if include == nil {
			include = pagePredicates.ShouldListLocal
		}

		w := &doctree.NodeShiftTreeWalker[contentNodeI]{
			Tree:   m.treePages,
			Prefix: prefix,
		}

		w.Handle = func(key string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if q.Recursive {
				if p, ok := n.(*pageState); ok && include(p) {
					pas = append(pas, p)
				}
				return false, nil
			}

			if p, ok := n.(*pageState); ok && include(p) {
				pas = append(pas, p)
			}

			if n.isContentNodeBranch() {
				currentBranch := key + "/"
				if otherBranch == "" || otherBranch != currentBranch {
					w.SkipPrefix(currentBranch)
				}
				otherBranch = currentBranch
			}
			return false, nil
		}

		err := w.Walk(context.Background())

		if err == nil {
			if q.IncludeSelf {
				if n := m.treePages.Get(q.Path); n != nil {
					if p, ok := n.(*pageState); ok && include(p) {
						pas = append(pas, p)
					}
				}
			}
			page.SortByDefault(pas)
		}

		return pas, err
	})
	if err != nil {
		panic(err)
	}

	return pages
}

func (m *pageMap) getPagesWithTerm(q pageMapQueryPagesBelowPath) page.Pages {
	key := q.Key()

	v, err := m.cachePages1.GetOrCreate(key, func(string) (page.Pages, error) {
		var pas page.Pages
		include := q.Include
		if include == nil {
			include = pagePredicates.ShouldListLocal
		}

		err := m.treeTaxonomyEntries.WalkPrefix(
			doctree.LockTypeNone,
			paths.AddTrailingSlash(q.Path),
			func(s string, n *weightedContentNode) (bool, error) {
				p := n.n.(*pageState)
				if !include(p) {
					return false, nil
				}
				pas = append(pas, pageWithWeight0{n.weight, p})
				return false, nil
			},
		)
		if err != nil {
			return nil, err
		}

		page.SortByDefault(pas)

		return pas, nil
	})
	if err != nil {
		panic(err)
	}

	return v
}

func (m *pageMap) getTermsForPageInTaxonomy(path, taxonomy string) page.Pages {
	prefix := paths.AddLeadingSlash(taxonomy)

	termPages, err := m.cacheGetTerms.GetOrCreate(prefix, func(string) (map[string]page.Pages, error) {
		mm := make(map[string]page.Pages)
		err := m.treeTaxonomyEntries.WalkPrefix(
			doctree.LockTypeNone,
			paths.AddTrailingSlash(prefix),
			func(s string, n *weightedContentNode) (bool, error) {
				mm[n.n.Path()] = append(mm[n.n.Path()], n.term)
				return false, nil
			},
		)
		if err != nil {
			return nil, err
		}

		// Sort the terms.
		for _, v := range mm {
			page.SortByDefault(v)
		}

		return mm, nil
	})
	if err != nil {
		panic(err)
	}

	return termPages[path]
}

func (m *pageMap) forEachResourceInPage(
	ps *pageState,
	lockType doctree.LockType,
	exact bool,
	handle func(resourceKey string, n contentNodeI, match doctree.DimensionFlag) (bool, error),
) error {
	keyPage := ps.Path()
	if keyPage == "/" {
		keyPage = ""
	}
	prefix := paths.AddTrailingSlash(ps.Path())
	isBranch := ps.IsNode()

	rw := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     m.treeResources,
		Prefix:   prefix,
		LockType: lockType,
		Exact:    exact,
	}

	rw.Handle = func(resourceKey string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
		if isBranch {
			// A resourceKey always represents a filename with extension.
			// A page key points to the logical path of a page, which when sourced from the filesystem
			// may represent a directory (bundles) or a single content file (e.g. p1.md).
			// So, to avoid any overlapping ambiguity, we start looking from the owning directory.
			ownerKey, _ := m.treePages.LongestPrefixAll(path.Dir(resourceKey))
			if ownerKey != keyPage {
				// Stop walking downwards, someone else owns this resource.
				rw.SkipPrefix(ownerKey + "/")
				return false, nil
			}

		}
		return handle(resourceKey, n, match)
	}

	return rw.Walk(context.Background())
}

func (m *pageMap) getResourcesForPage(ps *pageState) (resource.Resources, error) {
	var res resource.Resources
	m.forEachResourceInPage(ps, doctree.LockTypeNone, false, func(resourceKey string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
		rs := n.(*resourceSource)
		if rs.r != nil {
			res = append(res, rs.r)
		}
		return false, nil
	})
	return res, nil
}

func (m *pageMap) getOrCreateResourcesForPage(ps *pageState) resource.Resources {
	keyPage := ps.Path()
	if keyPage == "/" {
		keyPage = ""
	}
	key := keyPage + "/get-resources-for-page"
	v, err := m.cacheResources.GetOrCreate(key, func(string) (resource.Resources, error) {
		res, err := m.getResourcesForPage(ps)
		if err != nil {
			return nil, err
		}

		if translationKey := ps.m.pageConfig.TranslationKey; translationKey != "" {
			// This this should not be a very common case.
			// Merge in resources from the other languages.
			translatedPages, _ := m.s.h.translationKeyPages.Get(translationKey)
			for _, tp := range translatedPages {
				if tp == ps {
					continue
				}
				tps := tp.(*pageState)
				// Make sure we query from the correct language root.
				res2, err := tps.s.pageMap.getResourcesForPage(tps)
				if err != nil {
					return nil, err
				}
				// Add if Name not already in res.
				for _, r := range res2 {
					var found bool
					for _, r2 := range res {
						if resource.NameNormalizedOrName(r2) == resource.NameNormalizedOrName(r) {
							found = true
							break
						}
					}
					if !found {
						res = append(res, r)
					}
				}
			}
		}

		lessFunc := func(i, j int) bool {
			ri, rj := res[i], res[j]
			if ri.ResourceType() < rj.ResourceType() {
				return true
			}

			p1, ok1 := ri.(page.Page)
			p2, ok2 := rj.(page.Page)

			if ok1 != ok2 {
				// Pull pages behind other resources.

				return ok2
			}

			if ok1 {
				return page.DefaultPageSort(p1, p2)
			}

			// Make sure not to use RelPermalink or any of the other methods that
			// trigger lazy publishing.
			return ri.Name() < rj.Name()
		}

		sort.SliceStable(res, lessFunc)

		if len(ps.m.pageConfig.ResourcesMeta) > 0 {
			for i, r := range res {
				res[i] = resources.CloneWithMetadataFromMapIfNeeded(ps.m.pageConfig.ResourcesMeta, r)
			}
			sort.SliceStable(res, lessFunc)
		}

		return res, nil
	})
	if err != nil {
		panic(err)
	}

	return v
}

type weightedContentNode struct {
	n      contentNodeI
	weight int
	term   *pageWithOrdinal
}

type buildStateReseter interface {
	resetBuildState()
}

type contentNodeI interface {
	identity.IdentityProvider
	identity.ForEeachIdentityProvider
	Path() string
	isContentNodeBranch() bool
	buildStateReseter
	resource.StaleMarker
}

var _ contentNodeI = (*contentNodeIs)(nil)

type contentNodeIs []contentNodeI

func (n contentNodeIs) Path() string {
	return n[0].Path()
}

func (n contentNodeIs) isContentNodeBranch() bool {
	return n[0].isContentNodeBranch()
}

func (n contentNodeIs) GetIdentity() identity.Identity {
	return n[0].GetIdentity()
}

func (n contentNodeIs) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, nn := range n {
		if nn != nil {
			if nn.ForEeachIdentity(f) {
				return true
			}
		}
	}
	return false
}

func (n contentNodeIs) resetBuildState() {
	for _, nn := range n {
		if nn != nil {
			nn.resetBuildState()
		}
	}
}

func (n contentNodeIs) MarkStale() {
	for _, nn := range n {
		resource.MarkStale(nn)
	}
}

type contentNodeShifter struct {
	numLanguages int
}

func (s *contentNodeShifter) Delete(n contentNodeI, dimension doctree.Dimension) (contentNodeI, bool, bool) {
	lidx := dimension[0]
	switch v := n.(type) {
	case contentNodeIs:
		deleted := v[lidx]
		resource.MarkStale(deleted)
		wasDeleted := deleted != nil
		v[lidx] = nil
		isEmpty := true
		for _, vv := range v {
			if vv != nil {
				isEmpty = false
				break
			}
		}
		return deleted, wasDeleted, isEmpty
	case resourceSources:
		deleted := v[lidx]
		resource.MarkStale(deleted)
		wasDeleted := deleted != nil
		v[lidx] = nil
		isEmpty := true
		for _, vv := range v {
			if vv != nil {
				isEmpty = false
				break
			}
		}
		return deleted, wasDeleted, isEmpty
	case *resourceSource:
		if lidx != v.LangIndex() {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	case *pageState:
		if lidx != v.s.languagei {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	default:
		panic(fmt.Sprintf("unknown type %T", n))
	}
}

func (s *contentNodeShifter) Shift(n contentNodeI, dimension doctree.Dimension, exact bool) (contentNodeI, bool, doctree.DimensionFlag) {
	lidx := dimension[0]
	// How accurate is the match.
	accuracy := doctree.DimensionLanguage
	switch v := n.(type) {
	case contentNodeIs:
		if len(v) == 0 {
			panic("empty contentNodeIs")
		}
		vv := v[lidx]
		if vv != nil {
			return vv, true, accuracy
		}
		return nil, false, 0
	case resourceSources:
		vv := v[lidx]
		if vv != nil {
			return vv, true, doctree.DimensionLanguage
		}
		if exact {
			return nil, false, 0
		}
		// For non content resources, pick the first match.
		for _, vv := range v {
			if vv != nil {
				if vv.isPage() {
					return nil, false, 0
				}
				return vv, true, 0
			}
		}
	case *resourceSource:
		if v.LangIndex() == lidx {
			return v, true, doctree.DimensionLanguage
		}
		if !v.isPage() && !exact {
			return v, true, 0
		}
	case *pageState:
		if v.s.languagei == lidx {
			return n, true, doctree.DimensionLanguage
		}
	default:
		panic(fmt.Sprintf("unknown type %T", n))
	}
	return nil, false, 0
}

func (s *contentNodeShifter) ForEeachInDimension(n contentNodeI, d int, f func(contentNodeI) bool) {
	if d != doctree.DimensionLanguage.Index() {
		panic("only language dimension supported")
	}

	switch vv := n.(type) {
	case contentNodeIs:
		for _, v := range vv {
			if v != nil {
				if f(v) {
					return
				}
			}
		}
	default:
		f(vv)
	}
}

func (s *contentNodeShifter) InsertInto(old, new contentNodeI, dimension doctree.Dimension) (contentNodeI, contentNodeI, bool) {
	langi := dimension[doctree.DimensionLanguage.Index()]
	switch vv := old.(type) {
	case *pageState:
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.s.languagei == newp.s.languagei && newp.s.languagei == langi {
			return new, vv, true
		}
		is := make(contentNodeIs, s.numLanguages)
		is[vv.s.languagei] = old
		is[langi] = new
		return is, old, false
	case contentNodeIs:
		oldv := vv[langi]
		vv[langi] = new
		return vv, oldv, oldv != nil
	case resourceSources:
		oldv := vv[langi]
		vv[langi] = new.(*resourceSource)
		return vv, oldv, oldv != nil
	case *resourceSource:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.LangIndex() == newp.LangIndex() && newp.LangIndex() == langi {
			return new, vv, true
		}
		rs := make(resourceSources, s.numLanguages)
		rs[vv.LangIndex()] = vv
		rs[langi] = newp
		return rs, vv, false

	default:
		panic(fmt.Sprintf("unknown type %T", old))
	}
}

func (s *contentNodeShifter) Insert(old, new contentNodeI) (contentNodeI, contentNodeI, bool) {
	switch vv := old.(type) {
	case *pageState:
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.s.languagei == newp.s.languagei {
			if newp != old {
				resource.MarkStale(old)
			}
			return new, vv, true
		}
		is := make(contentNodeIs, s.numLanguages)
		is[newp.s.languagei] = new
		is[vv.s.languagei] = old
		return is, old, false
	case contentNodeIs:
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		oldp := vv[newp.s.languagei]
		if oldp != newp {
			resource.MarkStale(oldp)
		}
		vv[newp.s.languagei] = new
		return vv, oldp, oldp != nil
	case *resourceSource:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.LangIndex() == newp.LangIndex() {
			if vv != newp {
				resource.MarkStale(vv)
			}
			return new, vv, true
		}
		rs := make(resourceSources, s.numLanguages)
		rs[newp.LangIndex()] = newp
		rs[vv.LangIndex()] = vv
		return rs, vv, false
	case resourceSources:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		oldp := vv[newp.LangIndex()]
		if oldp != newp {
			resource.MarkStale(oldp)
		}
		vv[newp.LangIndex()] = newp
		return vv, oldp, oldp != nil
	default:
		panic(fmt.Sprintf("unknown type %T", old))
	}
}

func newPageMap(i int, s *Site, mcache *dynacache.Cache, pageTrees *pageTrees) *pageMap {
	var m *pageMap

	var taxonomiesConfig taxonomiesConfig = s.conf.Taxonomies

	m = &pageMap{
		pageTrees:              pageTrees.Shape(0, i),
		cachePages1:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag1/%d", i), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cachePages2:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag2/%d", i), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cacheGetTerms:          dynacache.GetOrCreatePartition[string, map[string]page.Pages](mcache, fmt.Sprintf("/gett/%d", i), dynacache.OptionsPartition{Weight: 5, ClearWhen: dynacache.ClearOnRebuild}),
		cacheResources:         dynacache.GetOrCreatePartition[string, resource.Resources](mcache, fmt.Sprintf("/ress/%d", i), dynacache.OptionsPartition{Weight: 60, ClearWhen: dynacache.ClearOnRebuild}),
		cacheContentRendered:   dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentSummary]](mcache, fmt.Sprintf("/cont/ren/%d", i), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		cacheContentPlain:      dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentPlainPlainWords]](mcache, fmt.Sprintf("/cont/pla/%d", i), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		contentTableOfContents: dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentTableOfContents]](mcache, fmt.Sprintf("/cont/toc/%d", i), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),

		contentDataFileSeenItems: maps.NewCache[string, map[uint64]bool](),

		cfg: contentMapConfig{
			lang:                 s.Lang(),
			taxonomyConfig:       taxonomiesConfig.Values(),
			taxonomyDisabled:     !s.conf.IsKindEnabled(kinds.KindTaxonomy),
			taxonomyTermDisabled: !s.conf.IsKindEnabled(kinds.KindTerm),
			pageDisabled:         !s.conf.IsKindEnabled(kinds.KindPage),
		},
		i: i,
		s: s,
	}

	m.pageReverseIndex = &contentTreeReverseIndex{
		initFn: func(rm map[any]contentNodeI) {
			add := func(k string, n contentNodeI) {
				existing, found := rm[k]
				if found && existing != ambiguousContentNode {
					rm[k] = ambiguousContentNode
				} else if !found {
					rm[k] = n
				}
			}

			w := &doctree.NodeShiftTreeWalker[contentNodeI]{
				Tree:     m.treePages,
				LockType: doctree.LockTypeRead,
				Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
					p := n.(*pageState)
					if p.PathInfo() != nil {
						add(p.PathInfo().BaseNameNoIdentifier(), p)
					}
					return false, nil
				},
			}

			if err := w.Walk(context.Background()); err != nil {
				panic(err)
			}
		},
		contentTreeReverseIndexMap: &contentTreeReverseIndexMap{},
	}

	return m
}

type contentTreeReverseIndex struct {
	initFn func(rm map[any]contentNodeI)
	*contentTreeReverseIndexMap
}

func (c *contentTreeReverseIndex) Reset() {
	c.contentTreeReverseIndexMap = &contentTreeReverseIndexMap{
		m: make(map[any]contentNodeI),
	}
}

func (c *contentTreeReverseIndex) Get(key any) contentNodeI {
	c.init.Do(func() {
		c.m = make(map[any]contentNodeI)
		c.initFn(c.contentTreeReverseIndexMap.m)
	})
	return c.m[key]
}

type contentTreeReverseIndexMap struct {
	init sync.Once
	m    map[any]contentNodeI
}

type sitePagesAssembler struct {
	*Site
	assembleChanges *WhatChanged
	ctx             context.Context
}

func (m *pageMap) debugPrint(prefix string, maxLevel int, w io.Writer) {
	noshift := false
	var prevKey string

	pageWalker := &doctree.NodeShiftTreeWalker[contentNodeI]{
		NoShift:     noshift,
		Tree:        m.treePages,
		Prefix:      prefix,
		WalkContext: &doctree.WalkContext[contentNodeI]{},
	}

	resourceWalker := pageWalker.Extend()
	resourceWalker.Tree = m.treeResources

	pageWalker.Handle = func(keyPage string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
		level := strings.Count(keyPage, "/")
		if level > maxLevel {
			return false, nil
		}
		const indentStr = " "
		p := n.(*pageState)
		s := strings.TrimPrefix(keyPage, paths.CommonDirPath(prevKey, keyPage))
		lenIndent := len(keyPage) - len(s)
		fmt.Fprint(w, strings.Repeat(indentStr, lenIndent))
		info := fmt.Sprintf("%s lm: %s (%s)", s, p.Lastmod().Format("2006-01-02"), p.Kind())
		fmt.Fprintln(w, info)
		switch p.Kind() {
		case kinds.KindTerm:
			m.treeTaxonomyEntries.WalkPrefix(
				doctree.LockTypeNone,
				keyPage+"/",
				func(s string, n *weightedContentNode) (bool, error) {
					fmt.Fprint(w, strings.Repeat(indentStr, lenIndent+4))
					fmt.Fprintln(w, s)
					return false, nil
				},
			)
		}

		isBranch := n.isContentNodeBranch()
		prevKey = keyPage
		resourceWalker.Prefix = keyPage + "/"

		resourceWalker.Handle = func(ss string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if isBranch {
				ownerKey, _ := pageWalker.Tree.LongestPrefix(ss, true, nil)
				if ownerKey != keyPage {
					// Stop walking downwards, someone else owns this resource.
					pageWalker.SkipPrefix(ownerKey + "/")
					return false, nil
				}
			}
			fmt.Fprint(w, strings.Repeat(indentStr, lenIndent+8))
			fmt.Fprintln(w, ss+" (resource)")
			return false, nil
		}

		return false, resourceWalker.Walk(context.Background())
	}

	err := pageWalker.Walk(context.Background())
	if err != nil {
		panic(err)
	}
}

func (h *HugoSites) dynacacheGCFilenameIfNotWatchedAndDrainMatching(filename string) {
	cpss := h.BaseFs.ResolvePaths(filename)
	if len(cpss) == 0 {
		return
	}
	// Compile cache busters.
	var cacheBusters []func(string) bool
	for _, cps := range cpss {
		if cps.Watch {
			continue
		}
		np := glob.NormalizePath(path.Join(cps.Component, cps.Path))
		g, err := h.ResourceSpec.BuildConfig().MatchCacheBuster(h.Log, np)
		if err == nil && g != nil {
			cacheBusters = append(cacheBusters, g)
		}
	}
	if len(cacheBusters) == 0 {
		return
	}
	cacheBusterOr := func(s string) bool {
		for _, cb := range cacheBusters {
			if cb(s) {
				return true
			}
		}
		return false
	}

	h.dynacacheGCCacheBuster(cacheBusterOr)

	// We want to avoid that evicted items in the above is considered in the next step server change.
	_ = h.MemCache.DrainEvictedIdentitiesMatching(func(ki dynacache.KeyIdentity) bool {
		return cacheBusterOr(ki.Key.(string))
	})
}

func (h *HugoSites) dynacacheGCCacheBuster(cachebuster func(s string) bool) {
	if cachebuster == nil {
		return
	}
	shouldDelete := func(k, v any) bool {
		var b bool
		if s, ok := k.(string); ok {
			b = cachebuster(s)
		}

		return b
	}

	h.MemCache.ClearMatching(nil, shouldDelete)
}

func (h *HugoSites) resolveAndClearStateForIdentities(
	ctx context.Context,
	l logg.LevelLogger,
	cachebuster func(s string) bool, changes []identity.Identity,
) error {
	h.Log.Debug().Log(logg.StringFunc(
		func() string {
			var sb strings.Builder
			for _, change := range changes {
				var key string
				if kp, ok := change.(resource.Identifier); ok {
					key = " " + kp.Key()
				}
				sb.WriteString(fmt.Sprintf("Direct dependencies of %q (%T%s) =>\n", change.IdentifierBase(), change, key))
				seen := map[string]bool{
					change.IdentifierBase(): true,
				}
				// Print the top level dependencies.
				identity.WalkIdentitiesDeep(change, func(level int, id identity.Identity) bool {
					if level > 1 {
						return true
					}
					if !seen[id.IdentifierBase()] {
						sb.WriteString(fmt.Sprintf("         %s%s\n", strings.Repeat(" ", level), id.IdentifierBase()))
					}
					seen[id.IdentifierBase()] = true
					return false
				})
			}
			return sb.String()
		}),
	)

	for _, id := range changes {
		if staler, ok := id.(resource.Staler); ok {
			var msgDetail string
			if p, ok := id.(*pageState); ok && p.File() != nil {
				msgDetail = fmt.Sprintf(" (%s)", p.File().Filename())
			}
			h.Log.Trace(logg.StringFunc(func() string { return fmt.Sprintf("Marking stale: %s (%T)%s\n", id, id, msgDetail) }))
			staler.MarkStale()
		}
	}

	// The order matters here:
	// 1. Handle the cache busters first, as those may produce identities for the page reset step.
	// 2. Then reset the page outputs, which may mark some resources as stale.
	// 3. Then GC the cache.
	if cachebuster != nil {
		if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
			ll := l.WithField("substep", "gc dynacache cachebuster")
			h.dynacacheGCCacheBuster(cachebuster)
			return ll, nil
		}); err != nil {
			return err
		}
	}

	// Drain the cache eviction stack.
	evicted := h.Deps.MemCache.DrainEvictedIdentities()
	if len(evicted) < 200 {
		for _, c := range evicted {
			changes = append(changes, c.Identity)
		}
	} else {
		// Mass eviction, we might as well invalidate everything.
		changes = []identity.Identity{identity.GenghisKhan}
	}

	// Remove duplicates
	seen := make(map[identity.Identity]bool)
	var n int
	for _, id := range changes {
		if !seen[id] {
			seen[id] = true
			changes[n] = id
			n++
		}
	}
	changes = changes[:n]

	if h.pageTrees.treePagesFromTemplateAdapters.LenRaw() > 0 {
		if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
			ll := l.WithField("substep", "resolve content adapter change set").WithField("changes", len(changes))
			checkedCount := 0
			matchCount := 0
			depsFinder := identity.NewFinder(identity.FinderConfig{})

			h.pageTrees.treePagesFromTemplateAdapters.WalkPrefixRaw(doctree.LockTypeRead, "",
				func(s string, n *pagesfromdata.PagesFromTemplate) (bool, error) {
					for _, id := range changes {
						checkedCount++
						if r := depsFinder.Contains(id, n.DependencyManager, 2); r > identity.FinderNotFound {
							n.MarkStale()
							matchCount++
							break
						}
					}
					return false, nil
				})

			ll = ll.WithField("checked", checkedCount).WithField("matches", matchCount)
			return ll, nil
		}); err != nil {
			return err
		}
	}

	if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
		// changesLeft: The IDs that the pages is dependent on.
		// changesRight: The IDs that the pages depend on.
		ll := l.WithField("substep", "resolve page output change set").WithField("changes", len(changes))

		checkedCount, matchCount, err := h.resolveAndResetDependententPageOutputs(ctx, changes)
		ll = ll.WithField("checked", checkedCount).WithField("matches", matchCount)
		return ll, err
	}); err != nil {
		return err
	}

	if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
		ll := l.WithField("substep", "gc dynacache")

		h.MemCache.ClearOnRebuild(changes...)
		h.Log.Trace(logg.StringFunc(func() string {
			var sb strings.Builder
			sb.WriteString("dynacache keys:\n")
			for _, key := range h.MemCache.Keys(nil) {
				sb.WriteString(fmt.Sprintf("   %s\n", key))
			}
			return sb.String()
		}))
		return ll, nil
	}); err != nil {
		return err
	}

	return nil
}

// The left change set is the IDs that the pages is dependent on.
// The right change set is the IDs that the pages depend on.
func (h *HugoSites) resolveAndResetDependententPageOutputs(ctx context.Context, changes []identity.Identity) (int, int, error) {
	if changes == nil {
		return 0, 0, nil
	}

	// This can be shared (many of the same IDs are repeated).
	depsFinder := identity.NewFinder(identity.FinderConfig{})

	h.Log.Trace(logg.StringFunc(func() string {
		var sb strings.Builder
		sb.WriteString("resolve page dependencies: ")
		for _, id := range changes {
			sb.WriteString(fmt.Sprintf(" %T: %s|", id, id.IdentifierBase()))
		}
		return sb.String()
	}))

	var (
		resetCounter   atomic.Int64
		checkedCounter atomic.Int64
	)

	resetPo := func(po *pageOutput, r identity.FinderResult) {
		if po.pco != nil {
			po.pco.Reset() // Will invalidate content cache.
		}

		po.renderState = 0
		po.p.resourcesPublishInit = &sync.Once{}
		if r == identity.FinderFoundOneOfMany {
			// Will force a re-render even in fast render mode.
			po.renderOnce = false
		}
		resetCounter.Add(1)
		h.Log.Trace(logg.StringFunc(func() string {
			p := po.p
			return fmt.Sprintf("Resetting page output %s for %s for output %s\n", p.Kind(), p.Path(), po.f.Name)
		}))
	}

	// This can be a relativeley expensive operations, so we do it in parallel.
	g := rungroup.Run[*pageState](ctx, rungroup.Config[*pageState]{
		NumWorkers: h.numWorkers,
		Handle: func(ctx context.Context, p *pageState) error {
			if !p.isRenderedAny() {
				// This needs no reset, so no need to check it.
				return nil
			}
			// First check the top level dependency manager.
			for _, id := range changes {
				checkedCounter.Add(1)
				if r := depsFinder.Contains(id, p.dependencyManager, 2); r > identity.FinderFoundOneOfManyRepetition {
					for _, po := range p.pageOutputs {
						resetPo(po, r)
					}
					// Done.
					return nil
				}
			}
			// Then do a more fine grained reset for each output format.
		OUTPUTS:
			for _, po := range p.pageOutputs {
				if !po.isRendered() {
					continue
				}
				for _, id := range changes {
					checkedCounter.Add(1)
					if r := depsFinder.Contains(id, po.dependencyManagerOutput, 50); r > identity.FinderFoundOneOfManyRepetition {
						resetPo(po, r)
						continue OUTPUTS
					}
				}
			}
			return nil
		},
	})

	h.withPage(func(s string, p *pageState) bool {
		var needToCheck bool
		for _, po := range p.pageOutputs {
			if po.isRendered() {
				needToCheck = true
				break
			}
		}
		if needToCheck {
			g.Enqueue(p)
		}
		return false
	})

	err := g.Wait()
	resetCount := int(resetCounter.Load())
	checkedCount := int(checkedCounter.Load())

	return checkedCount, resetCount, err
}

// Calculate and apply aggregate values to the page tree (e.g. dates, cascades).
func (sa *sitePagesAssembler) applyAggregates() error {
	sectionPageCount := map[string]int{}

	pw := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:        sa.pageMap.treePages,
		LockType:    doctree.LockTypeRead,
		WalkContext: &doctree.WalkContext[contentNodeI]{},
	}
	rw := pw.Extend()
	rw.Tree = sa.pageMap.treeResources
	sa.lastmod = time.Time{}
	rebuild := sa.s.h.isRebuild()

	pw.Handle = func(keyPage string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
		pageBundle := n.(*pageState)

		if pageBundle.Kind() == kinds.KindTerm {
			// Delay this until they're created.
			return false, nil
		}

		if pageBundle.IsPage() {
			rootSection := pageBundle.Section()
			sectionPageCount[rootSection]++
		}

		// Handle cascades first to get any default dates set.
		var cascade map[page.PageMatcher]maps.Params
		if keyPage == "" {
			// Home page gets it's cascade from the site config.
			cascade = sa.conf.Cascade.Config

			if pageBundle.m.pageConfig.CascadeCompiled == nil {
				// Pass the site cascade downwards.
				pw.WalkContext.Data().Insert(keyPage, cascade)
			}
		} else {
			_, data := pw.WalkContext.Data().LongestPrefix(keyPage)
			if data != nil {
				cascade = data.(map[page.PageMatcher]maps.Params)
			}
		}

		if rebuild {
			if (pageBundle.IsHome() || pageBundle.IsSection()) && pageBundle.m.setMetaPostCount > 0 {
				oldDates := pageBundle.m.pageConfig.Dates

				// We need to wait until after the walk to determine if any of the dates have changed.
				pw.WalkContext.AddPostHook(
					func() error {
						if oldDates != pageBundle.m.pageConfig.Dates {
							sa.assembleChanges.Add(pageBundle)
						}
						return nil
					},
				)
			}
		}

		// Combine the cascade map with front matter.
		if err := pageBundle.setMetaPost(cascade); err != nil {
			return false, err
		}

		// We receive cascade values from above. If this leads to a change compared
		// to the previous value, we need to mark the page and its dependencies as changed.
		if rebuild && pageBundle.m.setMetaPostCascadeChanged {
			sa.assembleChanges.Add(pageBundle)
		}

		const eventName = "dates"
		if n.isContentNodeBranch() {
			if pageBundle.m.pageConfig.CascadeCompiled != nil {
				// Pass it down.
				pw.WalkContext.Data().Insert(keyPage, pageBundle.m.pageConfig.CascadeCompiled)
			}

			wasZeroDates := pageBundle.m.pageConfig.Dates.IsAllDatesZero()
			if wasZeroDates || pageBundle.IsHome() {
				pw.WalkContext.AddEventListener(eventName, keyPage, func(e *doctree.Event[contentNodeI]) {
					sp, ok := e.Source.(*pageState)
					if !ok {
						return
					}

					if wasZeroDates {
						pageBundle.m.pageConfig.Dates.UpdateDateAndLastmodAndPublishDateIfAfter(sp.m.pageConfig.Dates)
					}

					if pageBundle.IsHome() {
						if pageBundle.m.pageConfig.Dates.Lastmod.After(pageBundle.s.lastmod) {
							pageBundle.s.lastmod = pageBundle.m.pageConfig.Dates.Lastmod
						}
						if sp.m.pageConfig.Dates.Lastmod.After(pageBundle.s.lastmod) {
							pageBundle.s.lastmod = sp.m.pageConfig.Dates.Lastmod
						}
					}
				})
			}
		}

		// Send the date info up the tree.
		pw.WalkContext.SendEvent(&doctree.Event[contentNodeI]{Source: n, Path: keyPage, Name: eventName})

		isBranch := n.isContentNodeBranch()
		rw.Prefix = keyPage + "/"

		rw.Handle = func(resourceKey string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if isBranch {
				ownerKey, _ := pw.Tree.LongestPrefix(resourceKey, true, nil)
				if ownerKey != keyPage {
					// Stop walking downwards, someone else owns this resource.
					rw.SkipPrefix(ownerKey + "/")
					return false, nil
				}
			}
			rs := n.(*resourceSource)
			if rs.isPage() {
				pageResource := rs.r.(*pageState)
				relPath := pageResource.m.pathInfo.BaseRel(pageBundle.m.pathInfo)
				pageResource.m.resourcePath = relPath
				var cascade map[page.PageMatcher]maps.Params
				// Apply cascade (if set) to the page.
				_, data := pw.WalkContext.Data().LongestPrefix(resourceKey)
				if data != nil {
					cascade = data.(map[page.PageMatcher]maps.Params)
				}
				if err := pageResource.setMetaPost(cascade); err != nil {
					return false, err
				}
			}

			return false, nil
		}
		return false, rw.Walk(sa.ctx)
	}

	if err := pw.Walk(sa.ctx); err != nil {
		return err
	}

	if err := pw.WalkContext.HandleEventsAndHooks(); err != nil {
		return err
	}

	if !sa.s.conf.C.IsMainSectionsSet() {
		var mainSection string
		var maxcount int
		for section, counter := range sectionPageCount {
			if section != "" && counter > maxcount {
				mainSection = section
				maxcount = counter
			}
		}
		sa.s.conf.C.SetMainSections([]string{mainSection})

	}

	return nil
}

func (sa *sitePagesAssembler) applyAggregatesToTaxonomiesAndTerms() error {
	walkContext := &doctree.WalkContext[contentNodeI]{}

	handlePlural := func(key string) error {
		var pw *doctree.NodeShiftTreeWalker[contentNodeI]
		pw = &doctree.NodeShiftTreeWalker[contentNodeI]{
			Tree:        sa.pageMap.treePages,
			Prefix:      key, // We also want to include the root taxonomy nodes, so no trailing slash.
			LockType:    doctree.LockTypeRead,
			WalkContext: walkContext,
			Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
				p := n.(*pageState)
				if p.Kind() != kinds.KindTerm {
					// The other kinds were handled in applyAggregates.
					if p.m.pageConfig.CascadeCompiled != nil {
						// Pass it down.
						pw.WalkContext.Data().Insert(s, p.m.pageConfig.CascadeCompiled)
					}
				}

				if p.Kind() != kinds.KindTerm && p.Kind() != kinds.KindTaxonomy {
					// Already handled.
					return false, nil
				}

				const eventName = "dates"

				if p.Kind() == kinds.KindTerm {
					var cascade map[page.PageMatcher]maps.Params
					_, data := pw.WalkContext.Data().LongestPrefix(s)
					if data != nil {
						cascade = data.(map[page.PageMatcher]maps.Params)
					}
					if err := p.setMetaPost(cascade); err != nil {
						return false, err
					}
					if !p.s.shouldBuild(p) {
						sa.pageMap.treePages.Delete(s)
						sa.pageMap.treeTaxonomyEntries.DeletePrefix(paths.AddTrailingSlash(s))
					} else if err := sa.pageMap.treeTaxonomyEntries.WalkPrefix(
						doctree.LockTypeRead,
						paths.AddTrailingSlash(s),
						func(ss string, wn *weightedContentNode) (bool, error) {
							// Send the date info up the tree.
							pw.WalkContext.SendEvent(&doctree.Event[contentNodeI]{Source: wn.n, Path: ss, Name: eventName})
							return false, nil
						},
					); err != nil {
						return false, err
					}
				}

				// Send the date info up the tree.
				pw.WalkContext.SendEvent(&doctree.Event[contentNodeI]{Source: n, Path: s, Name: eventName})

				if p.m.pageConfig.Dates.IsAllDatesZero() {
					pw.WalkContext.AddEventListener(eventName, s, func(e *doctree.Event[contentNodeI]) {
						sp, ok := e.Source.(*pageState)
						if !ok {
							return
						}

						p.m.pageConfig.Dates.UpdateDateAndLastmodAndPublishDateIfAfter(sp.m.pageConfig.Dates)
					})
				}

				return false, nil
			},
		}

		if err := pw.Walk(sa.ctx); err != nil {
			return err
		}
		return nil
	}

	for _, viewName := range sa.pageMap.cfg.taxonomyConfig.views {
		if err := handlePlural(viewName.pluralTreeKey); err != nil {
			return err
		}
	}

	if err := walkContext.HandleEventsAndHooks(); err != nil {
		return err
	}

	return nil
}

func (sa *sitePagesAssembler) assembleTermsAndTranslations() error {
	var (
		pages   = sa.pageMap.treePages
		entries = sa.pageMap.treeTaxonomyEntries
		views   = sa.pageMap.cfg.taxonomyConfig.views
	)

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     pages,
		LockType: lockType,
		Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			ps := n.(*pageState)

			if ps.m.noLink() {
				return false, nil
			}

			if sa.pageMap.cfg.taxonomyTermDisabled {
				return false, nil
			}

			for _, viewName := range views {
				vals := types.ToStringSlicePreserveString(getParam(ps, viewName.plural, false))
				if vals == nil {
					continue
				}

				w := getParamToLower(ps, viewName.plural+"_weight")
				weight, err := cast.ToIntE(w)
				if err != nil {
					sa.Log.Warnf("Unable to convert taxonomy weight %#v to int for %q", w, n.Path())
					// weight will equal zero, so let the flow continue
				}

				for i, v := range vals {
					if v == "" {
						continue
					}
					viewTermKey := "/" + viewName.plural + "/" + v
					pi := sa.Site.Conf.PathParser().Parse(files.ComponentFolderContent, viewTermKey+"/_index.md")
					term := pages.Get(pi.Base())
					if term == nil {
						m := &pageMeta{
							term:     v,
							singular: viewName.singular,
							s:        sa.Site,
							pathInfo: pi,
							pageMetaParams: &pageMetaParams{
								pageConfig: &pagemeta.PageConfig{
									Kind: kinds.KindTerm,
								},
							},
						}
						n, pi, err := sa.h.newPage(m)
						if err != nil {
							return false, err
						}
						pages.InsertIntoValuesDimension(pi.Base(), n)
						term = pages.Get(pi.Base())
					} else {
						m := term.(*pageState).m
						m.term = v
						m.singular = viewName.singular
					}

					if s == "" {
						// Consider making this the real value.
						s = "/"
					}

					key := pi.Base() + s

					entries.Insert(key, &weightedContentNode{
						weight: weight,
						n:      n,
						term:   &pageWithOrdinal{pageState: term.(*pageState), ordinal: i},
					})
				}
			}
			return false, nil
		},
	}

	return w.Walk(sa.ctx)
}

func (sa *sitePagesAssembler) assembleResources() error {
	pagesTree := sa.pageMap.treePages
	resourcesTree := sa.pageMap.treeResources

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     pagesTree,
		LockType: lockType,
		Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			ps := n.(*pageState)

			// This is a little out of place, but is conveniently put here.
			// Check if translationKey is set by user.
			// This is to support the manual way of setting the translationKey in front matter.
			if ps.m.pageConfig.TranslationKey != "" {
				sa.s.h.translationKeyPages.Append(ps.m.pageConfig.TranslationKey, ps)
			}

			// Prepare resources for this page.
			ps.shiftToOutputFormat(true, 0)
			targetPaths := ps.targetPaths()
			baseTarget := targetPaths.SubResourceBaseTarget
			duplicateResourceFiles := true
			if ps.m.pageConfig.ContentMediaType.IsMarkdown() {
				duplicateResourceFiles = ps.s.ContentSpec.Converters.GetMarkupConfig().Goldmark.DuplicateResourceFiles
			}

			duplicateResourceFiles = duplicateResourceFiles || ps.s.Conf.IsMultihost()

			err := sa.pageMap.forEachResourceInPage(
				ps, lockType,
				!duplicateResourceFiles,
				func(resourceKey string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
					rs := n.(*resourceSource)
					if !match.Has(doctree.DimensionLanguage) {
						// We got an alternative language version.
						// Clone this and insert it into the tree.
						rs = rs.clone()
						resourcesTree.InsertIntoCurrentDimension(resourceKey, rs)
					}
					if rs.r != nil {
						return false, nil
					}

					relPathOriginal := rs.path.Unnormalized().PathRel(ps.m.pathInfo.Unnormalized())
					relPath := rs.path.BaseRel(ps.m.pathInfo)

					var targetBasePaths []string
					if ps.s.Conf.IsMultihost() {
						baseTarget = targetPaths.SubResourceBaseLink
						// In multihost we need to publish to the lang sub folder.
						targetBasePaths = []string{ps.s.GetTargetLanguageBasePath()} // TODO(bep) we don't need this as a slice anymore.

					}

					if rs.rc != nil && rs.rc.Content.IsResourceValue() {
						if rs.rc.Name == "" {
							rs.rc.Name = relPathOriginal
						}
						r, err := ps.m.s.ResourceSpec.NewResourceWrapperFromResourceConfig(rs.rc)
						if err != nil {
							return false, err
						}
						rs.r = r
						return false, nil
					}

					var mt media.Type
					if rs.rc != nil {
						mt = rs.rc.ContentMediaType
					}

					rd := resources.ResourceSourceDescriptor{
						OpenReadSeekCloser:   rs.opener,
						Path:                 rs.path,
						GroupIdentity:        rs.path,
						TargetPath:           relPathOriginal, // Use the original path for the target path, so the links can be guessed.
						TargetBasePaths:      targetBasePaths,
						BasePathRelPermalink: targetPaths.SubResourceBaseLink,
						BasePathTargetPath:   baseTarget,
						NameNormalized:       relPath,
						NameOriginal:         relPathOriginal,
						MediaType:            mt,
						LazyPublish:          !ps.m.pageConfig.Build.PublishResources,
					}

					if rs.rc != nil {
						rc := rs.rc
						rd.OpenReadSeekCloser = rc.Content.ValueAsOpenReadSeekCloser()
						if rc.Name != "" {
							rd.NameNormalized = rc.Name
							rd.NameOriginal = rc.Name
						}
						if rc.Title != "" {
							rd.Title = rc.Title
						}
						rd.Params = rc.Params
					}

					r, err := ps.m.s.ResourceSpec.NewResource(rd)
					if err != nil {
						return false, err
					}
					rs.r = r
					return false, nil
				},
			)

			return false, err
		},
	}

	return w.Walk(sa.ctx)
}

func (sa *sitePagesAssembler) assemblePagesStep1(ctx context.Context) error {
	if err := sa.addMissingTaxonomies(); err != nil {
		return err
	}
	if err := sa.addMissingRootSections(); err != nil {
		return err
	}
	if err := sa.addStandalonePages(); err != nil {
		return err
	}
	if err := sa.applyAggregates(); err != nil {
		return err
	}
	return nil
}

func (sa *sitePagesAssembler) assemblePagesStep2() error {
	if err := sa.removeShouldNotBuild(); err != nil {
		return err
	}
	if err := sa.assembleTermsAndTranslations(); err != nil {
		return err
	}
	if err := sa.applyAggregatesToTaxonomiesAndTerms(); err != nil {
		return err
	}

	return nil
}

func (sa *sitePagesAssembler) assemblePagesStepFinal() error {
	if err := sa.assembleResources(); err != nil {
		return err
	}
	return nil
}

// Remove any leftover node that we should not build for some reason (draft, expired, scheduled in the future).
// Note that for the home and section kinds we just disable the nodes to preserve the structure.
func (sa *sitePagesAssembler) removeShouldNotBuild() error {
	s := sa.Site
	var keys []string
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		LockType: doctree.LockTypeRead,
		Tree:     sa.pageMap.treePages,
		Handle: func(key string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			p := n.(*pageState)
			if !s.shouldBuild(p) {
				switch p.Kind() {
				case kinds.KindHome, kinds.KindSection, kinds.KindTaxonomy:
					// We need to keep these for the structure, but disable
					// them so they don't get listed/rendered.
					(&p.m.pageConfig.Build).Disable()
				default:
					keys = append(keys, key)
				}
			}
			return false, nil
		},
	}
	if err := w.Walk(sa.ctx); err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	sa.pageMap.DeletePageAndResourcesBelow(keys...)

	return nil
}

// // Create the fixed output pages, e.g. sitemap.xml, if not already there.
func (sa *sitePagesAssembler) addStandalonePages() error {
	s := sa.Site
	m := s.pageMap
	tree := m.treePages

	commit := tree.Lock(true)
	defer commit()

	addStandalone := func(key, kind string, f output.Format) {
		if !s.Conf.IsMultihost() {
			switch kind {
			case kinds.KindSitemapIndex, kinds.KindRobotsTXT:
				// Only one for all languages.
				if s.languagei != 0 {
					return
				}
			}
		}

		if !sa.Site.conf.IsKindEnabled(kind) || tree.Has(key) {
			return
		}

		m := &pageMeta{
			s:        s,
			pathInfo: s.Conf.PathParser().Parse(files.ComponentFolderContent, key+f.MediaType.FirstSuffix.FullSuffix),
			pageMetaParams: &pageMetaParams{
				pageConfig: &pagemeta.PageConfig{
					Kind: kind,
				},
			},
			standaloneOutputFormat: f,
		}

		p, _, _ := s.h.newPage(m)

		tree.InsertIntoValuesDimension(key, p)
	}

	addStandalone("/404", kinds.KindStatus404, output.HTTPStatusHTMLFormat)

	if s.conf.EnableRobotsTXT {
		if m.i == 0 || s.Conf.IsMultihost() {
			addStandalone("/_robots", kinds.KindRobotsTXT, output.RobotsTxtFormat)
		}
	}

	sitemapEnabled := false
	for _, s := range s.h.Sites {
		if s.conf.IsKindEnabled(kinds.KindSitemap) {
			sitemapEnabled = true
			break
		}
	}

	if sitemapEnabled {
		of := output.SitemapFormat
		if s.conf.Sitemap.Filename != "" {
			of.BaseName = paths.Filename(s.conf.Sitemap.Filename)
		}
		addStandalone("/_sitemap", kinds.KindSitemap, of)

		skipSitemapIndex := s.Conf.IsMultihost() || !(s.Conf.DefaultContentLanguageInSubdir() || s.Conf.IsMultilingual())
		if !skipSitemapIndex {
			of = output.SitemapIndexFormat
			if s.conf.Sitemap.Filename != "" {
				of.BaseName = paths.Filename(s.conf.Sitemap.Filename)
			}
			addStandalone("/_sitemapindex", kinds.KindSitemapIndex, of)
		}
	}

	return nil
}

func (sa *sitePagesAssembler) addMissingRootSections() error {
	var hasHome bool

	// Add missing root sections.
	seen := map[string]bool{}
	var w *doctree.NodeShiftTreeWalker[contentNodeI]
	w = &doctree.NodeShiftTreeWalker[contentNodeI]{
		LockType: doctree.LockTypeWrite,
		Tree:     sa.pageMap.treePages,
		Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if n == nil {
				panic("n is nil")
			}

			ps := n.(*pageState)

			if ps.Lang() != sa.Lang() {
				panic(fmt.Sprintf("lang mismatch: %q: %s != %s", s, ps.Lang(), sa.Lang()))
			}

			if s == "" {
				hasHome = true
				sa.home = ps
				return false, nil
			}

			switch ps.Kind() {
			case kinds.KindPage, kinds.KindSection:
				// OK
			default:
				// Skip taxonomy nodes etc.
				return false, nil
			}

			p := ps.m.pathInfo
			section := p.Section()
			if section == "" || seen[section] {
				return false, nil
			}
			seen[section] = true

			// Try to preserve the original casing if possible.
			sectionUnnormalized := p.Unnormalized().Section()
			pth := sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, "/"+sectionUnnormalized+"/_index.md")
			nn := w.Tree.Get(pth.Base())

			if nn == nil {
				m := &pageMeta{
					s:        sa.Site,
					pathInfo: pth,
				}

				ps, pth, err := sa.h.newPage(m)
				if err != nil {
					return false, err
				}
				w.Tree.InsertIntoValuesDimension(pth.Base(), ps)
			}

			// /a/b, we don't need to walk deeper.
			if strings.Count(s, "/") > 1 {
				w.SkipPrefix(s + "/")
			}

			return false, nil
		},
	}

	if err := w.Walk(sa.ctx); err != nil {
		return err
	}

	if !hasHome {
		p := sa.Site.Conf.PathParser().Parse(files.ComponentFolderContent, "/_index.md")
		m := &pageMeta{
			s:        sa.Site,
			pathInfo: p,
			pageMetaParams: &pageMetaParams{
				pageConfig: &pagemeta.PageConfig{
					Kind: kinds.KindHome,
				},
			},
		}
		n, p, err := sa.h.newPage(m)
		if err != nil {
			return err
		}
		w.Tree.InsertIntoValuesDimensionWithLock(p.Base(), n)
		sa.home = n
	}

	return nil
}

func (sa *sitePagesAssembler) addMissingTaxonomies() error {
	if sa.pageMap.cfg.taxonomyDisabled && sa.pageMap.cfg.taxonomyTermDisabled {
		return nil
	}

	tree := sa.pageMap.treePages

	commit := tree.Lock(true)
	defer commit()

	for _, viewName := range sa.pageMap.cfg.taxonomyConfig.views {
		key := viewName.pluralTreeKey
		if v := tree.Get(key); v == nil {
			m := &pageMeta{
				s:        sa.Site,
				pathInfo: sa.Conf.PathParser().Parse(files.ComponentFolderContent, key+"/_index.md"),
				pageMetaParams: &pageMetaParams{
					pageConfig: &pagemeta.PageConfig{
						Kind: kinds.KindTaxonomy,
					},
				},
				singular: viewName.singular,
			}
			p, _, _ := sa.h.newPage(m)
			tree.InsertIntoValuesDimension(key, p)
		}
	}

	return nil
}

func (m *pageMap) CreateSiteTaxonomies(ctx context.Context) error {
	m.s.taxonomies = make(page.TaxonomyList)

	if m.cfg.taxonomyDisabled && m.cfg.taxonomyTermDisabled {
		return nil
	}

	for _, viewName := range m.cfg.taxonomyConfig.views {
		key := viewName.pluralTreeKey
		m.s.taxonomies[viewName.plural] = make(page.Taxonomy)
		w := &doctree.NodeShiftTreeWalker[contentNodeI]{
			Tree:     m.treePages,
			Prefix:   paths.AddTrailingSlash(key),
			LockType: doctree.LockTypeRead,
			Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
				p := n.(*pageState)

				switch p.Kind() {
				case kinds.KindTerm:
					if !p.m.shouldList(true) {
						return false, nil
					}
					taxonomy := m.s.taxonomies[viewName.plural]
					if taxonomy == nil {
						return true, fmt.Errorf("missing taxonomy: %s", viewName.plural)
					}
					if p.m.term == "" {
						panic("term is empty")
					}
					k := strings.ToLower(p.m.term)

					err := m.treeTaxonomyEntries.WalkPrefix(
						doctree.LockTypeRead,
						paths.AddTrailingSlash(s),
						func(ss string, wn *weightedContentNode) (bool, error) {
							taxonomy[k] = append(taxonomy[k], page.NewWeightedPage(wn.weight, wn.n.(page.Page), wn.term.Page()))
							return false, nil
						},
					)
					if err != nil {
						return true, err
					}

				default:
					return false, nil
				}

				return false, nil
			},
		}

		if err := w.Walk(ctx); err != nil {
			return err
		}
	}

	for _, taxonomy := range m.s.taxonomies {
		for _, v := range taxonomy {
			v.Sort()
		}
	}

	return nil
}

type viewName struct {
	singular      string // e.g. "category"
	plural        string // e.g. "categories"
	pluralTreeKey string
}

func (v viewName) IsZero() bool {
	return v.singular == ""
}
