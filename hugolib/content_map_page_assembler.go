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
	"cmp"
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
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
	r   para.Runner

	assembleChanges            *WhatChanged
	assembleSectionsInParallel bool

	pwRoot *doctree.NodeShiftTreeWalker[contentNode] // walks pages.
	rwRoot *doctree.NodeShiftTreeWalker[contentNode] // walks resources.

	// Walking state.
	seenTerms        *maps.Map[term, sitesmatrix.Vectors]
	droppedPages     *maps.Map[*Site, []string] // e.g. drafts, expired, future.
	seenRootSections *maps.Map[string, bool]
	seenHome         bool // set before we fan out to multiple goroutines.
}

func newAllPagesAssembler(
	ctx context.Context,
	h *HugoSites,
	m *pageMap,
	assembleChanges *WhatChanged,
) *allPagesAssembler {
	rw := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:        m.treeResources,
		LockType:    doctree.LockTypeNone,
		NoShift:     true,
		WalkContext: &doctree.WalkContext[contentNode]{},
	}
	pw := rw.Extend()
	pw.Tree = m.treePages

	seenRootSections := maps.NewMap[string, bool]()
	seenRootSections.Set("", true) // home.

	return &allPagesAssembler{
		ctx:                        ctx,
		h:                          h,
		m:                          m,
		assembleChanges:            assembleChanges,
		seenTerms:                  maps.NewMap[term, sitesmatrix.Vectors](),
		droppedPages:               maps.NewMap[*Site, []string](),
		seenRootSections:           seenRootSections,
		assembleSectionsInParallel: true,
		pwRoot:                     pw,
		rwRoot:                     rw,
	}
}

type sitePagesAssembler struct {
	s               *Site
	assembleChanges *WhatChanged
	a               *allPagesAssembler
	ctx             context.Context
}

func (a *allPagesAssembler) createAllPages() error {
	defer func() {
		for site, dropped := range a.droppedPages.All() {
			for _, s := range dropped {
				site.pageMap.treePages.Delete(s)
				site.pageMap.resourceTrees.DeletePrefix(paths.AddTrailingSlash(s))
			}
		}
	}()

	if a.h.Conf.Watching() {
		defer func() {
			if a.h.isRebuild() && a.h.previousSeenTerms != nil {
				// Find removed terms.
				for t := range a.h.previousSeenTerms.All() {
					if _, found := a.seenTerms.Lookup(t); !found {
						// This term has been removed.
						a.pwRoot.Tree.Delete(t.view.pluralTreeKey)
						a.rwRoot.Tree.DeletePrefix(t.view.pluralTreeKey + "/")
					}
				}
				// Find new terms.
				for t := range a.seenTerms.All() {
					if _, found := a.h.previousSeenTerms.Lookup(t); !found {
						// This term is new.
						if n, ok := a.pwRoot.Tree.GetRaw(t.view.pluralTreeKey); ok {
							a.assembleChanges.Add(cnh.GetIdentity(n))
						}
					}
				}
			}
			a.h.previousSeenTerms = a.seenTerms
		}()
	}
	workers := para.New(config.GetNumWorkerMultiplier())
	a.r, _ = workers.Start(context.Background())
	if err := cmp.Or(a.doCreatePages(""), a.r.Wait()); err != nil {
		return err
	}
	if err := a.pwRoot.WalkContext.HandleHooks1AndEventsAndHooks2(); err != nil {
		return err
	}
	return nil
}

