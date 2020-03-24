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
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/para"
	"github.com/pkg/errors"
)

func newPageMaps(h *HugoSites) *pageMaps {
	mps := make([]*pageMap, len(h.Sites))
	for i, s := range h.Sites {
		mps[i] = s.pageMap
	}
	return &pageMaps{
		workers: para.New(h.numWorkers),
		pmaps:   mps,
	}

}

type pageMap struct {
	s *Site
	*contentMap
}

func (m *pageMap) Len() int {
	l := 0
	for _, t := range m.contentMap.pageTrees {
		l += t.Len()
	}
	return l
}

func (m *pageMap) createMissingTaxonomyNodes() error {
	if m.cfg.taxonomyDisabled {
		return nil
	}
	m.taxonomyEntries.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)
		vi := n.viewInfo
		k := cleanTreeKey(vi.name.plural + "/" + vi.termKey)

		if _, found := m.taxonomies.Get(k); !found {
			vic := &contentBundleViewInfo{
				name:       vi.name,
				termKey:    vi.termKey,
				termOrigin: vi.termOrigin,
			}
			m.taxonomies.Insert(k, &contentNode{viewInfo: vic})
		}
		return false
	})

	return nil
}

func (m *pageMap) newPageFromContentNode(n *contentNode, parentBucket *pagesMapBucket, owner *pageState) (*pageState, error) {
	if n.fi == nil {
		panic("FileInfo must (currently) be set")
	}

	f, err := newFileInfo(m.s.SourceSpec, n.fi)
	if err != nil {
		return nil, err
	}

	meta := n.fi.Meta()
	content := func() (hugio.ReadSeekCloser, error) {
		return meta.Open()
	}

	bundled := owner != nil
	s := m.s

	sections := s.sectionsFromFile(f)

	kind := s.kindFromFileInfoOrSections(f, sections)
	if kind == page.KindTaxonomy {
		s.PathSpec.MakePathsSanitized(sections)
	}

	metaProvider := &pageMeta{kind: kind, sections: sections, bundled: bundled, s: s, f: f}

	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	if n.fi.Meta().GetBool(walkIsRootFileMetaKey) {
		// Make sure that the bundle/section we start walking from is always
		// rendered.
		// This is only relevant in server fast render mode.
		ps.forceRender = true
	}

	n.p = ps
	if ps.IsNode() {
		ps.bucket = newPageBucket(ps)
	}

	gi, err := s.h.gitInfoForPage(ps)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Git data")
	}
	ps.gitInfo = gi

	r, err := content()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	parseResult, err := pageparser.Parse(
		r,
		pageparser.Config{EnableEmoji: s.siteCfg.enableEmoji},
	)
	if err != nil {
		return nil, err
	}

	ps.pageContent = pageContent{
		source: rawPageContent{
			parsed:         parseResult,
			posMainContent: -1,
			posSummaryEnd:  -1,
			posBodyStart:   -1,
		},
	}

	ps.shortcodeState = newShortcodeHandler(ps, ps.s, nil)

	if err := ps.mapContent(parentBucket, metaProvider); err != nil {
		return nil, ps.wrapError(err)
	}

	if err := metaProvider.applyDefaultValues(n); err != nil {
		return nil, err
	}

	ps.init.Add(func() (interface{}, error) {
		pp, err := newPagePaths(s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		outputFormatsForPage := ps.m.outputFormats()

		// Prepare output formats for all sites.
		// We do this even if this page does not get rendered on
		// its own. It may be referenced via .Site.GetPage and
		// it will then need an output format.
		ps.pageOutputs = make([]*pageOutput, len(ps.s.h.renderFormats))
		created := make(map[string]*pageOutput)
		shouldRenderPage := !ps.m.noRender()

		for i, f := range ps.s.h.renderFormats {
			if po, found := created[f.Name]; found {
				ps.pageOutputs[i] = po
				continue
			}

			render := shouldRenderPage
			if render {
				_, render = outputFormatsForPage.GetByName(f.Name)
			}

			po := newPageOutput(ps, pp, f, render)

			// Create a content provider for the first,
			// we may be able to reuse it.
			if i == 0 {
				contentProvider, err := newPageContentOutput(ps, po)
				if err != nil {
					return nil, err
				}
				po.initContentProvider(contentProvider)
			}

			ps.pageOutputs[i] = po
			created[f.Name] = po

		}

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil
	})

	ps.parent = owner

	return ps, nil
}

