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
	"iter"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/hugofs"
)

// Used to mark ambiguous keys in reverse index lookups.
var ambiguousContentNode = &pageState{}

var trimCutsetDotSlashSpace = func(r rune) bool {
	return r == '.' || r == '/' || unicode.IsSpace(r)
}

type contentMapConfig struct {
	lang                 string
	taxonomyConfig       taxonomiesConfigValues
	taxonomyDisabled     bool
	taxonomyTermDisabled bool
	pageDisabled         bool
	isRebuild            bool
}

var _ contentNode = (*resourceSource)(nil)

type resourceSourceState int

const (
	resourceStateNew resourceSourceState = iota
	resourceStateAssigned
)

type resourceSource struct {
	state  resourceSourceState
	sv     sitesmatrix.Vector
	path   *paths.Path
	opener hugio.OpenReadSeekCloser
	fi     hugofs.FileMetaInfo
	rc     *pagemeta.ResourceConfig

	r resource.Resource
}

func (r *resourceSource) assignSiteVector(vec sitesmatrix.Vector) *resourceSource {
	if r.state == resourceStateAssigned {
		panic("cannot assign site vector to a resourceSource that is already assigned")
	}
	r.sv = vec
	r.state = resourceStateAssigned
	return r
}

func (r resourceSource) clone() *resourceSource {
	r.state = resourceStateNew
	r.r = nil
	return &r
}

func (r *resourceSource) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	return f(r.sv, r)
}

func (r *resourceSource) String() string {
	var sb strings.Builder
	if r.fi != nil {
		sb.WriteString("filename: " + r.fi.Meta().Filename)
		sb.WriteString(fmt.Sprintf(" matrix: %v", r.fi.Meta().SitesMatrix))
		sb.WriteString(fmt.Sprintf(" fallbacks: %v", r.fi.Meta().SitesFallbacks))
	}
	if r.rc != nil {
		sb.WriteString(fmt.Sprintf("rc matrix: %v", r.rc.SitesMatrix))
		sb.WriteString(fmt.Sprintf("rc fallbacks: %v", r.rc.SitesFallbacks))
	}
	sb.WriteString(fmt.Sprintf(" sv: %v", r.sv))
	return sb.String()
}

// TODO1 remove.
func (r *resourceSource) sitesMatrix() sitesmatrix.VectorProvider {
	return r.sv
}

func (r *resourceSource) siteVector() sitesmatrix.Vector {
	return r.sv
}

func (r *resourceSource) MarkStale() {
	resource.MarkStale(r.r)
}

func (r *resourceSource) resetBuildState() {
	if rr, ok := r.r.(buildStateReseter); ok {
		rr.resetBuildState()
	}
}

func (r *resourceSource) isPage() bool {
	_, ok := r.r.(page.Page)
	return ok
}

func (r *resourceSource) contentWeight() int {
	return 0
}

func (r *resourceSource) GetIdentity() identity.Identity {
	if r.r != nil {
		return r.r.(identity.IdentityProvider).GetIdentity()
	}
	return r.path
}

func (p *resourceSource) matchSiteVector(siteVector sitesmatrix.Vector) bool {
	if p.state >= resourceStateAssigned && p.sv == siteVector {
		return true
	}

	// A site has not been assigned yet.

	if p.rc != nil && p.rc.MatchSiteVector(siteVector) {
		return true
	}

	if p.rc != nil && p.rc.SitesMatrix.LenVectors() > 0 {
		// Do not consider file mount matrix if the resource config has its own.
		return false
	}

	if p.fi != nil && p.fi.Meta().SitesMatrix.HasAnyVector(siteVector) {
		return true
	}

	return false
}

