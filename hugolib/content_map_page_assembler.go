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
	"path"
	"slices"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/hdebug"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/spf13/cast"
)

// Handles the flow of creating all the pages and resources for all sites.
// This includes applying any cascade configuration passed down the tree.
type allPagesAssembler struct {
	// Dependencies.
	ctx context.Context
	h   *HugoSites
	m   *pageMap

	assembleChanges *WhatChanged

	// Internal state.
	pw *doctree.NodeShiftTreeWalker[contentNode] // walks pages.
	rw *doctree.NodeShiftTreeWalker[contentNode] // walks resources.
}

func newAllPagesAssembler(
	ctx context.Context,
	h *HugoSites,
	m *pageMap,
	assembleChanges *WhatChanged,
) *allPagesAssembler {
	rw := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:        m.treeResources,
		LockType:    doctree.LockTypeWrite,
		NoShift:     true,
		WalkContext: &doctree.WalkContext[contentNode]{},
	}
	pw := rw.Extend()
	pw.Tree = m.treePages

	return &allPagesAssembler{
		ctx:             ctx,
		h:               h,
		m:               m,
		assembleChanges: assembleChanges,

		pw: pw,
		rw: rw,
	}
}

type sitePagesAssembler struct {
	s               *Site
	assembleChanges *WhatChanged
	ctx             context.Context
}