func (m *pageMap) newResource(fim hugofs.FileMetaInfo, owner *pageState) (resource.Resource, error) {

	if owner == nil {
		panic("owner is nil")
	}
	// TODO(bep) consolidate with multihost logic + clean up
	outputFormats := owner.m.outputFormats()
	seen := make(map[string]bool)
	var targetBasePaths []string
	// Make sure bundled resources are published to all of the ouptput formats'
	// sub paths.
	for _, f := range outputFormats {
		p := f.Path
		if seen[p] {
			continue
		}
		seen[p] = true
		targetBasePaths = append(targetBasePaths, p)

	}

	meta := fim.Meta()
	r := func() (hugio.ReadSeekCloser, error) {
		return meta.Open()
	}

	target := strings.TrimPrefix(meta.Path(), owner.File().Dir())

	return owner.s.ResourceSpec.New(
		resources.ResourceSourceDescriptor{
			TargetPaths:        owner.getTargetPaths,
			OpenReadSeekCloser: r,
			FileInfo:           fim,
			RelTargetFilename:  target,
			TargetBasePaths:    targetBasePaths,
			LazyPublish:        !owner.m.buildConfig.PublishResources,
		})
}

func (m *pageMap) createSiteTaxonomies() error {
	m.s.taxonomies = make(TaxonomyList)
	m.taxonomies.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)
		t := n.viewInfo

		viewName := t.name

		if t.termKey == "" {
			m.s.taxonomies[viewName.plural] = make(Taxonomy)
		} else {
			taxonomy := m.s.taxonomies[viewName.plural]
			m.taxonomyEntries.WalkPrefix(s+"/", func(ss string, v interface{}) bool {
				b2 := v.(*contentNode)
				info := b2.viewInfo
				taxonomy.add(info.termKey, page.NewWeightedPage(info.weight, info.ref.p, n.p))

				return false
			})
		}

		return false
	})

	for _, taxonomy := range m.s.taxonomies {
		for _, v := range taxonomy {
			v.Sort()
		}
	}

	return nil
}

func (m *pageMap) createListAllPages() page.Pages {
	pages := make(page.Pages, 0)

	m.contentMap.pageTrees.Walk(func(s string, n *contentNode) bool {
		if n.p == nil {
			panic(fmt.Sprintf("BUG: page not set for %q", s))
		}
		if contentTreeNoListAlwaysFilter(s, n) {
			return false
		}
		pages = append(pages, n.p)
		return false
	})

	page.SortByDefault(pages)
	return pages
}

func (m *pageMap) assemblePages() error {
	m.taxonomyEntries.DeletePrefix("/")

	if err := m.assembleSections(); err != nil {
		return err
	}

	var err error

	if err != nil {
		return err
	}

	m.pages.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)

		var shouldBuild bool

		defer func() {
			// Make sure we always rebuild the view cache.
			if shouldBuild && err == nil && n.p != nil {
				m.attachPageToViews(s, n)
			}
		}()

		if n.p != nil {
			// A rebuild
			shouldBuild = true
			return false
		}

		var parent *contentNode
		var parentBucket *pagesMapBucket

		_, parent = m.getSection(s)
		if parent == nil {
			panic(fmt.Sprintf("BUG: parent not set for %q", s))
		}
		parentBucket = parent.p.bucket

		n.p, err = m.newPageFromContentNode(n, parentBucket, nil)
		if err != nil {
			return true
		}

		shouldBuild = !(n.p.Kind() == page.KindPage && m.cfg.pageDisabled) && m.s.shouldBuild(n.p)
		if !shouldBuild {
			m.deletePage(s)
			return false
		}

		n.p.treeRef = &contentTreeRef{
			m:   m,
			t:   m.pages,
			n:   n,
			key: s,
		}

		if err = m.assembleResources(s, n.p, parentBucket); err != nil {
			return true
		}

		return false
	})

	m.deleteOrphanSections()

	return err
}

