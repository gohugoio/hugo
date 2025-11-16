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

func (r *resourceSource) nodeCategorySingle() {
	// Marker method.
}

func (r *resourceSource) String() string {
	var sb strings.Builder
	if r.fi != nil {
		sb.WriteString("filename: " + r.fi.Meta().Filename)
		sb.WriteString(fmt.Sprintf(" matrix: %v", r.fi.Meta().SitesMatrix))
		sb.WriteString(fmt.Sprintf(" complements: %v", r.fi.Meta().SitesComplements))
	}
	if r.rc != nil {
		sb.WriteString(fmt.Sprintf("rc matrix: %v", r.rc.SitesMatrix))
		sb.WriteString(fmt.Sprintf("rc complements: %v", r.rc.SitesComplements))
	}
	sb.WriteString(fmt.Sprintf(" sv: %v", r.sv))
	return sb.String()
}

func (r *resourceSource) siteVector() sitesmatrix.Vector {
	return r.sv
}

func (r *resourceSource) MarkStale() {
	resource.MarkStale(r.r)
}

func (r *resourceSource) resetBuildState() {
	if rr, ok := r.r.(contentNodeBuildStateResetter); ok {
		rr.resetBuildState()
	}
}

func (r *resourceSource) GetIdentity() identity.Identity {
	if r.r != nil {
		return r.r.(identity.IdentityProvider).GetIdentity()
	}
	return r.path
}

func (p *resourceSource) nodeSourceEntryID() any {
	if p.rc != nil {
		return p.rc.ContentAdapterSourceEntryHash
	}
	if p.fi != nil {
		return p.fi.Meta().Filename
	}
	return p.path
}

func (p *resourceSource) lookupContentNode(v sitesmatrix.Vector) contentNode {
	if p.state >= resourceStateAssigned {
		if p.sv == v {
			return p
		}
		return nil
	}

	// A site has not been assigned yet.

	if p.rc != nil && p.rc.MatchSiteVector(v) {
		return p
	}

	if p.rc != nil && p.rc.SitesMatrix.LenVectors() > 0 {
		// Do not consider file mount matrix if the resource config has its own.
		return nil
	}

	if p.fi != nil && p.fi.Meta().SitesMatrix.HasVector(v) {
		return p
	}

	return nil
}

func (p *resourceSource) lookupContentNodes(siteVector sitesmatrix.Vector, fallback bool) iter.Seq[contentNodeForSite] {
	if siteVector == p.sv {
		return func(yield func(n contentNodeForSite) bool) {
			yield(p)
		}
	}

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
		if !pc.MatchLanguageCoarse(siteVector) {
			return nil
		}
		if !pc.MatchVersionCoarse(siteVector) {
			return nil
		}
		if !pc.MatchRoleCoarse(siteVector) {
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

func (r *resourceSource) Path() string {
	return r.path.Base()
}

func (r *resourceSource) PathInfo() *paths.Path {
	return r.path
}

func (cfg contentMapConfig) getTaxonomyConfig(s string) (v viewName) {
	for _, n := range cfg.taxonomyConfig.views {
		if strings.HasPrefix(s, n.pluralTreeKey) {
			return n
		}
	}
	return
}

func (m *pageMap) AddFi(fi hugofs.FileMetaInfo, buildConfig *BuildCfg) (pageSourceCount uint64, resourceSourceCount uint64, addErr error) {
	if fi.IsDir() {
		return
	}

	if m == nil {
		panic("nil pageMap")
	}

	h := m.s.h

	insertResource := func(fim hugofs.FileMetaInfo) error {
		resourceSourceCount++
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
			m.treeResources.Insert(key, pm)
		} else {
			r := func() (hugio.ReadSeekCloser, error) {
				return fim.Meta().Open()
			}
			// Create one dimension now, the rest later on demand.
			rs := &resourceSource{path: pi, opener: r, fi: fim, sv: fim.Meta().SitesMatrix.VectorSample()}
			m.treeResources.Insert(key, rs)
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
		pageSourceCount += pc
		resourceSourceCount += rc
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

		pageSourceCount++

		pm, err := h.newPageMetaSourceFromFile(fi)
		if err != nil {
			addErr = fmt.Errorf("failed to create page meta from file %q: %w", fi.Meta().Filename, err)
			return
		}

		m.treePages.InsertWithLock(pm.pathInfo.Base(), pm)
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
				DependencyManager: s.Conf.NewIdentityManager(),
				Watching:          s.Conf.Watching(),
				HandlePage: func(pt *pagesfromdata.PagesFromTemplate, pe *pagemeta.PageConfigEarly) error {
					s := pt.Site.(*Site)

					if err := pe.CompileForPagesFromDataPre(pt.GoTmplFi.Meta().PathInfo.Base(), m.s.Log, s.conf.MediaTypes.Config); err != nil {
						return err
					}

					ps, err := s.h.newPageMetaSourceForContentAdapter(fi, s.siteVector, pe)
					if err != nil {
						return err
					}

					if ps == nil {
						// Disabled page.
						return nil
					}

					u, n, replaced := s.pageMap.treePages.InsertWithLock(ps.pathInfo.Base(), ps)

					if h.isRebuild() {
						if replaced {
							pt.AddChange(cnh.GetIdentity(n))
						} else {
							pt.AddChange(cnh.GetIdentity(u))
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
					s := pt.Site.(*Site)
					if err := rc.Compile(
						pt.GoTmplFi.Meta().PathInfo.Base(),
						pt.GoTmplFi,
						s.Conf,
						s.conf.MediaTypes.Config,
					); err != nil {
						return err
					}

					// Create one dimension now, the rest later on demand.
					rs := &resourceSource{path: rc.PathInfo, rc: rc, opener: nil, fi: pt.GoTmplFi, sv: s.siteVector}

					_, n, updated := s.pageMap.treeResources.InsertWithLock(rs.path.Base(), rs)

					if h.isRebuild() && updated {
						pt.AddChange(cnh.GetIdentity(n))
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
		var iter iter.Seq[*Site]
		if bi.EnableAllLanguages {
			include := func(ss *Site) bool {
				return s.siteVector.Language() != ss.siteVector.Language()
			}
			iter = h.allSiteLanguages(include)
		} else {
			include := func(ss *Site) bool {
				return s.siteVector != ss.siteVector
			}
			iter = h.allSites(include)
		}

		for ss := range iter {
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