func (a *allPagesAssembler) createAllPages() error {
	var (
		sites             = a.h.sitesVersionsRolesMap
		h                 = a.h
		isRebuild         = a.h.isRebuild()
		printPathWarnings = !isRebuild && a.h.Configs.Base.PrintPathWarnings
		lockType          = doctree.LockTypeWrite
		treePages         = a.m.treePages
		treeResources     = a.m.treeResources

		getViews = func(vec sitesmatrix.Vector) []viewName {
			return h.languageSiteForSiteVector(vec).pageMap.cfg.taxonomyConfig.views
		}
	)

	resourceOwnerInfo := struct {
		n contentNode
		s string
	}{}

	type term struct {
		view viewName
		term string
	}

	seenTerms := map[term]sitesmatrix.Vectors{}

	newHomePageMetaSource := func() *pageMetaSource {
		pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, "/_index.md")
		return &pageMetaSource{
			pathInfo:       pi,
			siteMatrixBase: a.h.Conf.AllSitesMatrix(),
			pageConfigSource: &pagemeta.PageConfig{
				PageConfigEarly: pagemeta.PageConfigEarly{
					Kind: kinds.KindHome,
				},
			},
		}
	}

	if err := a.createMissingPages(); err != nil {
		return err
	}

	if treePages.Len() == 0 {
		// No pages, insert a home page to get something to walk on.
		p := newHomePageMetaSource()
		treePages.InsertRawWithLock(p.pathInfo.Base(), p)
	}

	getCascades := func(s string) []page.PageMatcherParamsConfig {
		var cascades []page.PageMatcherParamsConfig
		data := a.rw.WalkContext.Data()
		if s == "" {
			// Home page gets it's cascade from the site config.
			// TODO1 make sure this is called for auto generated home pages too.
			for s := range a.h.allSiteLanguages(nil) {
				cascades = append(cascades, s.conf.Cascade.Config.Cascades...)
			}
		} else {
			_, data := data.LongestPrefix(paths.Dir(s))
			if data != nil {
				cascades = data.([]page.PageMatcherParamsConfig)
			}
		}
		return cascades
	}

	transformPages := func(s string, n contentNode) (n2 contentNode, replaced bool, skip bool, terminate bool, err error) {
		cascades := getCascades(s)
		cascadesLen := len(cascades)

		defer func() {
			if len(cascades) > cascadesLen || s == "" {
				// New cascade values added, pass them downwards.
				a.rw.WalkContext.Data().Insert(s, cascades)
			}
		}()

		handlePageMetaSource := func(v any, is contentNodes[contentNodePage], replaceVector bool) (bool, bool, error) {
			var (
				replaced bool
				err      error
			)
			switch ms := v.(type) {
			case *pageMetaSource:
				if err := ms.initSitesMatrix(a.h, cascades); err != nil {
					return false, false, err
				}

				if ms.isContentNodeBranch() && ms.pageConfigSource.CascadeCompiled != nil {
					// Cascade on itself has higher priority than inherited ones,
					// so insert it first.
					cascades = slices.Insert(cascades, 0, ms.pageConfigSource.CascadeCompiled...)
				}

				sitesMatrix := ms.sitesMatrix()

				sitesMatrix.ForEeachVector(func(vec sitesmatrix.Vector) bool {
					site, found := sites[vec]
					if !found {
						panic(fmt.Sprintf("site not found for %v", vec))
					}

					var p *pageState
					p, err = site.newPageFromPageMetasource(ms)
					if err != nil {
						return false
					}

					if s == "" {
						site.home = p
					}

					// Combine the cascade map with front matter.
					if err = p.setMetaPost(s, cascades); err != nil {
						return true
					}

					// We receive cascade values from above. If this leads to a change compared
					// to the previous value, we need to mark the page and its dependencies as changed.
					if isRebuild && p.m.setMetaPostCascadeChanged {
						a.assembleChanges.Add(p)
					}

					pp, found := is[vec]

					replaced = replaced || found

					if found && !replaceVector && pp.contentWeight() > p.contentWeight() {
						return true
					}

					is[vec] = p
					return true
				})
				return true, replaced, err
			case *pageState:
				is[ms.s.siteVector] = ms
				return true, replaced, err
			}
			return false, replaced, err
		}

		switch v := n.(type) {
		case *pageState:
			// Nothing to do.
		case contentNodeSlice:
			var updated bool
			is := make(contentNodes[contentNodePage])
			for _, ms := range v {
				b, r, err := handlePageMetaSource(ms, is, false)
				if err != nil {
					return nil, false, false, false, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
				}

				updated = updated || b

				if r && printPathWarnings {
					// TODO1 I'm not sure this is practical when we get all the matrix in play.
					hdebug.Printf("Duplicate content path: %q", n.Path())
					/*if replaced && !m.s.h.isRebuild() && m.s.conf.PrintPathWarnings {
						var messageDetail string
						if p1, ok := n.(*pageState); ok && p1.File() != nil {
							messageDetail = fmt.Sprintf(" file: %q", p1.File().Filename())
						}
						if p2, ok := u.(*pageState); ok && p2.File() != nil {
							messageDetail += fmt.Sprintf(" file: %q", p2.File().Filename())
						}

						m.s.Log.Warnf("Duplicate content path: %q%s", s, messageDetail)
					}*/
				}

			}
			return is, updated, false, false, nil
		case *pageMetaSource:
			var updated bool
			is := make(contentNodes[contentNodePage])
			b, _, err := handlePageMetaSource(v, is, false)
			if err != nil {
				return nil, false, false, false, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
			}
			updated = updated || b
			return is, updated, false, false, nil
		case *pageMeta: // TODO1 remove.
			site, found := sites[v.sitesMatrix().FirstVector()]
			if !found {
				panic(fmt.Sprintf("site not found for %v", v))
			}
			p, err := site.newPageNew(v)
			return p, true, false, false, err
		case contentNodes[contentNodePage]:
			for i, vv := range v {
				switch m := vv.(type) {
				case *pageMeta:
					var err error
					site, found := sites[m.sitesMatrix().FirstVector()] // TODO1 get rid of this interface.
					if !found {
						panic(fmt.Sprintf("site not found for %v", m))
					}
					v[i], err = site.newPageNew(m)
					if err != nil {
						return nil, false, false, false, fmt.Errorf("failed to create page %s: %w", s, err)
					}
				case *pageMetaSource:
					_, _, err := handlePageMetaSource(m, v, true)
					if err != nil {
						return nil, false, false, false, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
					}
					return v, false, false, false, nil

				default:
					// Nothing to do.
				}
			}
		default:
			panic(fmt.Sprintf("unexpected type %T", n))
		}

		return n, false, false, false, nil
	}

	var (
		seenRootSections = map[string]bool{
			"": true, // The home section.
		}
		seenHome bool
	)

	foo := func(s string, n contentNode) (n2 contentNode, replaced bool, skip bool, terminate bool, err error) {
		level := strings.Count(s, "/")

		if s != "" && !seenHome {
			homePages, _, _, _, _ := transformPages("", newHomePageMetaSource())
			treePages.InsertRaw("", homePages)
			seenHome = true
		}

		n2, replaced, skip, terminate, err = transformPages(s, n)
		if err != nil || skip || terminate {
			return
		}

		isRootSetion := level == 1 && n.isContentNodeBranch()

		if isRootSetion {
			// This is a root section.
			seenRootSections[n.PathInfo().Section()] = true
		} else if level < 3 {
			p := n.PathInfo()
			rootSection := p.Section()
			if !seenRootSections[rootSection] {
				seenRootSections[rootSection] = true
				// Try to preserve the original casing if possible.
				sectionUnnormalized := p.Unnormalized().Section()
				rootSectionPath := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, "/"+sectionUnnormalized+"/_index.md")
				rootSectionPages, _, _, _, _ := transformPages(rootSectionPath.Base(), &pageMetaSource{
					pathInfo:       rootSectionPath,
					siteMatrixBase: n2.(contentNodeForSites).siteVectors(),
					pageConfigSource: &pagemeta.PageConfig{
						PageConfigEarly: pagemeta.PageConfigEarly{
							Kind: kinds.KindSection, // TODO1 also handle taxonomies here.
						},
					},
				})
				treePages.InsertRaw(rootSectionPath.Base(), rootSectionPages)
			}
		}

		const eventNameSitesMatrix = "sitesmatrix"

		if s == "" || isRootSetion {
			if s == "" {
				seenHome = true
			}

			// Every page needs a home and a root section (.FirstSection).
			// We don't know yet what language, version, role combination that will
			// be created below, so collect that information and create the missing pages
			// on demand.
			switch nn := n2.(type) {
			case *pageState:
				n2 = contentNodes[contentNodePage]{
					nn.s.siteVector: nn,
				}
			}

			nm := n2.(contentNodes[contentNodePage])
			missingVectors := sitesmatrix.Vectors{}

			a.rw.WalkContext.AddEventListener(eventNameSitesMatrix, s,
				func(e *doctree.Event[contentNode]) {
					n := e.Source
					e.StopPropagation()
					n.forEeachContentNode(
						func(vec sitesmatrix.Vector, nn contentNode) bool {
							if _, found := nm[vec]; !found {
								missingVectors[vec] = struct{}{}
							}
							return true
						})
				},
			)
			// We need to wait until after the walk to have a complete set.
			a.rw.WalkContext.AddPostHook(
				func() error {
					if i := len(missingVectors); i > 0 {
						vec := missingVectors.One()
						kind := kinds.KindSection
						if s == "" {
							kind = kinds.KindHome // TODO1 also handle taxonomies here.
						}
						pms := &pageMetaSource{
							pathInfo:       n.PathInfo(),
							siteMatrixBase: missingVectors,
							pageConfigSource: &pagemeta.PageConfig{
								PageConfigEarly: pagemeta.PageConfigEarly{
									Kind: kind,
								},
							},
						}
						nm[vec] = pms
						_, replaced, _, _, _ := transformPages(s, nm)
						if replaced {
							// Should not happen.
							panic(fmt.Sprintf("expected no replacement for %q", s))
						}
					}
					return nil
				},
			)
		}

		if s != "" {
			a.rw.WalkContext.SendEvent(&doctree.Event[contentNode]{Source: n2, Path: s, Name: eventNameSitesMatrix})
		}

		return
	}

	transformPagesAndCreateMissingStructuralNodes := func(s string, n contentNode) (n2 contentNode, replaced bool, skip bool, terminate bool, err error) {
		n2, replaced, skip, terminate, err = foo(s, n)
		if err != nil || skip || terminate {
			return
		}
		n2.forEeachContentNode(
			func(vec sitesmatrix.Vector, nn contentNode) bool {
				if ps, ok := nn.(*pageState); ok {
					for _, viewName := range getViews(vec) {
						vals := types.ToStringSlicePreserveString(getParam(ps, viewName.plural, false))
						if vals == nil {
							continue
						}
						for _, v := range vals {
							if v == "" {
								continue
							}
							t := term{view: viewName, term: v}
							if vectors, found := seenTerms[t]; found {
								if _, found := vectors[vec]; found {
									continue
								}
								vectors[vec] = struct{}{}
							} else {
								seenTerms[t] = sitesmatrix.Vectors{
									vec: struct{}{},
								}
							}

							if true {
								// TODO1 remove below.
								continue
							}
							viewTermKey := "/" + viewName.plural + "/" + v
							pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, viewTermKey+"/_index.md")
							termKey := pi.Base()

							term, _, _, _, _ := transformPages(termKey, &pageMetaSource{
								pathInfo:       pi,
								siteMatrixBase: vec,
								term:           v,
								singular:       viewName.singular,
								pageConfigSource: &pagemeta.PageConfig{
									PageConfigEarly: pagemeta.PageConfigEarly{
										Kind: kinds.KindTerm,
									},
								},
							})
							hdebug.Printf("add %q for %v => %T", termKey, vec, term)
							a.pw.Tree.InsertRaw(termKey, term)
						}
					}
				}
				return true
			},
		)

		return
	}

	shouldSkipOrTerminate := func(s string) (skip, terminate bool) {
		owner := resourceOwnerInfo.n
		if owner == nil {
			return false, true
		}
		if !owner.isContentNodeBranch() {
			return false, false
		}

		// A resourceKey always represents a filename with extension.
		// A page key points to the logical path of a page, which when sourced from the filesystem
		// may represent a directory (bundles) or a single content file (e.g. p1.md).
		// So, to avoid any overlapping ambiguity, we start looking from the owning directory.
		for {
			s = path.Dir(s)
			ownerKey, found := treePages.LongestPrefixAll(s)
			if !found {
				return false, true
			}
			if ownerKey == resourceOwnerInfo.s {
				break
			}

			if s != ownerKey && strings.HasPrefix(s, ownerKey) {
				// Keep looking
				continue
			}

			// Stop walking downwards, someone else owns this resource.
			a.rw.SkipPrefix(ownerKey + "/")
			return true, false
		}
		return false, false
	}

	forEeachPage := func(fn func(p *pageState) bool) bool {
		switch nn := resourceOwnerInfo.n.(type) {
		case *pageState:
			return fn(nn)
		case contentNodes[contentNodePage]:
			for _, p := range nn {
				if !fn(p.(*pageState)) {
					return false
				}
			}
		default:
			panic(fmt.Sprintf("unknown type %T", nn))
		}
		return true
	}

	a.rw = &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:        treeResources,
		LockType:    lockType,
		NoShift:     true,
		WalkContext: &doctree.WalkContext[contentNode]{},
		Transform: func(s string, n contentNode) (n2 contentNode, replaced bool, skip bool, terminate bool, err error) {
			if skip, terminate = shouldSkipOrTerminate(s); skip || terminate {
				return
			}

			if contentNodeHelper.isPageNode(n) {
				return transformPagesAndCreateMissingStructuralNodes(s, n)
			}

			// TODO1 avoid creating a map for one node.
			nodes := make(contentNodes[contentNode])
			n2 = nodes
			replaced = true
			forEeachPage(
				func(p *pageState) bool {
					if _, found := nodes[p.s.siteVector]; !found {
						var rs *resourceSource
						n.forEeachContentNode(
							func(vec sitesmatrix.Vector, nn contentNode) bool {
								if r, ok := nn.(*resourceSource); ok && r.matchSiteVector(p.s.siteVector) {
									rs = r
									return false
								}
								return true
							},
						)

						if rs != nil {
							if rs.state == resourceStateNew {
								nodes[p.s.siteVector] = rs.assignSiteVector(p.s.siteVector)
							} else {
								nodes[p.s.siteVector] = rs.clone().assignSiteVector(p.s.siteVector)
							}
						}
					}

					return true
				},
			)

			return
		},
	}

	a.pw.WalkContext.AddPostHook(
		func() error {
			for k, v := range seenTerms {
				viewTermKey := "/" + k.view.plural + "/" + k.term
				pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, viewTermKey+"/_index.md")
				termKey := pi.Base()
				n, found := a.pw.Tree.GetRaw(termKey)

				if found {
					// Merge.
					n.forEeachContentNode(
						func(vec sitesmatrix.Vector, nn contentNode) bool {
							delete(v, vec)
							return true
						},
					)
				}

				if len(v) > 0 {
					p := &pageMetaSource{
						pathInfo:       pi,
						siteMatrixBase: v,
						term:           k.term,
						singular:       k.view.singular,
						pageConfigSource: &pagemeta.PageConfig{
							PageConfigEarly: pagemeta.PageConfigEarly{
								Kind: kinds.KindTerm,
							},
						},
					}
					var n2 contentNode = contentNodeSlice{n, p}
					n2, replace, _, _, err := transformPages(termKey, n2)
					if err != nil {
						return fmt.Errorf("failed to create term page %q: %w", termKey, err)
					}
					if replace {
						a.pw.Tree.InsertRaw(termKey, n2)
					}
				}
			}
			return nil
		},
	)

	a.pw.Transform = func(s string, n contentNode) (n2 contentNode, replaced bool, skip bool, terminate bool, err error) {
		n2, replaced, skip, terminate, err = transformPagesAndCreateMissingStructuralNodes(s, n)
		if err != nil || skip || terminate {
			return
		}

		// Walk nested resources.
		resourceOwnerInfo.s = s
		resourceOwnerInfo.n = n2
		a.rw.Prefix = s + "/"
		if err := a.rw.Walk(a.ctx); err != nil {
			return nil, false, false, false, err
		}

		return
	}
	a.pw.Handle = nil

	if err := a.pw.Walk(a.ctx); err != nil {
		return err
	}

	if err := a.pw.WalkContext.HandleEventsAndHooks(); err != nil {
		return err
	}

	return nil
}