func (m *pageMap) assembleResources(s string, p *pageState, parentBucket *pagesMapBucket) error {
	var err error

	m.resources.WalkPrefix(s, func(s string, v interface{}) bool {
		n := v.(*contentNode)
		meta := n.fi.Meta()
		classifier := meta.Classifier()
		var r resource.Resource
		switch classifier {
		case files.ContentClassContent:
			var rp *pageState
			rp, err = m.newPageFromContentNode(n, parentBucket, p)
			if err != nil {
				return true
			}
			rp.m.resourcePath = filepath.ToSlash(strings.TrimPrefix(rp.Path(), p.File().Dir()))
			r = rp

		case files.ContentClassFile:
			r, err = m.newResource(n.fi, p)
			if err != nil {
				return true
			}
		default:
			panic(fmt.Sprintf("invalid classifier: %q", classifier))
		}

		p.resources = append(p.resources, r)
		return false
	})

	return err
}

func (m *pageMap) assembleSections() error {

	var sectionsToDelete []string
	var err error

	m.sections.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)

		var shouldBuild bool

		defer func() {
			// Make sure we always rebuild the view cache.
			if shouldBuild && err == nil && n.p != nil {
				m.attachPageToViews(s, n)
				if n.p.IsHome() {
					m.s.home = n.p
				}
			}
		}()

		sections := m.splitKey(s)

		if n.p != nil {
			if n.p.IsHome() {
				m.s.home = n.p
			}
			shouldBuild = true
			return false
		}

		var parent *contentNode
		var parentBucket *pagesMapBucket

		if s != "/" {
			_, parent = m.getSection(s)
			if parent == nil || parent.p == nil {
				panic(fmt.Sprintf("BUG: parent not set for %q", s))
			}
		}

		if parent != nil {
			parentBucket = parent.p.bucket
		}

		kind := page.KindSection
		if s == "/" {
			kind = page.KindHome
		}

		if n.fi != nil {
			n.p, err = m.newPageFromContentNode(n, parentBucket, nil)
			if err != nil {
				return true
			}
		} else {
			n.p = m.s.newPage(n, parentBucket, kind, "", sections...)
		}

		shouldBuild = m.s.shouldBuild(n.p)
		if !shouldBuild {
			sectionsToDelete = append(sectionsToDelete, s)
			return false
		}

		n.p.treeRef = &contentTreeRef{
			m:   m,
			t:   m.sections,
			n:   n,
			key: s,
		}

		if err = m.assembleResources(s+cmLeafSeparator, n.p, parentBucket); err != nil {
			return true
		}

		return false
	})

	for _, s := range sectionsToDelete {
		m.deleteSectionByPath(s)
	}

	return err
}

func (m *pageMap) assembleTaxonomies() error {

	var taxonomiesToDelete []string
	var err error

	m.taxonomies.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)

		if n.p != nil {
			return false
		}

		kind := n.viewInfo.kind()
		sections := n.viewInfo.sections()

		_, parent := m.getTaxonomyParent(s)
		if parent == nil || parent.p == nil {
			panic(fmt.Sprintf("BUG: parent not set for %q", s))
		}
		parentBucket := parent.p.bucket

		if n.fi != nil {
			n.p, err = m.newPageFromContentNode(n, parent.p.bucket, nil)
			if err != nil {
				return true
			}
		} else {
			title := ""
			if kind == page.KindTaxonomy {
				title = n.viewInfo.term()
			}
			n.p = m.s.newPage(n, parent.p.bucket, kind, title, sections...)
		}

		if !m.s.shouldBuild(n.p) {
			taxonomiesToDelete = append(taxonomiesToDelete, s)
			return false
		}

		n.p.treeRef = &contentTreeRef{
			m:   m,
			t:   m.taxonomies,
			n:   n,
			key: s,
		}

		if err = m.assembleResources(s+cmLeafSeparator, n.p, parentBucket); err != nil {
			return true
		}

		return false
	})

	for _, s := range taxonomiesToDelete {
		m.deleteTaxonomy(s)
	}

	return err

}