func (p *resourceSource) matchSiteVectorAll(siteVector sitesmatrix.Vector, fallback bool) iter.Seq[contentNodeForSite] {
	if siteVector == p.sv {
		return func(yield func(n contentNodeForSite) bool) {
			yield(p)
		}
	}

	/*if variant := p.getVariant(siteVector); variant != nil {
	// TODO1 remove all the variant stuff.
		return func(yield func(n contentNodeForSite) bool) {
			yield(variant)
		}
	}*/

	pc := p.rc

	var found bool
	if !fallback {
		if pc != nil && pc.MatchSiteVector(siteVector) {
			found = true
		} else {
			return nil
		}
	}

	if !found && pc != nil {
		if !pc.MatchLanguageOrLanguageFallback(siteVector) {
			return nil
		}
		if !pc.MatchVersionOrVersionFallback(siteVector) {
			return nil
		}
		if !pc.MatchRoleOrRoleFallback(siteVector) {
			return nil
		}
	}

	if !found && !fallback {
		return nil
	}

	return func(yield func(n contentNodeForSite) bool) {
		if !yield(p) {
			return
		}
	}
}

func (r *resourceSource) ForEeachIdentity(f func(identity.Identity) bool) bool {
	return f(r.GetIdentity())
}

func (r *resourceSource) Path() string {
	return r.path.Base()
}

func (r *resourceSource) isContentNodeBranch() bool {
	return false
}

var _ contentNode = (*resourceSources)(nil)

type resourceSources map[sitesmatrix.Vector]*resourceSource

func (n resourceSources) MarkStale() {
	for _, r := range n {
		if r != nil {
			r.MarkStale()
		}
	}
}

func (r resourceSources) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for _, rs := range r {
		if !f(rs.sv, rs) {
			return false
		}
	}
	return true
}

func (r resourceSources) contentWeight() int {
	return 0
}

func (n resourceSources) Path() string {
	panic("not supported")
}

func (n resourceSources) isContentNodeBranch() bool {
	return false
}

func (n resourceSources) resetBuildState() {
	for _, r := range n {
		if r != nil {
			r.resetBuildState()
		}
	}
}

func (n resourceSources) sitesMatrix() sitesmatrix.VectorProvider {
	panic("not supported")
}

func (n resourceSources) GetIdentity() identity.Identity {
	for _, r := range n {
		if r != nil {
			return r.GetIdentity()
		}
	}
	return nil
}

func (n resourceSources) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, r := range n {
		if r != nil {
			if f(r.GetIdentity()) {
				return true
			}
		}
	}
	return false
}

type pageMetaSourcesSlice []contentNode // TODO1 pool?

func (m *pageMetaSourcesSlice) nodeCategoryPage() {
	// Marker method.
}

func (n pageMetaSourcesSlice) MarkStale() {
	panic("not supported")
}

func (n pageMetaSourcesSlice) contentWeight() int {
	return 0
}

func (n pageMetaSourcesSlice) Path() string {
	return n.one().Path()
}

func (n pageMetaSourcesSlice) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for _, rs := range n {
		if !f(rs.sitesMatrix().FirstVector(), rs) {
			return false
		}
	}
	return true
}

func (n pageMetaSourcesSlice) isContentNodeBranch() bool {
	return false
}

func (n pageMetaSourcesSlice) resetBuildState() {
	panic("not supported")
}

func (n pageMetaSourcesSlice) sitesMatrix() sitesmatrix.VectorProvider {
	panic("not supported")
}

func (n pageMetaSourcesSlice) one() contentNode {
	if len(n) == 0 {
		panic("pageMetaSourcesSlice is empty")
	}
	return n[0]
}

func (n pageMetaSourcesSlice) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
}

func (n pageMetaSourcesSlice) ForEeachIdentity(f func(identity.Identity) bool) bool {
	panic("not supported")
}

type resourceSourcesSlice []*resourceSource

func (n resourceSourcesSlice) MarkStale() {
	for _, r := range n {
		if r != nil {
			r.MarkStale()
		}
	}
}

func (r resourceSourcesSlice) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for _, rs := range r {
		if !f(rs.sv, rs) {
			return false
		}
	}
	return true
}

func (n resourceSourcesSlice) ForEeachInAllDimensions(f func(contentNode) bool) {
	for _, nn := range n {
		if f(nn) {
			return
		}
	}
}

