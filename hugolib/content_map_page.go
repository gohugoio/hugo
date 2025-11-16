// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"sync/atomic"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugofs/hglob"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var pagePredicates = struct {
	KindPage         predicate.PR[*pageState]
	KindSection      predicate.PR[*pageState]
	KindHome         predicate.PR[*pageState]
	KindTerm         predicate.PR[*pageState]
	ShouldListLocal  predicate.PR[*pageState]
	ShouldListGlobal predicate.PR[*pageState]
	ShouldListAny    predicate.PR[*pageState]
	ShouldLink       predicate.PR[page.Page]
}{
	KindPage: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.Kind() == kinds.KindPage)
	},
	KindSection: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.Kind() == kinds.KindSection)
	},
	KindHome: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.Kind() == kinds.KindHome)
	},
	KindTerm: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.Kind() == kinds.KindTerm)
	},
	ShouldListLocal: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.m.shouldList(false))
	},
	ShouldListGlobal: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.m.shouldList(true))
	},
	ShouldListAny: func(p *pageState) predicate.Match {
		return predicate.BoolMatch(p.m.shouldListAny())
	},
	ShouldLink: func(p page.Page) predicate.Match {
		return predicate.BoolMatch(!p.(*pageState).m.noLink())
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
	treePages *doctree.NodeShiftTree[contentNode]

	// This tree contains Resources bundled in pages.
	treeResources *doctree.NodeShiftTree[contentNode]

	// All pages and resources.
	treePagesResources doctree.WalkableTrees[contentNode]

	// This tree contains all taxonomy entries, e.g "/tags/blue/page1"
	treeTaxonomyEntries *doctree.TreeShiftTreeSlice[*weightedContentNode]

	// Stores the state for _content.gotmpl files.
	// Mostly releveant for rebuilds.
	treePagesFromTemplateAdapters *doctree.TreeShiftTreeSlice[*pagesfromdata.PagesFromTemplate]

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
	cb := func(n contentNode) bool {
		if n == nil {
			return false
		}
		cnh.markStale(n)
		if nCount > 0 {
			return true
		}
		nCount++

		cnh.toForEachIdentityProvider(n).ForEeachIdentity(func(id identity.Identity) bool {
			ids = append(ids, id)
			return false
		})

		return false
	}
	tree := t.treePages
	nCount = 0
	tree.ForEeachInAllDimensions(key, cb)

	tree = t.treeResources
	nCount = 0
	tree.ForEeachInAllDimensions(key, cb)

	if p.Component() == files.ComponentFolderContent {
		// It may also be a bundled content resource.
		key := p.ForType(paths.TypeContentResource).Base()
		tree = t.treeResources
		nCount = 0
		tree.ForEeachInAllDimensions(key, cb)

	}
	return ids
}

// collectIdentitiesSurrounding collects all identities surrounding the given key.
func (t *pageTrees) collectIdentitiesSurrounding(key string, maxSamplesPerTree int) []identity.Identity {
	ids := t.collectIdentitiesSurroundingIn(key, maxSamplesPerTree, t.treePages)
	ids = append(ids, t.collectIdentitiesSurroundingIn(key, maxSamplesPerTree, t.treeResources)...)
	return ids
}