func (m *pageMap) attachPageToViews(s string, b *contentNode) {
	if m.cfg.taxonomyDisabled {
		return
	}

	for _, viewName := range m.cfg.taxonomyConfig {
		vals := types.ToStringSlicePreserveString(getParam(b.p, viewName.plural, false))
		if vals == nil {
			continue
		}

		w := getParamToLower(b.p, viewName.plural+"_weight")
		weight, err := cast.ToIntE(w)
		if err != nil {
			m.s.Log.ERROR.Printf("Unable to convert taxonomy weight %#v to int for %q", w, b.p.Path())
			// weight will equal zero, so let the flow continue
		}

		for _, v := range vals {
			termKey := m.s.getTaxonomyKey(v)

			bv := &contentNode{
				viewInfo: &contentBundleViewInfo{
					name:       viewName,
					termKey:    termKey,
					termOrigin: v,
					weight:     weight,
					ref:        b,
				},
			}

			if s == "/" {
				// To avoid getting an empty key.
				s = page.KindHome
			}
			key := cleanTreeKey(path.Join(viewName.plural, termKey, s))
			m.taxonomyEntries.Insert(key, bv)
		}
	}
}

type pageMapQuery struct {
	Prefix string
	Filter contentTreeNodeCallback
}

func (m *pageMap) collectPages(query pageMapQuery, fn func(c *contentNode)) error {
	if query.Filter == nil {
		query.Filter = contentTreeNoListAlwaysFilter
	}

	m.pages.WalkQuery(query, func(s string, n *contentNode) bool {
		fn(n)
		return false
	})

	return nil
}

func (m *pageMap) collectPagesAndSections(query pageMapQuery, fn func(c *contentNode)) error {
	if err := m.collectSections(query, fn); err != nil {
		return err
	}

	query.Prefix = query.Prefix + cmBranchSeparator
	if err := m.collectPages(query, fn); err != nil {
		return err
	}

	return nil
}

func (m *pageMap) collectSections(query pageMapQuery, fn func(c *contentNode)) error {
	var level int
	isHome := query.Prefix == "/"

	if !isHome {
		level = strings.Count(query.Prefix, "/")
	}

	return m.collectSectionsFn(query, func(s string, c *contentNode) bool {
		if s == query.Prefix {
			return false
		}

		if (strings.Count(s, "/") - level) != 1 {
			return false
		}

		fn(c)

		return false
	})
}

func (m *pageMap) collectSectionsFn(query pageMapQuery, fn func(s string, c *contentNode) bool) error {

	if !strings.HasSuffix(query.Prefix, "/") {
		query.Prefix += "/"
	}

	m.sections.WalkQuery(query, func(s string, n *contentNode) bool {
		return fn(s, n)
	})

	return nil
}

func (m *pageMap) collectSectionsRecursiveIncludingSelf(query pageMapQuery, fn func(c *contentNode)) error {
	return m.collectSectionsFn(query, func(s string, c *contentNode) bool {
		fn(c)
		return false
	})
}

func (m *pageMap) collectTaxonomies(prefix string, fn func(c *contentNode)) error {
	m.taxonomies.WalkQuery(pageMapQuery{Prefix: prefix}, func(s string, n *contentNode) bool {
		fn(n)
		return false
	})
	return nil
}

// withEveryBundlePage applies fn to every Page, including those bundled inside
// leaf bundles.
func (m *pageMap) withEveryBundlePage(fn func(p *pageState) bool) {
	m.bundleTrees.Walk(func(s string, n *contentNode) bool {
		if n.p != nil {
			return fn(n.p)
		}
		return false
	})
}

type pageMaps struct {
	workers *para.Workers
	pmaps   []*pageMap
}

// deleteSection deletes the entire section from s.
func (m *pageMaps) deleteSection(s string) {
	m.withMaps(func(pm *pageMap) error {
		pm.deleteSectionByPath(s)
		return nil
	})
}