func (n resourceSourcesSlice) one() *resourceSource {
	if len(n) == 0 {
		panic("resourceSourcesSlice is empty")
	}
	return n[0]
}

func (n resourceSourcesSlice) Path() string {
	return n.one().Path()
}

func (n resourceSourcesSlice) isContentNodeBranch() bool {
	return false
}

func (m resourceSourcesSlice) contentWeight() int {
	return 0
}

func (n resourceSourcesSlice) resetBuildState() {
	for _, r := range n {
		if r != nil {
			r.resetBuildState()
		}
	}
}

func (n resourceSourcesSlice) sitesMatrix() sitesmatrix.VectorProvider {
	panic("not supported")
}

func (n resourceSourcesSlice) GetIdentity() identity.Identity {
	for _, r := range n {
		if r != nil {
			return r.GetIdentity()
		}
	}
	return nil
}

func (n resourceSourcesSlice) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, r := range n {
		if r != nil {
			if f(r.GetIdentity()) {
				return true
			}
		}
	}
	return false
}

func (cfg contentMapConfig) getTaxonomyConfig(s string) (v viewName) {
	for _, n := range cfg.taxonomyConfig.views {
		if strings.HasPrefix(s, n.pluralTreeKey) {
			return n
		}
	}
	return
}

func (m *pageMap) insertPageMetaSourceWithLock(p *pageMetaSource) (contentNode, contentNode, bool) {
	u, n, replaced := m.treePages.InsertIntoValuesDimensionWithLock(p.pathInfo.Base(), p)

	// TODO1
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

	return u, n, replaced
}

func (m *pageMap) insertResourceWithLock(s string, r contentNode) (contentNode, contentNode, bool) {
	u, n, replaced := m.treeResources.InsertIntoValuesDimensionWithLock(s, r)
	if replaced {
		// TODO1
		m.handleDuplicateResourcePath(s, r, n)
	}
	return u, n, replaced
}

func (m *pageMap) insertResource(s string, r contentNode) (contentNode, contentNode, bool) {
	u, n, replaced := m.treeResources.InsertIntoValuesDimension(s, r)
	if replaced {
		m.handleDuplicateResourcePath(s, r, n)
	}
	return u, n, replaced
}

func (m *pageMap) handleDuplicateResourcePath(s string, updated, existing contentNode) {
	if m.s.h.isRebuild() || !m.s.conf.PrintPathWarnings {
		return
	}
	var messageDetail string
	if r1, ok := existing.(*resourceSource); ok && r1.fi != nil {
		messageDetail = fmt.Sprintf(" file: %q", r1.fi.Meta().Filename)
	}
	if r2, ok := updated.(*resourceSource); ok && r2.fi != nil {
		messageDetail += fmt.Sprintf(" file: %q", r2.fi.Meta().Filename)
	}

	m.s.Log.Warnf("Duplicate resource path: %q%s", s, messageDetail)
}