func (a *allPagesAssembler) doCreatePages(prefix string) error {
	var (
		sites     = a.h.sitesVersionsRolesMap
		h         = a.h
		treePages = a.m.treePages

		getViews = func(vec sitesmatrix.Vector) []viewName {
			return h.languageSiteForSiteVector(vec).pageMap.cfg.taxonomyConfig.views
		}

		pw *doctree.NodeShiftTreeWalker[contentNode] // walks pages.
		rw *doctree.NodeShiftTreeWalker[contentNode] // walks resources.

		isRootWalk = prefix == ""
	)

	if isRootWalk {
		pw = a.pwRoot
		rw = a.rwRoot
	} else {
		// Sub-walkers for a specific prefix.
		pw = a.pwRoot.Extend()
		pw.Prefix = prefix + "/"

		rw = a.rwRoot.Extend() // rw will get its prefix(es) set later.

		pw.TransformDelayInsert = true
		rw.TransformDelayInsert = true

	}

	resourceOwnerInfo := struct {
		n contentNode
		s string
	}{}

	newHomePageMetaSource := func() *pageMetaSource {
		pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, "/_index.md")
		return &pageMetaSource{
			pathInfo:            pi,
			sitesMatrixBase:     a.h.Conf.AllSitesMatrix(),
			sitesMatrixBaseOnly: true,
			pageConfigSource: &pagemeta.PageConfigEarly{
				Kind: kinds.KindHome,
			},
		}
	}

	if isRootWalk {
		if err := a.createMissingPages(); err != nil {
			return err
		}

		if treePages.Len() == 0 {
			// No pages, insert a home page to get something to walk on.
			p := newHomePageMetaSource()
			treePages.InsertRaw(p.pathInfo.Base(), p)
		}
	}

	getCascades := func(wctx *doctree.WalkContext[contentNode], s string) *page.PageMatcherParamsConfigs {
		if wctx == nil {
			panic("nil walk context")
		}
		var cascades *page.PageMatcherParamsConfigs
		data := wctx.Data()
		if s != "" {
			_, data := data.LongestPrefix(s)

			if data != nil {
				cascades = data.(*page.PageMatcherParamsConfigs)
			}
		}

		if cascades == nil {
			for s := range a.h.allSiteLanguages(nil) {
				cascades = cascades.Append(s.conf.Cascade)
			}
		}
		return cascades
	}

	doTransformPages := func(s string, n contentNode, cascades *page.PageMatcherParamsConfigs) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		defer func() {
			cascadesLen := cascades.Len()
			if n2 == nil {
				if ns < doctree.NodeTransformStateSkip {
					ns = doctree.NodeTransformStateSkip
				}
			} else {
				n2.forEeachContentNode(
					func(vec sitesmatrix.Vector, nn contentNode) bool {
						if pms, ok := nn.(contentNodeCascadeProvider); ok {
							cascades = cascades.Prepend(pms.getCascade())
						}
						return true
					},
				)
			}

			if s == "" || cascades.Len() > cascadesLen {
				// New cascade values added, pass them downwards.
				rw.WalkContext.Data().Insert(paths.AddTrailingSlash(s), cascades)
			}
		}()

		var handlePageMetaSource func(v contentNode, n contentNodesMap, replaceVector bool) (n2 contentNode, err error)
		handlePageMetaSource = func(v contentNode, n contentNodesMap, replaceVector bool) (n2 contentNode, err error) {
			if n != nil {
				n2 = n
			}

			switch ms := v.(type) {
			case *pageMetaSource:
				if err := ms.initEarly(a.h, cascades); err != nil {
					return nil, err
				}

				sitesMatrix := ms.pageConfigSource.SitesMatrix

				sitesMatrix.ForEachVector(func(vec sitesmatrix.Vector) bool {
					site, found := sites[vec]
					if !found {
						panic(fmt.Sprintf("site not found for %v", vec))
					}

					var p *pageState
					p, err = site.newPageFromPageMetasource(ms, cascades)
					if err != nil {
						return false
					}

					var drop bool
					if !site.shouldBuild(p) {
						switch p.Kind() {
						case kinds.KindHome, kinds.KindSection, kinds.KindTaxonomy:
							// We need to keep these for the structure, but disable
							// them so they don't get listed/rendered.
							(&p.m.pageConfig.Build).Disable()
						default:
							// Skip this page.
							a.droppedPages.WithWriteLock(
								func(m map[*Site][]string) {
									m[site] = append(m[site], s)
								},
							)

							drop = true
						}
					}

					if !drop && n == nil {
						if n2 == nil {
							// Avoid creating a map for one node.
							n2 = p
						} else {
							// Upgrade to a map.
							n = make(contentNodesMap)
							ps := n2.(*pageState)
							n[ps.s.siteVector] = ps
							n2 = n
						}
					}

					if n == nil {
						return true
					}

					pp, found := n[vec]

					var w1, w2 int
					if wp, ok := pp.(contentNodeContentWeightProvider); ok {
						w1 = wp.contentWeight()
					}
					w2 = p.contentWeight()

					if found && !replaceVector && w1 > w2 {
						return true
					}

					n[vec] = p
					return true
				})
			case *pageState:
				if n == nil {
					n2 = ms
					return
				}
				n[ms.s.siteVector] = ms
			case contentNodesMap:
				for _, vv := range ms {
					var err error
					n2, err = handlePageMetaSource(vv, n, replaceVector)
					if err != nil {
						return nil, err
					}
					if m, ok := n2.(contentNodesMap); ok {
						n = m
					}

				}
			default:
				panic(fmt.Sprintf("unexpected type %T", v))
			}

			return
		}

		// The common case.
		ns = doctree.NodeTransformStateReplaced

		handleContentNodeSeq := func(v contentNodeSeq) (contentNode, doctree.NodeTransformState, error) {
			is := make(contentNodesMap)
			for ms := range v {
				_, err := handlePageMetaSource(ms, is, false)
				if err != nil {
					return nil, 0, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
				}
			}
			return is, ns, nil
		}

		switch v := n.(type) {
		case contentNodeSeq, contentNodes:
			return handleContentNodeSeq(contentNodeToSeq(v))
		case *pageMetaSource:
			n2, err = handlePageMetaSource(v, nil, false)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
			}
			return
		case *pageState:
			// Nothing to do.
			ns = doctree.NodeTransformStateNone
			return v, ns, nil
		case contentNodesMap:
			ns = doctree.NodeTransformStateNone
			for _, vv := range v {
				switch m := vv.(type) {
				case *pageMetaSource:
					ns = doctree.NodeTransformStateUpdated
					_, err := handlePageMetaSource(m, v, true)
					if err != nil {
						return nil, 0, fmt.Errorf("failed to create page from pageMetaSource %s: %w", s, err)
					}
				default:
					// Nothing to do.
				}
			}

			return v, ns, nil

		default:
			panic(fmt.Sprintf("unexpected type %T", n))
		}
	}

	var unpackPageMetaSources func(n contentNode) contentNode
	unpackPageMetaSources = func(n contentNode) contentNode {
		switch nn := n.(type) {
		case *pageState:
			return nn.m.pageMetaSource
		case contentNodesMap:
			if len(nn) == 0 {
				return nil
			}
			var iter contentNodeSeq = func(yield func(contentNode) bool) {
				seen := map[*pageMetaSource]struct{}{}
				for _, v := range nn {
					vv := unpackPageMetaSources(v)
					pms := vv.(*pageMetaSource)
					if _, found := seen[pms]; !found {
						if !yield(pms) {
							return
						}
						seen[pms] = struct{}{}
					}
				}
			}
			return iter
		case contentNodes:
			if len(nn) == 0 {
				return nil
			}
			var iter contentNodeSeq = func(yield func(contentNode) bool) {
				seen := map[*pageMetaSource]struct{}{}
				for _, v := range nn {
					vv := unpackPageMetaSources(v)
					pms := vv.(*pageMetaSource)
					if _, found := seen[pms]; !found {
						if !yield(pms) {
							return
						}
						seen[pms] = struct{}{}
					}
				}
			}
			return iter
		case *pageMetaSource:
			return nn
		default:
			panic(fmt.Sprintf("unexpected type %T", n))
		}
	}

	transformPages := func(s string, n contentNode, cascades *page.PageMatcherParamsConfigs) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		if a.h.isRebuild() {
			cascadesPrevious := getCascades(h.previousPageTreesWalkContext, s)
			h1, h2 := cascadesPrevious.SourceHash(), cascades.SourceHash()
			if h1 != h2 {
				// Force rebuild from the source.
				n = unpackPageMetaSources(n)
			}
		}
		if n == nil {
			panic("nil node")
		}
		n2, ns, err = doTransformPages(s, n, cascades)

		return
	}

	transformPagesAndCreateMissingHome := func(s string, n contentNode, isResource bool, cascades *page.PageMatcherParamsConfigs) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		if n == nil {
			panic("nil node " + s)
		}
		level := strings.Count(s, "/")

		if s == "" {
			a.seenHome = true
		}

		if !isResource && s != "" && !a.seenHome {
			a.seenHome = true
			var homePages contentNode
			homePages, ns, err = transformPages("", newHomePageMetaSource(), cascades)
			if err != nil {
				return
			}
			treePages.InsertRaw("", homePages)
		}

		n2, ns, err = transformPages(s, n, cascades)
		if err != nil || ns >= doctree.NodeTransformStateSkip {
			return
		}

		if n2 == nil {
			ns = doctree.NodeTransformStateSkip
			return
		}

		if isResource {
			// Done.
			return
		}

		isTaxonomy := !a.h.getFirstTaxonomyConfig(s).IsZero()
		isRootSection := !isTaxonomy && level == 1 && cnh.isBranchNode(n)

		if isRootSection {
			// This is a root section.
			a.seenRootSections.SetIfAbsent(cnh.PathInfo(n).Section(), true)
		} else if !isTaxonomy {
			p := cnh.PathInfo(n)
			rootSection := p.Section()
			_, err := a.seenRootSections.GetOrCreate(rootSection, func() (bool, error) {
				// Try to preserve the original casing if possible.
				sectionUnnormalized := p.Unnormalized().Section()
				rootSectionPath := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, "/"+sectionUnnormalized+"/_index.md")
				var rootSectionPages contentNode
				rootSectionPages, _, err = transformPages(rootSectionPath.Base(), &pageMetaSource{
					pathInfo:        rootSectionPath,
					sitesMatrixBase: n2.(contentNodeForSites).siteVectors(),
					pageConfigSource: &pagemeta.PageConfigEarly{
						Kind: kinds.KindSection,
					},
				}, cascades)
				if err != nil {
					return true, err
				}
				treePages.InsertRaw(rootSectionPath.Base(), rootSectionPages)
				return true, nil
			})
			if err != nil {
				return nil, 0, err
			}

		}

		const eventNameSitesMatrix = "sitesmatrix"

		if s == "" || isRootSection {

			// Every page needs a home and a root section (.FirstSection).
			// We don't know yet what language, version, role combination that will
			// be created below, so collect that information and create the missing pages
			// on demand.
			nm, replaced := contentNodeToContentNodesPage(n2)

			missingVectorsForHomeOrRootSection := sitesmatrix.Vectors{}

			if s == "" {
				// We need a complete set of home pages.
				a.h.Conf.AllSitesMatrix().ForEachVector(func(vec sitesmatrix.Vector) bool {
					if _, found := nm[vec]; !found {
						missingVectorsForHomeOrRootSection[vec] = struct{}{}
					}
					return true
				})
			} else {
				pw.WalkContext.AddEventListener(eventNameSitesMatrix, s,
					func(e *doctree.Event[contentNode]) {
						n := e.Source
						e.StopPropagation()
						n.forEeachContentNode(
							func(vec sitesmatrix.Vector, nn contentNode) bool {
								if _, found := nm[vec]; !found {
									missingVectorsForHomeOrRootSection[vec] = struct{}{}
								}
								return true
							})
					},
				)
			}

			// We need to wait until after the walk to have a complete set.
			pw.WalkContext.HooksPost2().Push(
				func() error {
					if i := len(missingVectorsForHomeOrRootSection); i > 0 {
						// Pick one, the rest will be created later.
						vec := missingVectorsForHomeOrRootSection.VectorSample()

						kind := kinds.KindSection
						if s == "" {
							kind = kinds.KindHome
						}

						pms := &pageMetaSource{
							pathInfo:            cnh.PathInfo(n),
							sitesMatrixBase:     missingVectorsForHomeOrRootSection,
							sitesMatrixBaseOnly: true,
							pageConfigSource: &pagemeta.PageConfigEarly{
								Kind: kind,
							},
						}
						nm[vec] = pms

						_, ns, err := transformPages(s, nm, cascades)
						if err != nil {
							return err
						}
						if ns == doctree.NodeTransformStateReplaced {
							// Should not happen.
							panic(fmt.Sprintf("expected no replacement for %q", s))
						}

						if replaced {
							pw.Tree.InsertRaw(s, nm)
						}
					}
					return nil
				},
			)
		}

		if s != "" {
			rw.WalkContext.SendEvent(&doctree.Event[contentNode]{Source: n2, Path: s, Name: eventNameSitesMatrix})
		}

		return
	}

	transformPagesAndCreateMissingStructuralNodes := func(s string, n contentNode, isResource bool) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		cascades := getCascades(pw.WalkContext, s)
		n2, ns, err = transformPagesAndCreateMissingHome(s, n, isResource, cascades)
		if err != nil || ns >= doctree.NodeTransformStateSkip {
			return
		}
		n2.forEeachContentNode(
			func(vec sitesmatrix.Vector, nn contentNode) bool {
				if ps, ok := nn.(*pageState); ok {
					if ps.m.noLink() {
						return true
					}
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
							a.seenTerms.WithWriteLock(func(m map[term]sitesmatrix.Vectors) {
								vectors, found := m[t]
								if !found {
									m[t] = sitesmatrix.Vectors{
										vec: struct{}{},
									}
									return
								}
								vectors[vec] = struct{}{}
							})
						}
					}
				}
				return true
			},
		)

		return
	}

	shouldSkipOrTerminate := func(s string) (ns doctree.NodeTransformState) {
		owner := resourceOwnerInfo.n
		if owner == nil {
			return doctree.NodeTransformStateTerminate
		}
		if !cnh.isBranchNode(owner) {
			return
		}

		// A resourceKey always represents a filename with extension.
		// A page key points to the logical path of a page, which when sourced from the filesystem
		// may represent a directory (bundles) or a single content file (e.g. p1.md).
		// So, to avoid any overlapping ambiguity, we start looking from the owning directory.
		for {
			s = path.Dir(s)
			ownerKey, found := treePages.LongestPrefixRaw(s)
			if !found {
				return doctree.NodeTransformStateTerminate
			}
			if ownerKey == resourceOwnerInfo.s {
				break
			}

			if s != ownerKey && strings.HasPrefix(s, ownerKey) {
				// Keep looking
				continue
			}

			// Stop walking downwards, someone else owns this resource.
			rw.SkipPrefix(ownerKey + "/")
			return doctree.NodeTransformStateSkip
		}
		return
	}

	forEeachResourceOwnerPage := func(fn func(p *pageState) bool) bool {
		switch nn := resourceOwnerInfo.n.(type) {
		case *pageState:
			return fn(nn)
		case contentNodesMap:
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

	rw.Transform = func(s string, n contentNode) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		if ns = shouldSkipOrTerminate(s); ns >= doctree.NodeTransformStateSkip {
			return
		}

		if cnh.isPageNode(n) {
			return transformPagesAndCreateMissingStructuralNodes(s, n, true)
		}

		nodes := make(contentNodesMap)
		ns = doctree.NodeTransformStateReplaced
		n2 = nodes

		forEeachResourceOwnerPage(
			func(p *pageState) bool {
				duplicateResourceFiles := a.h.Cfg.IsMultihost()
				if !duplicateResourceFiles && p.m.pageConfigSource.ContentMediaType.IsMarkdown() {
					duplicateResourceFiles = p.s.ContentSpec.Converters.GetMarkupConfig().Goldmark.DuplicateResourceFiles
				}

				if _, found := nodes[p.s.siteVector]; !found {
					var rs *resourceSource
					match := cnh.findContentNodeForSiteVector(p.s.siteVector, duplicateResourceFiles, contentNodeToSeq(n))
					if match == nil {
						return true
					}

					rs = match.(*resourceSource)

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
	}

	// Create  missing term pages.
	pw.WalkContext.HooksPost2().Push(
		func() error {
			for k, v := range a.seenTerms.All() {
				viewTermKey := "/" + k.view.plural + "/" + k.term

				pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, viewTermKey+"/_index.md")
				termKey := pi.Base()

				n, found := pw.Tree.GetRaw(termKey)

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
						pathInfo:        pi,
						sitesMatrixBase: v,
						pageConfigSource: &pagemeta.PageConfigEarly{
							Kind: kinds.KindTerm,
						},
					}
					var n2 contentNode = p
					if found {
						n2 = contentNodes{n, p}
					}
					n2, ns, err := transformPages(termKey, n2, getCascades(pw.WalkContext, termKey))
					if err != nil {
						return fmt.Errorf("failed to create term page %q: %w", termKey, err)
					}

					switch ns {
					case doctree.NodeTransformStateReplaced:
						pw.Tree.InsertRaw(termKey, n2)
					}

				}
			}
			return nil
		},
	)

	pw.Transform = func(s string, n contentNode) (n2 contentNode, ns doctree.NodeTransformState, err error) {
		n2, ns, err = transformPagesAndCreateMissingStructuralNodes(s, n, false)

		if err != nil || ns >= doctree.NodeTransformStateSkip {
			return
		}

		if iep, ok := n2.(contentNodeIsEmptyProvider); ok && iep.isEmpty() {
			ns = doctree.NodeTransformStateDeleted
		}

		if ns == doctree.NodeTransformStateDeleted {
			return
		}

		// Walk nested resources.
		resourceOwnerInfo.s = s
		resourceOwnerInfo.n = n2
		rw = rw.WithPrefix(s + "/")
		if err := rw.Walk(a.ctx); err != nil {
			return nil, 0, err
		}

		if a.assembleSectionsInParallel && prefix == "" && s != "" && cnh.isBranchNode(n) && a.h.getFirstTaxonomyConfig(s).IsZero() {
			// Handle this branch's descendants in its own goroutine.
			pw.SkipPrefix(s + "/")
			a.r.Run(func() error {
				return a.doCreatePages(s)
			})
		}

		return
	}

	pw.Handle = nil

	if err := pw.Walk(a.ctx); err != nil {
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

	pw.Handle = func(keyPage string, n contentNode) (bool, error) {
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
			if pageBundle.IsHome() || pageBundle.IsSection() {
				oldDates := pageBundle.m.pageConfig.Dates

				// We need to wait until after the walk to determine if any of the dates have changed.
				pw.WalkContext.HooksPost2().Push(
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
		if cnh.isBranchNode(n) {
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

		isBranch := cnh.isBranchNode(n)
		rw.Prefix = keyPage + "/"
		rw.IncludeRawFilter = func(s string, n contentNode) bool {
			// Only page nodes.
			return cnh.isPageNode(n)
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

		rw.Handle = func(resourceKey string, n contentNode) (bool, error) {
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

	if err := pw.WalkContext.HandleHooks1AndEventsAndHooks2(); err != nil {
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
			Handle: func(s string, n contentNode) (bool, error) {
				p := n.(*pageState)

				if p.Kind() != kinds.KindTerm && p.Kind() != kinds.KindTaxonomy {
					// Already handled.
					return false, nil
				}

				const eventName = "dates"

				if p.Kind() == kinds.KindTerm {
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

	if err := walkContext.HandleHooks1AndEventsAndHooks2(); err != nil {
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
		Handle: func(s string, n contentNode) (bool, error) {
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
					termNode := pages.Get(pi.Base())
					if termNode == nil {
						// This means that the term page has been disabled (e.g. a draft).
						continue
					}

					m := termNode.(*pageState).m
					m.term = v
					m.singular = viewName.singular

					if s == "" {
						s = "/"
					}

					key := pi.Base() + s

					entries.Insert(key, &weightedContentNode{
						weight: weight,
						n:      n,
						term:   &pageWithOrdinal{pageState: termNode.(*pageState), ordinal: i},
					})
				}
			}

			return false, nil
		},
	}

	if err := w.Walk(sa.ctx); err != nil {
		return err
	}

	return nil
}

func (sa *sitePagesAssembler) assemblePagesStepFinal() error {
	if err := sa.assembleResourcesAndSetHome(); err != nil {
		return err
	}
	return nil
}

func (sa *sitePagesAssembler) assembleResourcesAndSetHome() error {
	pagesTree := sa.s.pageMap.treePages

	lockType := doctree.LockTypeWrite
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:     pagesTree,
		LockType: lockType,
		Handle: func(s string, n contentNode) (bool, error) {
			ps := n.(*pageState)

			if s == "" {
				sa.s.home = ps
			} else if ps.s.home == nil {
				panic(fmt.Sprintf("[%v] expected home page to be set for %q", sa.s.siteVector, s))
			}

			// This is a little out of place, but is conveniently put here.
			// Check if translationKey is set by user.
			// This is to support the manual way of setting the translationKey in front matter.
			if ps.m.pageConfig.TranslationKey != "" {
				sa.s.h.translationKeyPages.Append(ps.m.pageConfig.TranslationKey, ps)
			}

			if !sa.s.h.isRebuild() {
				if ps.hasRenderableOutput() {
					// For multi output pages this will not be complete, but will have to do for now.
					sa.s.h.progressReporter.numPagesToRender.Add(1)
				}
			}

			// Prepare resources for this page.
			ps.shiftToOutputFormat(true, 0)
			targetPaths := ps.targetPaths()
			baseTarget := targetPaths.SubResourceBaseTarget

			err := sa.s.pageMap.forEachResourceInPage(
				ps, lockType,
				false,
				nil,
				func(resourceKey string, n contentNode) (bool, error) {
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
	if err := sa.applyAggregates(); err != nil {
		return err
	}
	return nil
}

func (sa *sitePagesAssembler) assemblePagesStep2() error {
	if err := sa.assembleTerms(); err != nil {
		return err
	}

	if err := sa.applyAggregatesToTaxonomiesAndTerms(); err != nil {
		return err
	}

	return nil
}

// No locking.
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

	for viewName, languages := range viewLanguages {
		key := viewName.pluralTreeKey
		if a.h.isRebuild() {
			if v := tree.Get(key); v != nil {
				// Already there.
				continue
			}
		}

		matrixAllForLanguages := sitesmatrix.NewIntSetsBuilder(a.h.Conf.ConfiguredDimensions()).WithLanguageIndices(languages...).WithAllIfNotSet().Build()

		pi := a.h.Conf.PathParser().Parse(files.ComponentFolderContent, key+"/_index.md")
		p := &pageMetaSource{
			pathInfo:        pi,
			sitesMatrixBase: matrixAllForLanguages,
			pageConfigSource: &pagemeta.PageConfigEarly{
				Kind: kinds.KindTaxonomy,
			},
		}
		tree.AppendRaw(key, p)
	}

	return nil
}

// Create the fixed output pages, e.g. sitemap.xml, if not already there.
// TODO2 revise, especially around disabled and config per language/site.
// No locking.
func (a *allPagesAssembler) createMissingStandalonePages() error {
	m := a.m
	tree := m.treePages
	oneSiteStore := (sitesmatrix.Vectors{sitesmatrix.Vector{0, 0, 0}: struct{}{}}).ToVectorStore()

	addStandalone := func(key, kind string, f output.Format) {
		if !a.h.Conf.IsKindEnabled(kind) || tree.Has(key) {
			return
		}

		var sitesMatrixBase sitesmatrix.VectorIterator

		sitesMatrixBase = a.h.Conf.AllSitesMatrix()
		if !a.h.Conf.IsMultihost() {
			switch kind {
			case kinds.KindSitemapIndex, kinds.KindRobotsTXT:
				// First site only.
				sitesMatrixBase = oneSiteStore
			}
		}

		p := &pageMetaSource{
			pathInfo:               a.h.Conf.PathParser().Parse(files.ComponentFolderContent, key+f.MediaType.FirstSuffix.FullSuffix),
			standaloneOutputFormat: f,
			sitesMatrixBase:        sitesMatrixBase,
			sitesMatrixBaseOnly:    true,
			pageConfigSource: &pagemeta.PageConfigEarly{
				Kind: kind,
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

// No locking.
func (a *allPagesAssembler) createMissingPages() error {
	if err := a.createMissingTaxonomies(); err != nil {
		return err
	}

	if !a.h.isRebuild() {
		if err := a.createMissingStandalonePages(); err != nil {
			return err
		}
	}
	return nil
}