func (m *pageMaps) AssemblePages() error {
	return m.withMaps(func(pm *pageMap) error {
		if err := pm.CreateMissingNodes(); err != nil {
			return err
		}

		if err := pm.assemblePages(); err != nil {
			return err
		}

		if err := pm.createMissingTaxonomyNodes(); err != nil {
			return err
		}

		// Handle any new sections created in the step above.
		if err := pm.assembleSections(); err != nil {
			return err
		}

		if pm.s.home == nil {
			// Home is disabled, everything is.
			pm.bundleTrees.DeletePrefix("")
			return nil
		}

		if err := pm.assembleTaxonomies(); err != nil {
			return err
		}

		if err := pm.createSiteTaxonomies(); err != nil {
			return err
		}

		a := (&sectionWalker{m: pm.contentMap}).applyAggregates()
		_, mainSectionsSet := pm.s.s.Info.Params()["mainsections"]
		if !mainSectionsSet && a.mainSection != "" {
			mainSections := []string{a.mainSection}
			pm.s.s.Info.Params()["mainSections"] = mainSections
			pm.s.s.Info.Params()["mainsections"] = mainSections
		}

		pm.s.lastmod = a.datesAll.Lastmod()
		if resource.IsZeroDates(pm.s.home) {
			pm.s.home.m.Dates = a.datesAll
		}

		return nil
	})
}

func (m *pageMaps) walkBundles(fn func(n *contentNode) bool) {
	_ = m.withMaps(func(pm *pageMap) error {
		pm.bundleTrees.Walk(func(s string, n *contentNode) bool {
			return fn(n)
		})
		return nil
	})
}

func (m *pageMaps) walkBranchesPrefix(prefix string, fn func(s string, n *contentNode) bool) {
	_ = m.withMaps(func(pm *pageMap) error {
		pm.branchTrees.WalkPrefix(prefix, func(s string, n *contentNode) bool {
			return fn(s, n)
		})
		return nil
	})
}

func (m *pageMaps) withMaps(fn func(pm *pageMap) error) error {
	g, _ := m.workers.Start(context.Background())
	for _, pm := range m.pmaps {
		pm := pm
		g.Run(func() error {
			return fn(pm)
		})
	}
	return g.Wait()
}

type pagesMapBucket struct {
	// Cascading front matter.
	cascade maps.Params

	owner *pageState // The branch node

	*pagesMapBucketPages
}

type pagesMapBucketPages struct {
	pagesInit sync.Once
	pages     page.Pages

	pagesAndSectionsInit sync.Once
	pagesAndSections     page.Pages

	sectionsInit sync.Once
	sections     page.Pages
}

func (b *pagesMapBucket) getPages() page.Pages {
	b.pagesInit.Do(func() {
		b.pages = b.owner.treeRef.getPages()
		page.SortByDefault(b.pages)
	})
	return b.pages
}

func (b *pagesMapBucket) getPagesRecursive() page.Pages {
	pages := b.owner.treeRef.getPagesRecursive()
	page.SortByDefault(pages)
	return pages
}

func (b *pagesMapBucket) getPagesAndSections() page.Pages {
	b.pagesAndSectionsInit.Do(func() {
		b.pagesAndSections = b.owner.treeRef.getPagesAndSections()
	})
	return b.pagesAndSections
}

func (b *pagesMapBucket) getSections() page.Pages {
	b.sectionsInit.Do(func() {
		if b.owner.treeRef == nil {
			return
		}
		b.sections = b.owner.treeRef.getSections()
	})

	return b.sections
}

func (b *pagesMapBucket) getTaxonomies() page.Pages {
	b.sectionsInit.Do(func() {
		var pas page.Pages
		ref := b.owner.treeRef
		ref.m.collectTaxonomies(ref.key+"/", func(c *contentNode) {
			pas = append(pas, c.p)
		})
		page.SortByDefault(pas)
		b.sections = pas
	})

	return b.sections
}

func (b *pagesMapBucket) getTaxonomyEntries() page.Pages {
	var pas page.Pages
	ref := b.owner.treeRef
	viewInfo := ref.n.viewInfo
	prefix := strings.ToLower("/" + viewInfo.name.plural + "/" + viewInfo.termKey + "/")
	ref.m.taxonomyEntries.WalkPrefix(prefix, func(s string, v interface{}) bool {
		n := v.(*contentNode)
		pas = append(pas, n.viewInfo.ref.p)
		return false
	})
	page.SortByDefault(pas)
	return pas
}

type sectionAggregate struct {
	datesAll             resource.Dates
	datesSection         resource.Dates
	pageCount            int
	mainSection          string
	mainSectionPageCount int
}

type sectionAggregateHandler struct {
	sectionAggregate
	sectionPageCount int

	// Section
	b *contentNode
	s string
}