func (m *pageMap) AddFi(fi hugofs.FileMetaInfo, buildConfig *BuildCfg) (pageCount uint64, resourceCount uint64, addErr error) {
	if fi.IsDir() {
		return
	}

	if m == nil {
		panic("nil pageMap")
	}

	h := m.s.h

	insertResource := func(fim hugofs.FileMetaInfo) error {
		resourceCount++
		pi := fi.Meta().PathInfo
		key := pi.Base()
		tree := m.treeResources

		commit := tree.Lock(true)
		defer commit()

		if pi.IsContent() {
			pm, err := h.newPageMetaSourceFromFile(fi)
			if err != nil {
				return fmt.Errorf("failed to create page from file %q: %w", fi.Meta().Filename, err)
			}
			pm.bundled = true
			_, _, _ = m.insertResource(key, pm)
		} else {
			r := func() (hugio.ReadSeekCloser, error) {
				return fim.Meta().Open()
			}

			// TODO1 siteVector vs matrix.
			rs := &resourceSource{path: pi, opener: r, fi: fim, sv: fim.Meta().SitesMatrix.FirstVector()}
			_, _, _ = m.insertResource(key, rs)
		}

		return nil
	}

	meta := fi.Meta()
	pi := meta.PathInfo

	switch pi.Type() {
	case paths.TypeFile, paths.TypeContentResource:
		m.s.Log.Trace(logg.StringFunc(
			func() string {
				return fmt.Sprintf("insert resource: %q", fi.Meta().Filename)
			},
		))
		if err := insertResource(fi); err != nil {
			addErr = err
			return
		}
	case paths.TypeContentData:
		pc, rc, err := m.addPagesFromGoTmplFi(fi, buildConfig)
		pageCount += pc
		resourceCount += rc
		if err != nil {
			addErr = err
			return
		}

	default:
		m.s.Log.Trace(logg.StringFunc(
			func() string {
				return fmt.Sprintf("insert bundle: %q", fi.Meta().Filename)
			},
		))

		pageCount++ // TODO1 this vs dims.

		pm, err := h.newPageMetaSourceFromFile(fi)
		if err != nil {
			addErr = fmt.Errorf("failed to create page meta from file %q: %w", fi.Meta().Filename, err)
			return
		}

		m.insertPageMetaSourceWithLock(pm)
	}
	return
}