// Calculate and apply aggregate values to the page tree (e.g. dates).
func (sa *sitePagesAssembler) applyAggregates() error {
	sectionPageCount := map[string]int{}

	pw := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:        sa.s.pageMap.treePages,
		LockType:    doctree.LockTypeRead,
		WalkContext: &doctree.WalkContext[contentNode]{},
	}
	rw := pw.Extend()
	rw.Tree = sa.s.pageMap.treeResources
	sa.s.lastmod = time.Time{}
	rebuild := sa.s.h.isRebuild()

	pw.Handle = func(keyPage string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
		pageBundle := n.(*pageState)

		if pageBundle.Kind() == kinds.KindTerm {
			// Delay this until they're created.
			return false, nil
		}

		if pageBundle.IsPage() {
			rootSection := pageBundle.Section()
			sectionPageCount[rootSection]++
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

		const eventName = "dates"
		if n.isContentNodeBranch() {
			wasZeroDates := pageBundle.m.pageConfig.Dates.IsAllDatesZero()
			if wasZeroDates || pageBundle.IsHome() {
				pw.WalkContext.AddEventListener(eventName, keyPage, func(e *doctree.Event[contentNode]) {
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
		pw.WalkContext.SendEvent(&doctree.Event[contentNode]{Source: n, Path: keyPage, Name: eventName})

		isBranch := n.isContentNodeBranch()
		rw.Prefix = keyPage + "/"
		rw.IncludeRawFilter = func(s string, n contentNode) bool {
			// TODO1 do some filtering here for performance.
			return true
		}

		rw.IncludeFilter = func(s string, n contentNode) bool {
			switch n.(type) {
			case *pageState:
				return true
			default:
				// We only want to handle page nodes here.
				return false
			}
		}

		rw.Handle = func(resourceKey string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
			if isBranch {
				ownerKey, _ := pw.Tree.LongestPrefix(resourceKey, false, nil)
				if ownerKey != keyPage {
					// Stop walking downwards, someone else owns this resource.
					rw.SkipPrefix(ownerKey + "/")
					return false, nil
				}
			}
			switch rs := n.(type) {
			case *pageState:
				relPath := rs.m.pathInfo.BaseRel(pageBundle.m.pathInfo)
				rs.m.resourcePath = relPath
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
	walkContext := &doctree.WalkContext[contentNode]{}

	handlePlural := func(key string) error {
		var pw *doctree.NodeShiftTreeWalker[contentNode]
		pw = &doctree.NodeShiftTreeWalker[contentNode]{
			Tree:        sa.s.pageMap.treePages,
			Prefix:      key, // We also want to include the root taxonomy nodes, so no trailing slash.
			LockType:    doctree.LockTypeRead,
			WalkContext: walkContext,
			Handle: func(s string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
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
					// TODO1
					/*var cascade *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig]
					_, data := pw.WalkContext.Data().LongestPrefix(s)
					if data != nil {
						cascade = data.(*maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig])
					}
					if err := p.setMetaPost(cascade); err != nil {
						return false, err
					}*/
					if !p.s.shouldBuild(p) {
						sa.s.pageMap.treePages.Delete(s)
						sa.s.pageMap.treeTaxonomyEntries.DeletePrefix(paths.AddTrailingSlash(s))
					} else if err := sa.s.pageMap.treeTaxonomyEntries.WalkPrefix(
						doctree.LockTypeRead,
						paths.AddTrailingSlash(s),
						func(ss string, wn *weightedContentNode) (bool, error) {
							// Send the date info up the tree.
							pw.WalkContext.SendEvent(&doctree.Event[contentNode]{Source: wn.n, Path: ss, Name: eventName})
							return false, nil
						},
					); err != nil {
						return false, err
					}
				}

				// Send the date info up the tree.
				pw.WalkContext.SendEvent(&doctree.Event[contentNode]{Source: n, Path: s, Name: eventName})

				if p.m.pageConfig.Dates.IsAllDatesZero() {
					pw.WalkContext.AddEventListener(eventName, s, func(e *doctree.Event[contentNode]) {
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

func (sa *sitePagesAssembler) assembleTerms() error {
	if sa.s.pageMap.cfg.taxonomyTermDisabled {
		return nil
	}

	var (
		pages   = sa.s.pageMap.treePages
		entries = sa.s.pageMap.treeTaxonomyEntries
		views   = sa.s.pageMap.cfg.taxonomyConfig.views
	)

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     pages,
		LockType: lockType,
		Handle: func(s string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
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
						panic(fmt.Sprintf("missing term page for %q", viewTermKey))
					}

					m := term.(*pageState).m
					m.term = v
					m.singular = viewName.singular

					if s == "" {
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

func (sa *sitePagesAssembler) assemblePagesStepFinal() error {
	return sa.assembleResources()
}

func (sa *sitePagesAssembler) assembleResources() error {
	pagesTree := sa.s.pageMap.treePages

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     pagesTree,
		LockType: lockType,
		Handle: func(s string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
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
			/*duplicateResourceFiles := true
			if ps.m.pageConfig.ContentMediaType.IsMarkdown() {
				duplicateResourceFiles = ps.s.ContentSpec.Converters.GetMarkupConfig().Goldmark.DuplicateResourceFiles
			}*/

			// TODO1 duplicateResourceFiles = duplicateResourceFiles || ps.s.Conf.IsMultihost()

			err := sa.s.pageMap.forEachResourceInPage(
				ps, lockType,
				false,
				nil,
				func(resourceKey string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
					if _, ok := n.(*pageState); ok {
						return false, nil
					}
					rs := n.(*resourceSource)

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
							rs.rc.Name = relPathOriginal // TODO1, this is shared.
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

	/*if err := sa.addMissingTaxonomies(); err != nil {
		return err
	}

	if err := sa.addMissingRootSections(); err != nil { // TODO1 see above.
		return err
	}

	if err := sa.addStandalonePages(); err != nil {
		return err
	}*/

	if err := sa.applyAggregates(); err != nil {
		return err
	}
	return nil
}

func (sa *sitePagesAssembler) assemblePagesStep2() error {
	if err := sa.removeShouldNotBuild(); err != nil { // TODO1
		return err
	}
	if err := sa.assembleTerms(); err != nil {
		return err
	}

	return nil
}

// Remove any leftover node that we should not build for some reason (draft, expired, scheduled in the future).
// Note that for the home and section kinds we just disable the nodes to preserve the structure.
func (sa *sitePagesAssembler) removeShouldNotBuild() error {
	s := sa.s
	var keys []string
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		LockType: doctree.LockTypeRead,
		Tree:     sa.s.pageMap.treePages,
		Handle: func(key string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
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

func (sa *sitePagesAssembler) addMissingRootSections() error {
	var hasHome bool

	// Add missing root sections.
	seen := map[string]bool{}
	var w *doctree.NodeShiftTreeWalker[contentNode]
	w = &doctree.NodeShiftTreeWalker[contentNode]{
		LockType: doctree.LockTypeWrite,
		Tree:     sa.s.pageMap.treePages,
		Handle: func(s string, n contentNode, match sitesmatrix.Dimension) (bool, error) {
			if n == nil {
				panic("n is nil")
			}

			ps := n.(*pageState)

			if s == "" {
				hasHome = true
				sa.s.home = ps
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
					pageMetaSource: &pageMetaSource{
						pathInfo: pth,
					},
				}
				ps, err := sa.s.newPageNew(m)
				if err != nil {
					return false, err
				}

				w.Tree.InsertIntoCurrentDimension(ps.PathInfo().Base(), ps)
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
		p := sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, "/_index.md")

		m := &pageMeta{
			pageMetaParams: &pageMetaParams{},
			pageMetaSource: &pageMetaSource{
				pathInfo: p,
				pageConfigSource: &pagemeta.PageConfig{
					PageConfigEarly: pagemeta.PageConfigEarly{
						Kind: kinds.KindHome,
					},
				},
			},
		}
		n, err := sa.s.newPageNew(m)
		if err != nil {
			return err
		}

		w.Tree.InsertIntoCurrentDimensionWithLock(n.PathInfo().Base(), n)
		sa.s.home = n
	}

	return nil
}

func (a *allPagesAssembler) createMissingTaxonomies() error {
	if a.m.cfg.taxonomyDisabled && a.m.cfg.taxonomyTermDisabled {
		return nil
	}

	tree := a.m.treePages

	viewLanguages := map[viewName][]int{}
	for _, s := range a.h.sitesLanguages {
		if s.pageMap.cfg.taxonomyDisabled && s.pageMap.cfg.taxonomyTermDisabled {
			continue
		}
		for _, viewName := range s.pageMap.cfg.taxonomyConfig.views {
			viewLanguages[viewName] = append(viewLanguages[viewName], s.siteVector.Language())
		}
	}

	commit := tree.Lock(true) // TODO1 revise locking for this flow.
	defer commit()

	for viewName, languages := range viewLanguages {
		key := viewName.pluralTreeKey
		matrixAllForLanguages := sitesmatrix.NewIntSetsBuilder(a.h.Conf.ConfiguredDimensions()).WithLanguageIndices(languages...).WithAllIfNotSet().Build()

		pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, key+"/_index.md")
		p := &pageMetaSource{
			pathInfo:       pi,
			singular:       viewName.singular,
			siteMatrixBase: matrixAllForLanguages,
			pageConfigSource: &pagemeta.PageConfig{
				PageConfigEarly: pagemeta.PageConfigEarly{
					Kind: kinds.KindTaxonomy,
				},
			},
		}
		tree.AppendRaw(key, p)

	}

	return nil
}

// Create the fixed output pages, e.g. sitemap.xml, if not already there.
// TODO2 revise, especially around disabled and config per language/site.
func (a *allPagesAssembler) createMissingStandalonePages() error {
	m := a.m
	tree := m.treePages

	commit := tree.Lock(true)
	defer commit()

	addStandalone := func(key, kind string, f output.Format) {
		if !a.h.Conf.IsMultihost() {
			switch kind {
			case kinds.KindSitemapIndex, kinds.KindRobotsTXT:
				// Only one for all dimensions. TODO1
				/*if !s.siteVector.IsFirst() {
					return
				}*/
			}
		}

		// TODO1 per site.
		/*if !sa.h.Conf.IsKindEnabled(kind) || tree.Has(key) {
			return
		}*/

		p := &pageMetaSource{
			pathInfo:               a.h.Conf.PathParser().Parse(files.ComponentFolderContent, key+f.MediaType.FirstSuffix.FullSuffix),
			standaloneOutputFormat: f,
			pageConfigSource: &pagemeta.PageConfig{
				PageConfigEarly: pagemeta.PageConfigEarly{
					Kind: kind,
				},
			},
		}

		tree.InsertRaw(key, p)
	}

	addStandalone("/404", kinds.KindStatus404, output.HTTPStatus404HTMLFormat)

	if a.h.Configs.Base.EnableRobotsTXT {
		if m.i == 0 || a.h.Conf.IsMultihost() {
			addStandalone("/_robots", kinds.KindRobotsTXT, output.RobotsTxtFormat)
		}
	}

	sitemapEnabled := false
	for _, s := range a.h.Sites {
		if s.conf.IsKindEnabled(kinds.KindSitemap) {
			sitemapEnabled = true
			break
		}
	}

	if sitemapEnabled {
		of := output.SitemapFormat
		if a.h.Configs.Base.Sitemap.Filename != "" {
			of.BaseName = paths.Filename(a.h.Configs.Base.Sitemap.Filename)
		}
		addStandalone("/_sitemap", kinds.KindSitemap, of)

		skipSitemapIndex := a.h.Conf.IsMultihost() || !(a.h.Conf.DefaultContentLanguageInSubdir() || a.h.Conf.IsMultilingual())
		if !skipSitemapIndex {
			of = output.SitemapIndexFormat
			if a.h.Configs.Base.Sitemap.Filename != "" {
				of.BaseName = paths.Filename(a.h.Configs.Base.Sitemap.Filename)
			}
			addStandalone("/_sitemapindex", kinds.KindSitemapIndex, of)
		}
	}

	return nil
}

func (a *allPagesAssembler) createMissingPages() error {
	if err := a.createMissingTaxonomies(); err != nil {
		return err
	}
	if err := a.createMissingStandalonePages(); err != nil {
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
			pi := sa.s.Conf.PathParser().Parse(files.ComponentFolderContent, key+"/_index.md")
			m := &pageMeta{
				pageMetaSource: &pageMetaSource{
					pathInfo: pi,
					singular: viewName.singular,
					pageConfigSource: &pagemeta.PageConfig{
						PageConfigEarly: pagemeta.PageConfigEarly{
							Kind: kinds.KindTaxonomy,
						},
					},
				},
				pageMetaParams: &pageMetaParams{},
			}
			p, err := sa.s.newPageNew(m)
			if err != nil {
				return fmt.Errorf("failed to create taxonomy %s: %w", viewName.plural, err)
			}
			tree.InsertIntoValuesDimension(key, p)
		}
	}

	return nil
}