func (h *sectionAggregateHandler) isRootSection() bool {
	return h.s != "/" && strings.Count(h.s, "/") == 1
}

func (h *sectionAggregateHandler) handleNested(v sectionWalkHandler) error {
	nested := v.(*sectionAggregateHandler)
	h.sectionPageCount += nested.pageCount
	h.pageCount += h.sectionPageCount
	h.datesAll.UpdateDateAndLastmodIfAfter(nested.datesAll)
	h.datesSection.UpdateDateAndLastmodIfAfter(nested.datesAll)
	return nil
}

func (h *sectionAggregateHandler) handlePage(s string, n *contentNode) error {
	h.sectionPageCount++

	var d resource.Dated
	if n.p != nil {
		d = n.p
	} else if n.viewInfo != nil && n.viewInfo.ref != nil {
		d = n.viewInfo.ref.p
	} else {
		return nil
	}

	h.datesAll.UpdateDateAndLastmodIfAfter(d)
	h.datesSection.UpdateDateAndLastmodIfAfter(d)
	return nil
}

func (h *sectionAggregateHandler) handleSectionPost() error {
	if h.sectionPageCount > h.mainSectionPageCount && h.isRootSection() {
		h.mainSectionPageCount = h.sectionPageCount
		h.mainSection = strings.TrimPrefix(h.s, "/")
	}

	if resource.IsZeroDates(h.b.p) {
		h.b.p.m.Dates = h.datesSection
	}

	h.datesSection = resource.Dates{}

	return nil
}

func (h *sectionAggregateHandler) handleSectionPre(s string, b *contentNode) error {
	h.s = s
	h.b = b
	h.sectionPageCount = 0
	h.datesAll.UpdateDateAndLastmodIfAfter(b.p)
	return nil
}

type sectionWalkHandler interface {
	handleNested(v sectionWalkHandler) error
	handlePage(s string, b *contentNode) error
	handleSectionPost() error
	handleSectionPre(s string, b *contentNode) error
}

type sectionWalker struct {
	err error
	m   *contentMap
}

func (w *sectionWalker) applyAggregates() *sectionAggregateHandler {
	return w.walkLevel("/", func() sectionWalkHandler {
		return &sectionAggregateHandler{}
	}).(*sectionAggregateHandler)

}

func (w *sectionWalker) walkLevel(prefix string, createVisitor func() sectionWalkHandler) sectionWalkHandler {

	level := strings.Count(prefix, "/")
	visitor := createVisitor()

	w.m.taxonomies.WalkPrefix(prefix, func(s string, v interface{}) bool {
		currentLevel := strings.Count(s, "/")
		if currentLevel > level {
			return false
		}

		n := v.(*contentNode)

		if w.err = visitor.handleSectionPre(s, n); w.err != nil {
			return true
		}

		if currentLevel == 1 {
			nested := w.walkLevel(s+"/", createVisitor)
			if w.err = visitor.handleNested(nested); w.err != nil {
				return true
			}
		} else {
			w.m.taxonomyEntries.WalkPrefix(s, func(ss string, v interface{}) bool {
				n := v.(*contentNode)
				w.err = visitor.handlePage(ss, n)
				return w.err != nil
			})
		}

		w.err = visitor.handleSectionPost()

		return w.err != nil
	})

	w.m.sections.WalkPrefix(prefix, func(s string, v interface{}) bool {
		currentLevel := strings.Count(s, "/")
		if currentLevel > level {
			return false
		}

		n := v.(*contentNode)

		if w.err = visitor.handleSectionPre(s, n); w.err != nil {
			return true
		}

		w.m.pages.WalkPrefix(s+cmBranchSeparator, func(s string, v interface{}) bool {
			w.err = visitor.handlePage(s, v.(*contentNode))
			return w.err != nil
		})

		if w.err != nil {
			return true
		}

		if s != "/" {
			nested := w.walkLevel(s+"/", createVisitor)
			if w.err = visitor.handleNested(nested); w.err != nil {
				return true
			}
		}

		w.err = visitor.handleSectionPost()

		return w.err != nil
	})

	return visitor

}

type viewName struct {
	singular string // e.g. "category"
	plural   string // e.g. "categories"
}

func (v viewName) IsZero() bool {
	return v.singular == ""
}