func (m *pageMap) addPagesFromGoTmplFi(fi hugofs.FileMetaInfo, buildConfig *BuildCfg) (pageCount uint64, resourceCount uint64, addErr error) {
	meta := fi.Meta()
	pi := meta.PathInfo

	m.s.Log.Trace(logg.StringFunc(
		func() string {
			return fmt.Sprintf("insert pages from data file: %q", fi.Meta().Filename)
		},
	))

	if !files.IsGoTmplExt(pi.Ext()) {
		addErr = fmt.Errorf("unsupported data file extension %q", pi.Ext())
		return
	}

	sitesMatrix := fi.Meta().SitesMatrix

	s := m.s.h.resolveFirstSite(sitesMatrix)
	if s == nil {
		panic("TODO1")
	}
	h := s.h

	contentAdapter := s.pageMap.treePagesFromTemplateAdapters.Get(pi.Base())
	var rebuild bool
	if contentAdapter != nil {
		// Rebuild
		contentAdapter = contentAdapter.CloneForGoTmpl(fi)
		rebuild = true
	} else {
		contentAdapter = pagesfromdata.NewPagesFromTemplate(
			pagesfromdata.PagesFromTemplateOptions{
				GoTmplFi: fi,
				Site:     s,
				DepsFromSite: func(s page.Site) pagesfromdata.PagesFromTemplateDeps {
					ss := s.(*Site)
					return pagesfromdata.PagesFromTemplateDeps{
						TemplateStore: ss.GetTemplateStore(),
					}
				},
				DependencyManager: s.Conf.NewIdentityManager("pagesfromdata"),
				Watching:          s.Conf.Watching(),
				HandlePage: func(pt *pagesfromdata.PagesFromTemplate, pc *pagemeta.PageConfig) error {
					defer herrors.Recover()

					s := pt.Site.(*Site)
					if err := pc.CompileForPagesFromDataPre(pt.GoTmplFi.Meta().PathInfo.Base(), m.s.Log, s.conf.MediaTypes.Config); err != nil {
						return err
					}

					sitesMatrixFile := fi.Meta().SitesMatrix

					if !sitesMatrixFile.HasLanguage(s.siteVector.Language()) {
						sitesMatrixFile = sitesMatrixFile.WithLanguageIndex(s.siteVector.Language())
					}

					ps, err := s.h.newPageMetaSourceForContentAdapter(fi, sitesMatrixFile, pc) // TODO1 fi vs f one instance?
					if err != nil {
						return err
					}

					if ps == nil {
						// Disabled page.
						return nil
					}

					// TODO1 replace logic.

					u, n, replaced := s.pageMap.insertPageMetaSourceWithLock(ps)

					if h.isRebuild() {
						if replaced {
							pt.AddChange(n.GetIdentity())
						} else {
							pt.AddChange(u.GetIdentity())
							// New content not in use anywhere.
							// To make sure that these gets listed in any site.RegularPages ranges or similar
							// we could invalidate everything, but first try to collect a sample set
							// from the surrounding pages.
							var surroundingIDs []identity.Identity
							ids := h.pageTrees.collectIdentitiesSurrounding(pi.Base(), 10)
							if len(ids) > 0 {
								surroundingIDs = append(surroundingIDs, ids...)
							} else {
								// No surrounding pages found, so invalidate everything.
								surroundingIDs = []identity.Identity{identity.GenghisKhan}
							}
							for _, id := range surroundingIDs {
								pt.AddChange(id)
							}
						}
					}

					return nil
				},
				HandleResource: func(pt *pagesfromdata.PagesFromTemplate, rc *pagemeta.ResourceConfig) error {
					// TODO1 we need to somehow take the current Site into account here, not sure how.
					s := pt.Site.(*Site)
					if err := rc.Compile(
						pt.GoTmplFi.Meta().PathInfo.Base(),
						pt.GoTmplFi,
						s.Conf,
						s.conf.MediaTypes.Config,
					); err != nil {
						return err
					}

					// TODO1 dims vs ResourceConfig and file. Make this behave like page front matter.
					// Also check the other place where resources are added.
					// Assign this resource to the first site found. It will be lazily copied to the other sites
					// when needed.
					siteVector := rc.SitesMatrix.FirstVector()
					rs := &resourceSource{path: rc.PathInfo, rc: rc, opener: nil, fi: pt.GoTmplFi, sv: siteVector}

					_, n, replaced := s.pageMap.insertResourceWithLock(rc.PathInfo.Base(), rs)

					if h.isRebuild() && replaced {
						pt.AddChange(n.GetIdentity())
					}
					return nil
				},
			},
		)

		s.pageMap.treePagesFromTemplateAdapters.Insert(pi.Base(), contentAdapter)

	}

	handleBuildInfo := func(s *Site, bi pagesfromdata.BuildInfo) {
		resourceCount += bi.NumResourcesAdded
		pageCount += bi.NumPagesAdded
		s.handleContentAdapterChanges(bi, buildConfig)
	}

	bi, err := contentAdapter.Execute(context.Background())
	if err != nil {
		addErr = err
		return
	}
	handleBuildInfo(s, bi)

	if !rebuild && (bi.EnableAllLanguages || bi.EnableAllDimensions) {
		// Clone and insert the adapter for the other sites.
		iter := h.allSites()
		if bi.EnableAllLanguages {
			skio := func(ss *Site) bool {
				return s == ss
			}
			iter = h.allSiteLanguages(skio)
		}

		for ss := range iter {
			if s == ss {
				continue
			}

			clone := contentAdapter.CloneForSite(ss)

			// Make sure it gets executed for the first time.
			bi, err := clone.Execute(context.Background())
			if err != nil {

				addErr = err
				return
			}
			handleBuildInfo(ss, bi)

			// Insert into the correct language tree so it get rebuilt on changes.
			ss.pageMap.treePagesFromTemplateAdapters.Insert(pi.Base(), clone)

		}
	}

	return
}

// The home page is represented with the zero string.
// All other keys starts with a leading slash. No trailing slash.
// Slashes are Unix-style.
func cleanTreeKey(elem ...string) string {
	var s string
	if len(elem) > 0 {
		s = elem[0]
		if len(elem) > 1 {
			s = path.Join(elem...)
		}
	}
	s = strings.TrimFunc(s, trimCutsetDotSlashSpace)
	s = filepath.ToSlash(strings.ToLower(paths.Sanitize(s)))
	if s == "" || s == "/" {
		return ""
	}
	if s[0] != '/' {
		s = "/" + s
	}
	return s
}
