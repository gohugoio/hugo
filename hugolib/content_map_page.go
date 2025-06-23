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
	"iter"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/hugolib/sitematrix"
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

// Invoked on rebuilds.
func (m *pageMap) Reset() {
	m.pageReverseIndex.Reset()
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
	dims := sitematrix.Vector{} // TODO1 fix the below.
	tree.ForEeachInDimension(key, dims, sitematrix.Language.Index(),
		cb,
	)

	tree = t.treeResources
	nCount = 0
	tree.ForEeachInDimension(key, dims, sitematrix.Language.Index(),
		cb,
	)

	if p.Component() == files.ComponentFolderContent {
		// It may also be a bundled content resource.
		key := p.ForType(paths.TypeContentResource).Base()
		tree = t.treeResources
		nCount = 0
		tree.ForEeachInDimension(key, dims, sitematrix.Language.Index(),
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
		Handle: func(key string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
		Handle: func(key string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
			Tree:            m.treePages,
			Prefix:          prefix,
			DelegeeFallback: true,
		}

		w.Handle = func(key string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
	handle func(resourceKey string, n contentNodeI, match sitematrix.Dimension) (bool, error),
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

	rw.Handle = func(resourceKey string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
		if isBranch {
			// A resourceKey always represents a filename with extension.
			// A page key points to the logical path of a page, which when sourced from the filesystem
			// may represent a directory (bundles) or a single content file (e.g. p1.md).
			// So, to avoid any overlapping ambiguity, we start looking from the owning directory.
			s := resourceKey

			for {
				s = path.Dir(s)
				ownerKey, found := m.treePages.LongestPrefixAll(s)
				if !found {
					return true, nil
				}
				if ownerKey == keyPage {
					break
				}

				if s != ownerKey && strings.HasPrefix(s, ownerKey) {
					// Keep looking
					continue
				}

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
	m.forEachResourceInPage(ps, doctree.LockTypeNone, false, func(resourceKey string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
	contentWeight() int
	matchDirectOrInDelegees(sitematrix.Vector) (contentNodeI, sitematrix.Vector)
	Dims() sitematrix.VectorProvider // TODO1 Can we unexport this?
	buildStateReseter
	resource.StaleMarker
}

var (
	_ contentNodeI = (*contentNodeIs)(nil)
	_ contentNodeI = (*contentNodeIs2)(nil)
)

type contentNodeIs2 []contentNodeI // TODO1 remove this type.

func (n contentNodeIs2) Dims() sitematrix.VectorProvider {
	panic("Dims(): not supported")
}

func (n contentNodeIs2) Path() string {
	return n.first().Path()
}

func (n contentNodeIs2) isContentNodeBranch() bool {
	return n.first().isContentNodeBranch()
}

func (n contentNodeIs2) contentWeight() int {
	return 0
}

func (n contentNodeIs2) first() contentNodeI {
	if len(n) == 0 {
		panic("empty contentNodeIs2")
	}
	return n[0]
}

func (n contentNodeIs2) matchDirectOrInDelegees(q sitematrix.Vector) (contentNodeI, sitematrix.Vector) {
	panic("not supported")
}

func (n contentNodeIs2) GetIdentity() identity.Identity {
	return n.first().GetIdentity()
}

func (n contentNodeIs2) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, nn := range n {
		if nn != nil {
			if nn.ForEeachIdentity(f) {
				return true
			}
		}
	}
	return false
}

func (n contentNodeIs2) resetBuildState() {
	for _, nn := range n {
		if nn != nil {
			nn.resetBuildState()
		}
	}
}

func (n contentNodeIs2) MarkStale() {
	for _, nn := range n {
		resource.MarkStale(nn)
	}
}

// TODO1 remove me.
type contentNodeIs map[sitematrix.Vector]contentNodeI

func (n contentNodeIs) one() contentNodeI {
	for _, nn := range n {
		return nn
	}
	return nil
}

// TODO1 remove this from the contentNodeI interface.
func (n contentNodeIs) Dims() sitematrix.VectorProvider {
	panic("not supported")
}

func (n contentNodeIs) contentWeight() int {
	return 0
}

func (n contentNodeIs) Path() string {
	return n.one().Path()
}

func (n contentNodeIs) isContentNodeBranch() bool {
	return n.one().isContentNodeBranch()
}

func (p contentNodeIs) matchDirectOrInDelegees(sitematrix.Vector) (contentNodeI, sitematrix.Vector) {
	panic("not implemented")
}

func (n contentNodeIs) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
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
	numLanguages int                // TODO1 remove.
	conf         config.AllProvider // Used for logging/debugging.
}

func (s *contentNodeShifter) Delete(n contentNodeI, dims sitematrix.Vector) (contentNodeI, bool, bool) {
	switch v := n.(type) {
	case contentNodeIs:
		deleted := v[dims]
		resource.MarkStale(deleted)
		wasDeleted := deleted != nil
		v[dims] = nil
		isEmpty := true
		for _, vv := range v {
			if vv != nil {
				isEmpty = false
				break
			}
		}
		return deleted, wasDeleted, isEmpty
	case contentNodeIs2:
		panic("TODO1 Delete")
		// TODO1 implement me.
		return nil, false, false

	case resourceSources:
		deleted := v[dims]
		resource.MarkStale(deleted)
		wasDeleted := deleted != nil
		v[dims] = nil
		isEmpty := true
		for _, vv := range v {
			if vv != nil {
				isEmpty = false
				break
			}
		}
		return deleted, wasDeleted, isEmpty
	case *resourceSource:
		if dims != v.Dims() {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	case *pageState:
		if dims != v.s.dims {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	default:
		panic(fmt.Sprintf("Delete: unknown type %T", n))
	}
}

func (s *contentNodeShifter) findDelegee(q sitematrix.Vector, candidates iter.Seq[contentNodeI]) contentNodeI {
	var (
		best         contentNodeI = nil
		bestDistance int
	)
	for n := range candidates {
		// The order of candidates is unstable, so we need to compare the matches to
		// get stable output. This compare will also make sure that we pick
		// language, version and role according to their individual sort order:
		// Closer is better, and matches above are better than matches below.
		if nn, dims := n.matchDirectOrInDelegees(q); nn != nil {
			distance := q.Distance(dims)
			if best == nil {
				best = nn
				bestDistance = distance
			} else {
				distanceAbs := absint(distance)
				bestDistanceAbs := absint(bestDistance)
				if distanceAbs < bestDistanceAbs {
					// Closer is better.
					best = nn
					bestDistance = distance
				} else if distanceAbs == bestDistanceAbs && distance > 0 {
					// Positive distance is better than negative.
					best = nn
					bestDistance = distance
				}
			}
		}
	}
	return best
}

func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (s *contentNodeShifter) Shift(n contentNodeI, dims sitematrix.Vector, exact, delegeeFallback bool) (contentNodeI, bool, sitematrix.Dimension) {
	accuracy := sitematrix.Language
	switch v := n.(type) {
	case contentNodeIs:
		if len(v) == 0 {
			panic("empty contentNodeIs")
		}
		vv := v[dims]
		if vv != nil {
			return vv, true, accuracy
		}
		if !delegeeFallback {
			return nil, false, 0
		}
		iter := func(yield func(n contentNodeI) bool) {
			for _, nn := range v {
				if !yield(nn) {
					return
				}
			}
		}
		if vv = s.findDelegee(dims, iter); vv != nil {
			return vv, true, accuracy
		}
		return nil, false, 0
	case contentNodeIs2:
		panic("not supported") // TODO1 remove this type.
		for _, vv := range v {
			if vv.Dims().HasVector(dims) {
				return vv, true, sitematrix.Language
			}
		}

		if !delegeeFallback {
			return nil, false, 0
		}

		// TODO1 performance for iter.Seq.
		iter := func(yield func(n contentNodeI) bool) {
			for _, nn := range v {
				if !yield(nn) {
					return
				}
			}
		}
		if vvv := s.findDelegee(dims, iter); vvv != nil {
			return vvv, true, accuracy
		}
		return nil, false, 0
	case resourceSources: // TODO1 remove this type.
		vv := v[dims]
		if vv != nil {
			return vv, true, sitematrix.Language
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
	case resourceSourcesSlice:
		for _, vv := range v {
			if vv.Dims().HasVector(dims) {
				return vv, true, sitematrix.Language
			}
			if !exact && vv.isPage() {
				// For non content resources, pick the first match.
				return vv, true, 0
			}
		}

	case *resourceSource:
		// TODO1 think.
		if v.Dims().FirstVector()[sitematrix.Language.Index()] == dims[sitematrix.Language.Index()] {
			return v, true, sitematrix.Language // TODO1
		}

		if !v.isPage() && !exact {
			return v, true, 0
		}
	case *pageMeta:
		panic("TODO1 remove me") // TODO1 remove this type.
		if v.Dims().HasVector(dims) {
			return v, true, sitematrix.Language
		}
	case *pageState:
		// TODO1 think.
		if v.s.dims == dims {
			return n, true, sitematrix.Language
		}
	default:
		panic(fmt.Sprintf("Shift: unsupported type %T", n))
	}
	return nil, false, 0
}

func (s *contentNodeShifter) ForEeachInDimension(n contentNodeI, dims sitematrix.Vector, d int, f func(contentNodeI) bool) {
	switch vv := n.(type) {
	case contentNodeIs:
	LOOP1:
		for dims2, v := range vv {
			if v != nil {
				for i, v := range dims2 {
					if i != d && v != dims[i] {
						continue LOOP1
					}
				}
				if f(v) {
					return
				}
			}
		}

	case contentNodeIs2: //
		panic("not supported") // TODO1 remove this type.
		seen := make(map[sitematrix.Vector]bool)
		for _, v := range vv {
			v.Dims().ForEeachVector(func(dims2 sitematrix.Vector) bool {
				seen[dims2] = true
				var match bool
				for i, vv := range dims2 {
					if i != d && vv != dims[i] {
						match = false
						break
					}
					match = true
				}
				if match {
					if f(v) {
						return true
					}
				}
				return false
			})
		}
	default:
		if n == nil {
			return
		}

		n.Dims().ForEeachVector(func(dims2 sitematrix.Vector) bool {
			var match bool
			for i, v := range dims2 {
				if i != d && v != dims[i] {
					match = false
					break
				}
				match = true
			}
			if match {
				if f(n) {
					return true
				}
			}
			return false
		})

	}
}

func (s *contentNodeShifter) InsertInto(old, new contentNodeI, dimension sitematrix.Vector) (contentNodeI, contentNodeI, bool) {
	switch vv := old.(type) {
	case *pageState:
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("InsertInto: unknown type %T", new))
		}
		if vv.s.dims == newp.s.dims && newp.s.dims == dimension {
			return new, vv, true
		}
		is := make(contentNodeIs, s.numLanguages)
		is[vv.s.dims] = old
		is[dimension] = new
		return is, old, false
	case contentNodeIs:
		oldv := vv[dimension]
		vv[dimension] = new
		return vv, oldv, oldv != nil
	case resourceSources:
		oldv := vv[dimension]
		vv[dimension] = new.(*resourceSource)
		return vv, oldv, oldv != nil
	case *resourceSource:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.Dims() == newp.Dims() && newp.Dims() == dimension {
			return new, vv, true
		}
		rs := make(resourceSources)
		rs[vv.Dims().FirstVector()] = vv
		rs[dimension] = newp
		return rs, vv, false

	default:
		panic(fmt.Sprintf("InsertInto: unknown type %T", old))
	}
}

func (s *contentNodeShifter) Insert(old, new contentNodeI) (contentNodeI, contentNodeI, bool) {
	deb("Insert: old %T, new %T", old, new)
	switch vv := old.(type) {
	case *pageMetaSource:
		return pageMetaSourcesSlice{vv, new.(*pageMetaSource)}, old, false
	case *pageMeta:
		switch new := new.(type) {
		case *pageState:
			return new, old, true
		case *pageMeta:
			is := make(contentNodeIs)
			// TODO1 remove s from pageMeta.
			vv.dims.ForEeachVector(func(dims sitematrix.Vector) bool {
				if vvv, ok := is[dims]; ok && vvv.contentWeight() > vv.contentWeight() {
					return true
				}
				deb("1 Insert: inserting pageMeta %s with dims %v/%v\t%v", old.Path(), dims, vv.dims, s.conf.ConfiguredDimensions().ResolveNames(dims))
				is[dims] = vv
				return true
			})
			new.dims.ForEeachVector(func(dims sitematrix.Vector) bool {
				if vvv, ok := is[dims]; ok && vvv.contentWeight() > new.contentWeight() {
					return true
				}
				deb("2 Insert: inserting pageMeta %s with dims %v/%v %v", new.Path(), dims, new.dims, s.conf.ConfiguredDimensions().ResolveNames(dims))
				is[dims] = new
				return true
			})

			// TODO1 stale + updated.
			return is, old, false
		default:
			panic(fmt.Sprintf("Insert: unknown type %T", new))
		}

	case *pageState: // TODO1 remove?
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		if vv.s.dims == newp.s.dims {
			if newp != old {
				resource.MarkStale(old)
			}
			return new, vv, true
		}
		is := make(contentNodeIs)
		is[vv.s.dims] = old
		is[newp.s.dims] = new
		return is, old, false
	case contentNodeIs:
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		oldp := vv[newp.s.dims]
		if oldp != newp {
			resource.MarkStale(oldp)
		}
		vv[newp.s.dims] = new
		return vv, oldp, oldp != nil
	case contentNodeIs2:
		panic("not supported") // TODO1 remove this type.
		newp, ok := new.(*pageState)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}

		// TODO1.
		/*
			if vv.Dims().EqualsVector(newp.Dims()) {
				if vv != newp {
					resource.MarkStale(vv)
				}
				return new, vv, true
			}*/

		rs := append(vv, newp)

		return rs, vv, false

	case *resourceSource:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}

		if vv.Dims().EqualsVector(newp.Dims()) {
			if vv != newp {
				resource.MarkStale(vv)
			}
			return new, vv, true
		}

		rs := resourceSourcesSlice{
			vv,
			newp,
		}

		return rs, vv, false

	case resourceSources:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		oldp := vv[newp.Dims().FirstVector()]
		if oldp != newp {
			resource.MarkStale(oldp)
		}
		vv[newp.Dims().FirstVector()] = newp
		return vv, oldp, oldp != nil
	default:
		panic(fmt.Sprintf("Insert: unknown type %T", old))
	}
}

func newPageMap(sitei, versioni, rolei int, s *Site, mcache *dynacache.Cache, pageTrees *pageTrees) *pageMap {
	var m *pageMap

	roleVersionSite := fmt.Sprintf("s%d/%d&%d", rolei, versioni, sitei)

	var taxonomiesConfig taxonomiesConfig = s.conf.Taxonomies

	m = &pageMap{
		pageTrees:              pageTrees.Shape(sitematrix.Language.Index(), sitei).Shape(sitematrix.Version.Index(), versioni).Shape(sitematrix.Role.Index(), rolei),
		cachePages1:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag1/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cachePages2:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag2/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cacheGetTerms:          dynacache.GetOrCreatePartition[string, map[string]page.Pages](mcache, fmt.Sprintf("/gett/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 5, ClearWhen: dynacache.ClearOnRebuild}),
		cacheResources:         dynacache.GetOrCreatePartition[string, resource.Resources](mcache, fmt.Sprintf("/ress/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 60, ClearWhen: dynacache.ClearOnRebuild}),
		cacheContentRendered:   dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentSummary]](mcache, fmt.Sprintf("/cont/ren/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		cacheContentPlain:      dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentPlainPlainWords]](mcache, fmt.Sprintf("/cont/pla/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		contentTableOfContents: dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentTableOfContents]](mcache, fmt.Sprintf("/cont/toc/%s", roleVersionSite), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),

		contentDataFileSeenItems: maps.NewCache[string, map[uint64]bool](),

		cfg: contentMapConfig{
			lang:                 s.Lang(),
			taxonomyConfig:       taxonomiesConfig.Values(),
			taxonomyDisabled:     !s.conf.IsKindEnabled(kinds.KindTaxonomy),
			taxonomyTermDisabled: !s.conf.IsKindEnabled(kinds.KindTerm),
			pageDisabled:         !s.conf.IsKindEnabled(kinds.KindPage),
		},
		i: sitei,
		s: s,
	}

	m.pageReverseIndex = newContentTreeTreverseIndex(func(get func(key any) (contentNodeI, bool), set func(key any, val contentNodeI)) {
		add := func(k string, n contentNodeI) {
			existing, found := get(k)
			if found && existing != ambiguousContentNode {
				set(k, ambiguousContentNode)
			} else if !found {
				set(k, n)
			}
		}

		w := &doctree.NodeShiftTreeWalker[contentNodeI]{
			Tree:     m.treePages,
			LockType: doctree.LockTypeRead,
			Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
	})

	return m
}

func newContentTreeTreverseIndex(init func(get func(key any) (contentNodeI, bool), set func(key any, val contentNodeI))) *contentTreeReverseIndex {
	return &contentTreeReverseIndex{
		initFn: init,
		mm:     maps.NewCache[any, contentNodeI](),
	}
}

type contentTreeReverseIndex struct {
	initFn func(get func(key any) (contentNodeI, bool), set func(key any, val contentNodeI))
	mm     *maps.Cache[any, contentNodeI]
}

func (c *contentTreeReverseIndex) Reset() {
	c.mm.Reset()
}

func (c *contentTreeReverseIndex) Get(key any) contentNodeI {
	v, _ := c.mm.InitAndGet(key, func(get func(key any) (contentNodeI, bool), set func(key any, val contentNodeI)) error {
		c.initFn(get, set)
		return nil
	})
	return v
}

type sitePagesAssembler struct {
	s               *Site
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

	pageWalker.Handle = func(keyPage string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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

		resourceWalker.Handle = func(ss string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
			if isBranch {
				ownerKey, _ := pageWalker.Tree.LongestPrefix(ss, true, false, nil)
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
	// Drain the cache eviction stack to start fresh.
	evictedStart := h.Deps.MemCache.DrainEvictedIdentities()

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
	// 1. Then GC the cache, which may produce changes.
	// 2. Then reset the page outputs, which may mark some resources as stale.
	if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
		ll := l.WithField("substep", "gc dynacache")

		predicate := func(k any, v any) bool {
			if cachebuster != nil {
				if s, ok := k.(string); ok {
					return cachebuster(s)
				}
			}
			return false
		}

		h.MemCache.ClearOnRebuild(predicate, changes...)
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

	// Drain the cache eviction stack.
	evicted := h.Deps.MemCache.DrainEvictedIdentities()
	if len(evicted) < 200 {
		for _, c := range evicted {
			changes = append(changes, c.Identity)
		}

		if len(evictedStart) > 0 {
			// In low memory situations and/or very big sites, there can be a lot of unrelated evicted items,
			// but there's a chance that some of them are related to the changes we are about to process,
			// so check.
			depsFinder := identity.NewFinder(identity.FinderConfig{})
			var addends []identity.Identity
			for _, ev := range evictedStart {
				for _, id := range changes {
					if cachebuster != nil && cachebuster(ev.Key.(string)) {
						addends = append(addends, ev.Identity)
						break
					}
					if r := depsFinder.Contains(id, ev.Identity, -1); r > 0 {
						addends = append(addends, ev.Identity)
						break
					}
				}
			}
			changes = append(changes, addends...)
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

	resetPo := func(po *pageOutput, rebuildContent bool, r identity.FinderResult) {
		if rebuildContent && po.pco != nil {
			po.pco.Reset() // Will invalidate content cache.
		}

		po.renderState = 0
		if r == identity.FinderFoundOneOfMany || po.f.Name == output.HTTPStatus404HTMLFormat.Name {
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
	g := rungroup.Run(ctx, rungroup.Config[*pageState]{
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
						// Note that p.dependencyManager is used when rendering content, so reset that.
						resetPo(po, true, r)
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
						// Note that dependencyManagerOutput is not used when rendering content, so don't reset that.
						resetPo(po, false, r)
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
		Tree:        sa.s.pageMap.treePages,
		LockType:    doctree.LockTypeRead,
		WalkContext: &doctree.WalkContext[contentNodeI]{},
	}
	rw := pw.Extend()
	rw.Tree = sa.s.pageMap.treeResources
	sa.s.lastmod = time.Time{}
	rebuild := sa.s.h.isRebuild()

	pw.Handle = func(keyPage string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
		var cascade *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig]
		if keyPage == "" {
			// Home page gets it's cascade from the site config.
			cascade = sa.s.conf.Cascade.Config
			if pageBundle.m.pageConfig.CascadeCompiled == nil {
				// Pass the site cascade downwards.
				pw.WalkContext.Data().Insert(keyPage, cascade)
			}
		} else {
			_, data := pw.WalkContext.Data().LongestPrefix(paths.Dir(keyPage))
			if data != nil {
				cascade = data.(*maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig])
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

		rw.Handle = func(resourceKey string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
			if isBranch {
				ownerKey, _ := pw.Tree.LongestPrefix(resourceKey, true, false, nil)
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
				var cascade *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig]
				// Apply cascade (if set) to the page.
				_, data := pw.WalkContext.Data().LongestPrefix(resourceKey)
				if data != nil {
					cascade = data.(*maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig])
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
			Tree:        sa.s.pageMap.treePages,
			Prefix:      key, // We also want to include the root taxonomy nodes, so no trailing slash.
			LockType:    doctree.LockTypeRead,
			WalkContext: walkContext,
			Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
					var cascade *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig]
					_, data := pw.WalkContext.Data().LongestPrefix(s)
					if data != nil {
						cascade = data.(*maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig])
					}
					if err := p.setMetaPost(cascade); err != nil {
						return false, err
					}
					if !p.s.shouldBuild(p) {
						sa.s.pageMap.treePages.Delete(s)
						sa.s.pageMap.treeTaxonomyEntries.DeletePrefix(paths.AddTrailingSlash(s))
					} else if err := sa.s.pageMap.treeTaxonomyEntries.WalkPrefix(
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

	for _, viewName := range sa.s.pageMap.cfg.taxonomyConfig.views {
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
	if sa.s.pageMap.cfg.taxonomyTermDisabled {
		return nil
	}

	var (
		pages   = sa.s.pageMap.treePages
		entries = sa.s.pageMap.treeTaxonomyEntries
		views   = sa.s.pageMap.cfg.taxonomyConfig.views
	)

	rebuild := sa.s.h.isRebuild()

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     pages,
		LockType: lockType,
		Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
			ps := n.(*pageState)

			if ps.m.noLink() {
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
					sa.s.Log.Warnf("Unable to convert taxonomy weight %#v to int for %q", w, n.Path())
					// weight will equal zero, so let the flow continue
				}

				for i, v := range vals {
					if v == "" {
						continue
					}
					viewTermKey := "/" + viewName.plural + "/" + v
					pi := sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, viewTermKey+"/_index.md")
					term := pages.Get(pi.Base())
					if term == nil {
						if rebuild {
							// A new tag was added in server mode.
							taxonomy := pages.Get(viewName.pluralTreeKey)
							if taxonomy != nil {
								sa.assembleChanges.Add(taxonomy.GetIdentity())
							}
						}

						m := &pageMeta{
							pageMetaSource: &pageMetaSource{
								pathInfo: pi,
							},
							term:     v,
							singular: viewName.singular,
							// TODO1 s:        sa.s,
							pathInfo: pi,
							pageMetaParams: &pageMetaParams{
								pageConfig: &pagemeta.PageConfig{
									PageConfigEarly: pagemeta.PageConfigEarly{
										Kind: kinds.KindTerm,
									},
								},
							},
						}
						ps, err := sa.s.newPageNew(m)
						if err != nil {
							return false, err
						}
						pages.InsertIntoValuesDimension(ps.PathInfo().Base(), ps)
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
	pagesTree := sa.s.pageMap.treePages
	resourcesTree := sa.s.pageMap.treeResources

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree:     pagesTree,
		LockType: lockType,
		Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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

			err := sa.s.pageMap.forEachResourceInPage(
				ps, lockType,
				!duplicateResourceFiles,
				func(resourceKey string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
					rs := n.(*resourceSource)
					if !match.Has(sitematrix.Language) {
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
						r, err := ps.s.ResourceSpec.NewResourceWrapperFromResourceConfig(rs.rc)
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

					var filename string
					if rs.fi != nil {
						filename = rs.fi.Meta().Filename
					}

					rd := resources.ResourceSourceDescriptor{
						OpenReadSeekCloser:   rs.opener,
						Path:                 rs.path,
						GroupIdentity:        rs.path,
						TargetPath:           relPathOriginal, // Use the original path for the target path, so the links can be guessed.
						TargetBasePaths:      targetBasePaths,
						BasePathRelPermalink: targetPaths.SubResourceBaseLink,
						BasePathTargetPath:   baseTarget,
						SourceFilenameOrPath: filename,
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

					r, err := ps.s.ResourceSpec.NewResource(rd)
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

func (sa *sitePagesAssembler) assemblePagesStep1() error {
	defer herrors.Recover()

	// TODO1

	/*if err := sa.addMissingTaxonomies(); err != nil {
		return err
	}*/
	/*if err := sa.addMissingRootSections(); err != nil {
		return err
	}*/
	/*
		if err := sa.addStandalonePages(); err != nil {
			return err
		}
	*/

	if err := sa.addMissingRootSections(); err != nil { // TODO1 see above.
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
	s := sa.s
	var keys []string
	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		LockType: doctree.LockTypeRead,
		Tree:     sa.s.pageMap.treePages,
		Handle: func(key string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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

	sa.s.pageMap.DeletePageAndResourcesBelow(keys...)

	return nil
}

// // Create the fixed output pages, e.g. sitemap.xml, if not already there.
func (sa *sitePagesAssembler) addStandalonePages() error {
	s := sa.s
	m := s.pageMap
	tree := m.treePages

	commit := tree.Lock(true)
	defer commit()

	addStandalone := func(key, kind string, f output.Format) {
		if !s.Conf.IsMultihost() {
			switch kind {
			case kinds.KindSitemapIndex, kinds.KindRobotsTXT:
				// Only one for all languages.
				if s.dims.Language() != 0 {
					return
				}
			}
		}

		if !sa.s.conf.IsKindEnabled(kind) || tree.Has(key) {
			return
		}

		m := &pageMeta{
			// TODO1 s:        s,
			pathInfo: s.Conf.PathParser().Parse(files.ComponentFolderContent, key+f.MediaType.FirstSuffix.FullSuffix),
			pageMetaParams: &pageMetaParams{
				pageConfig: &pagemeta.PageConfig{
					PageConfigEarly: pagemeta.PageConfigEarly{
						Kind: kind,
					},
				},
			},
			standaloneOutputFormat: f,
		}

		p, _, _ := s.h.newPage(m)

		tree.InsertIntoValuesDimension(key, p)
	}

	addStandalone("/404", kinds.KindStatus404, output.HTTPStatus404HTMLFormat)

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
		Tree:     sa.s.pageMap.treePages,
		Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
			if n == nil {
				panic("n is nil")
			}

			ps := n.(*pageState)

			if s == "" {
				hasHome = true
				sa.s.home = ps
				return false, nil
			}

			// TODO1
			if true {
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
					// TODO1 s:        sa.s,
					pathInfo: pth,
				}

				ps, pth, err := sa.s.h.newPage(m)
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

	// TODO1
	if false && !hasHome {
		p := sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, "/_index.md")
		m := &pageMeta{
			// TODO1 s:        sa.s,
			pathInfo: p,
			pageMetaParams: &pageMetaParams{
				pageConfig: &pagemeta.PageConfig{
					PageConfigEarly: pagemeta.PageConfigEarly{
						Kind: kinds.KindHome,
					},
				},
			},
		}
		n, p, err := sa.s.h.newPage(m)
		if err != nil {
			return err
		}
		w.Tree.InsertIntoValuesDimensionWithLock(p.Base(), n)
		sa.s.home = n
	}

	return nil
}

func (sa *sitePagesAssembler) createPages() error {
	sites := sa.s.h.sitesVersionsRolesMap

	var w *doctree.NodeShiftTreeWalker[contentNodeI]
	w = &doctree.NodeShiftTreeWalker[contentNodeI]{
		LockType: doctree.LockTypeWrite,
		Tree:     sa.s.pageMap.treePages,
		NoShift:  true,

		Transform: func(s string, n contentNodeI) (contentNodeI, bool, bool, error) {
			handlePageMetaSource := func(ms *pageMetaSource, is contentNodeIs) error {
				var err error
				d := ms.f.FileInfo().Meta().SiteInts
				if d == nil {
					panic(fmt.Sprintf("pageMetaSource %s has no site dimension", ms.f.FileInfo().Meta().Filename))
				}
				d.ForEeachVector(func(vec sitematrix.Vector) bool {
					site, found := sites[vec]
					if !found {
						panic(fmt.Sprintf("site not found for %v", vec))
					}
					var p *pageState
					p, err = site.newPageFromPageMetasource(ms)
					if err != nil {
						return false
					}

					if pp, found := is[vec]; found && pp.contentWeight() > p.contentWeight() {
						return true
					}

					is[vec] = p
					return true
				})
				return err
			}

			switch v := n.(type) {
			case pageMetaSourcesSlice:
				is := make(contentNodeIs)
				for _, ms := range v {
					if err := handlePageMetaSource(ms, is); err != nil {
						return nil, false, false, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
					}
				}
				return is, true, false, nil
			case *pageMetaSource:
				is := make(contentNodeIs)
				if err := handlePageMetaSource(v, is); err != nil {
					return nil, false, false, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
				}
				return is, true, false, nil
			case *pageMeta: // TODO1 remove.
				site, found := sites[v.dims.FirstVector()]
				if !found {
					panic(fmt.Sprintf("site not found for %v", v))
				}
				p, err := site.newPageNew(v)
				return p, true, false, err
			case contentNodeIs:
				for i, vv := range v {
					if m, ok := vv.(*pageMeta); ok {
						var err error
						site, found := sites[m.dims.FirstVector()] // TODO1 get rid of this interface.
						if !found {
							panic(fmt.Sprintf("site not found for %v", m))
						}
						v[i], err = site.newPageNew(m)
						if err != nil {
							return nil, false, false, fmt.Errorf("failed to create page %s: %w", s, err)
						}
					}
				}
			default:
				panic(fmt.Sprintf("unexpected contentNodeI type %T for %s", n, s))
			}

			return n, false, false, nil
		},
		Handle: nil,
	}

	if err := w.Walk(sa.ctx); err != nil {
		return err
	}

	return nil
}

func (sa *sitePagesAssembler) addMissingTaxonomies() error {
	if sa.s.pageMap.cfg.taxonomyDisabled && sa.s.pageMap.cfg.taxonomyTermDisabled {
		return nil
	}

	tree := sa.s.pageMap.treePages

	commit := tree.Lock(true)
	defer commit()

	for _, viewName := range sa.s.pageMap.cfg.taxonomyConfig.views {
		key := viewName.pluralTreeKey
		if v := tree.Get(key); v == nil {
			m := &pageMeta{
				// TODO1 s:        sa.s,
				pathInfo: sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, key+"/_index.md"),
				pageMetaParams: &pageMetaParams{
					pageConfig: &pagemeta.PageConfig{
						PageConfigEarly: pagemeta.PageConfigEarly{
							Kind: kinds.KindTaxonomy,
						},
					},
				},
				singular: viewName.singular,
			}
			p, _, _ := sa.s.h.newPage(m)
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
			Handle: func(s string, n contentNodeI, match sitematrix.Dimension) (bool, error) {
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