func (t *pageTrees) collectIdentitiesSurroundingIn(key string, maxSamples int, tree *doctree.NodeShiftTree[contentNode]) []identity.Identity {
	var ids []identity.Identity
	section, ok := tree.LongestPrefixRaw(path.Dir(key))
	if ok {
		count := 0
		prefix := section + "/"
		level := strings.Count(prefix, "/")
		tree.WalkPrefixRaw(prefix, func(s string, n contentNode) bool {
			if level != strings.Count(s, "/") {
				return false
			}
			cnh.toForEachIdentityProvider(n).ForEeachIdentity(func(id identity.Identity) bool {
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

func (t pageTrees) Shape(v sitesmatrix.Vector) *pageTrees {
	t.treePages = t.treePages.Shape(v)
	t.treeResources = t.treeResources.Shape(v)
	t.treeTaxonomyEntries = t.treeTaxonomyEntries.Shape(v)
	t.treePagesFromTemplateAdapters = t.treePagesFromTemplateAdapters.Shape(v)
	t.createMutableTrees()

	return &t
}

func (t *pageTrees) createMutableTrees() {
	t.treePagesResources = doctree.WalkableTrees[contentNode]{
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
func (m *pageMap) forEachPage(include predicate.PR[*pageState], fn func(p *pageState) (bool, error)) error {
	if include == nil {
		include = func(p *pageState) predicate.Match {
			return predicate.True
		}
	}
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     m.treePages,
		LockType: doctree.LockTypeRead,
		Handle: func(key string, n contentNode) (bool, error) {
			if p, ok := n.(*pageState); ok && include(p).OK() {
				if terminate, err := fn(p); terminate || err != nil {
					return terminate, err
				}
			}
			return false, nil
		},
	}

	return w.Walk(context.Background())
}

func (m *pageMap) forEeachPageIncludingBundledPages(include predicate.PR[*pageState], fn func(p *pageState) (bool, error)) error {
	if include == nil {
		include = func(p *pageState) predicate.Match {
			return predicate.True
		}
	}

	if err := m.forEachPage(include, fn); err != nil {
		return err
	}

	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     m.treeResources,
		LockType: doctree.LockTypeRead,
		Handle: func(key string, n contentNode) (bool, error) {
			if p, ok := n.(*pageState); ok && include(p).OK() {
				if terminate, err := fn(p); terminate || err != nil {
					return terminate, err
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
			include = pagePredicates.ShouldListLocal.BoolFunc()
		}

		w := &doctree.NodeShiftTreeWalker[contentNode]{
			Tree:     m.treePages,
			Prefix:   prefix,
			Fallback: true,
		}

		w.Handle = func(key string, n contentNode) (bool, error) {
			if q.Recursive {
				if p, ok := n.(*pageState); ok && include(p) {
					pas = append(pas, p)
				}
				return false, nil
			}

			if p, ok := n.(*pageState); ok && include(p) {
				pas = append(pas, p)
			}

			if cnh.isBranchNode(n) {
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
			include = pagePredicates.ShouldListLocal.BoolFunc()
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
	fallback bool,
	transform func(resourceKey string, n contentNode) (n2 contentNode, ns doctree.NodeTransformState, err error),
	handle func(resourceKey string, n contentNode) (bool, error),
) error {
	keyPage := ps.Path()
	if keyPage == "/" {
		keyPage = ""
	}

	prefix := paths.AddTrailingSlash(ps.Path())

	isBranch := ps.IsNode()

	rwr := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     m.treeResources,
		Prefix:   prefix,
		LockType: lockType,
		Fallback: fallback,
	}

	shouldSkipOrTerminate := func(resourceKey string) (ns doctree.NodeTransformState) {
		if !isBranch {
			return
		}

		// A resourceKey always represents a filename with extension.
		// A page key points to the logical path of a page, which when sourced from the filesystem
		// may represent a directory (bundles) or a single content file (e.g. p1.md).
		// So, to avoid any overlapping ambiguity, we start looking from the owning directory.
		s := resourceKey

		for {
			s = path.Dir(s)
			ownerKey, found := m.treePages.LongestPrefixRaw(s)

			if !found {
				return doctree.NodeTransformStateTerminate
			}
			if ownerKey == keyPage {
				break
			}

			if s != ownerKey && strings.HasPrefix(s, ownerKey) {
				// Keep looking
				continue
			}

			// Stop walking downwards, someone else owns this resource.
			rwr.SkipPrefix(ownerKey + "/")
			return doctree.NodeTransformStateSkip
		}
		return
	}

	if transform != nil {
		rwr.Transform = func(resourceKey string, n contentNode) (n2 contentNode, ns doctree.NodeTransformState, err error) {
			if ns = shouldSkipOrTerminate(resourceKey); ns >= doctree.NodeTransformStateSkip {
				return
			}
			return transform(resourceKey, n)
		}
	}

	rwr.Handle = func(resourceKey string, n contentNode) (terminate bool, err error) {
		if transform == nil {
			if ns := shouldSkipOrTerminate(resourceKey); ns >= doctree.NodeTransformStateSkip {
				return ns == doctree.NodeTransformStateTerminate, nil
			}
		}
		return handle(resourceKey, n)
	}

	return rwr.Walk(context.Background())
}

func (m *pageMap) getResourcesForPage(ps *pageState) (resource.Resources, error) {
	var res resource.Resources
	m.forEachResourceInPage(ps, doctree.LockTypeNone, true, nil, func(resourceKey string, n contentNode) (bool, error) {
		switch n := n.(type) {
		case *resourceSource:
			r := n.r
			if r == nil {
				panic(fmt.Sprintf("getResourcesForPage: resource %q for page %q has no resource, sites matrix %v/%v", resourceKey, ps.Path(), ps.siteVector(), n.sv))
			}

			res = append(res, r)
		case *pageState:
			res = append(res, n)
		default:
			panic(fmt.Sprintf("getResourcesForPage: unknown type %T", n))
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

var _ doctree.Transformer[contentNode] = (*contentNodeTransformerRaw)(nil)

type contentNodeTransformerRaw struct{}

func (t *contentNodeTransformerRaw) Append(n contentNode, ns ...contentNode) (contentNode, bool) {
	if n == nil {
		if len(ns) == 1 {
			return ns[0], true
		}
		var ss contentNodes = ns
		return ss, true
	}

	switch v := n.(type) {
	case contentNodes:
		v = append(v, ns...)
		return v, true
	default:
		ss := make(contentNodes, 0, 1+len(ns))
		ss = append(ss, v)
		ss = append(ss, ns...)
		return ss, true
	}
}

func newPageMap(s *Site, mcache *dynacache.Cache, pageTrees *pageTrees) *pageMap {
	var m *pageMap

	vec := s.siteVector
	languageVersionRole := fmt.Sprintf("s%d/%d&%d", vec.Language(), vec.Version(), vec.Role())

	var taxonomiesConfig taxonomiesConfig = s.conf.Taxonomies

	m = &pageMap{
		pageTrees:              pageTrees.Shape(vec),
		cachePages1:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag1/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cachePages2:            dynacache.GetOrCreatePartition[string, page.Pages](mcache, fmt.Sprintf("/pag2/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild}),
		cacheGetTerms:          dynacache.GetOrCreatePartition[string, map[string]page.Pages](mcache, fmt.Sprintf("/gett/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 5, ClearWhen: dynacache.ClearOnRebuild}),
		cacheResources:         dynacache.GetOrCreatePartition[string, resource.Resources](mcache, fmt.Sprintf("/ress/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 60, ClearWhen: dynacache.ClearOnRebuild}),
		cacheContentRendered:   dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentSummary]](mcache, fmt.Sprintf("/cont/ren/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		cacheContentPlain:      dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentPlainPlainWords]](mcache, fmt.Sprintf("/cont/pla/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		contentTableOfContents: dynacache.GetOrCreatePartition[string, *resources.StaleValue[contentTableOfContents]](mcache, fmt.Sprintf("/cont/toc/%s", languageVersionRole), dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),

		contentDataFileSeenItems: maps.NewCache[string, map[uint64]bool](),

		cfg: contentMapConfig{
			lang:                 s.Lang(),
			taxonomyConfig:       taxonomiesConfig.Values(),
			taxonomyDisabled:     !s.conf.IsKindEnabled(kinds.KindTaxonomy),
			taxonomyTermDisabled: !s.conf.IsKindEnabled(kinds.KindTerm),
			pageDisabled:         !s.conf.IsKindEnabled(kinds.KindPage),
		},
		i: s.siteVector.Language(),
		s: s,
	}

	m.pageReverseIndex = newContentTreeTreverseIndex(func(get func(key any) (contentNode, bool), set func(key any, val contentNode)) {
		add := func(k string, n contentNode) {
			existing, found := get(k)
			if found && existing != ambiguousContentNode {
				set(k, ambiguousContentNode)
			} else if !found {
				set(k, n)
			}
		}

		w := &doctree.NodeShiftTreeWalker[contentNode]{
			Tree:     m.treePages,
			LockType: doctree.LockTypeRead,
			Handle: func(s string, n contentNode) (bool, error) {
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

func newContentTreeTreverseIndex(init func(get func(key any) (contentNode, bool), set func(key any, val contentNode))) *contentTreeReverseIndex {
	return &contentTreeReverseIndex{
		initFn: init,
		mm:     maps.NewCache[any, contentNode](),
	}
}

type contentTreeReverseIndex struct {
	initFn func(get func(key any) (contentNode, bool), set func(key any, val contentNode))
	mm     *maps.Cache[any, contentNode]
}

func (c *contentTreeReverseIndex) Reset() {
	c.mm.Reset()
}

func (c *contentTreeReverseIndex) Get(key any) contentNode {
	v, _ := c.mm.InitAndGet(key, func(get func(key any) (contentNode, bool), set func(key any, val contentNode)) error {
		c.initFn(get, set)
		return nil
	})
	return v
}

func (m *pageMap) debugPrint(prefix string, maxLevel int, w io.Writer) {
	noshift := false

	pageWalker := &doctree.NodeShiftTreeWalker[contentNode]{
		NoShift:     noshift,
		Tree:        m.treePages,
		Prefix:      prefix,
		WalkContext: &doctree.WalkContext[contentNode]{},
	}

	resourceWalker := pageWalker.Extend()
	resourceWalker.Tree = m.treeResources

	pageWalker.Handle = func(keyPage string, n contentNode) (bool, error) {
		level := strings.Count(keyPage, "/")
		if level > maxLevel {
			return false, nil
		}
		const indentStr = " "
		p := n.(*pageState)
		lenIndent := 0
		info := fmt.Sprintf("%s lm: %s (%s)", keyPage, p.Lastmod().Format("2006-01-02"), p.Kind())
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

		isBranch := cnh.isBranchNode(n)
		resourceWalker.Prefix = keyPage + "/"

		resourceWalker.Handle = func(ss string, n contentNode) (bool, error) {
			if isBranch {
				ownerKey, _ := pageWalker.Tree.LongestPrefix(ss, false, nil)
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
		np := hglob.NormalizePath(path.Join(cps.Component, cps.Path))
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
			return fmt.Sprintf("%s Resetting page output %q for %q for output %q\n", p.s.resolveDimensionNames(), p.Kind(), p.Path(), po.f.Name)
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

func (m *pageMap) CreateSiteTaxonomies(ctx context.Context) (page.TaxonomyList, error) {
	taxonomies := make(page.TaxonomyList)

	if m.cfg.taxonomyDisabled && m.cfg.taxonomyTermDisabled {
		return taxonomies, nil
	}

	for _, viewName := range m.cfg.taxonomyConfig.views {
		key := viewName.pluralTreeKey
		taxonomies[viewName.plural] = make(page.Taxonomy)
		w := &doctree.NodeShiftTreeWalker[contentNode]{
			Tree:     m.treePages,
			Prefix:   paths.AddTrailingSlash(key),
			LockType: doctree.LockTypeRead,
			Handle: func(s string, n contentNode) (bool, error) {
				p := n.(*pageState)

				switch p.Kind() {
				case kinds.KindTerm:
					if !p.m.shouldList(true) {
						return false, nil
					}
					taxonomy := taxonomies[viewName.plural]
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
			return nil, err
		}
	}

	for _, taxonomy := range taxonomies {
		for _, v := range taxonomy {
			v.Sort()
		}
	}

	return taxonomies, nil
}

type term struct {
	view viewName
	term string
}

type viewName struct {
	singular      string // e.g. "category"
	plural        string // e.g. "categories"
	pluralTreeKey string
}

func (v viewName) IsZero() bool {
	return v.singular == ""
}
