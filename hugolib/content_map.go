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
	"unicode"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/source"

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

var _ contentNodeI = (*resourceSource)(nil)

type resourceSource struct {
	langIndex int
	path      *paths.Path
	opener    hugio.OpenReadSeekCloser
	fi        hugofs.FileMetaInfo
	rc        *pagemeta.ResourceConfig

	r resource.Resource
}

func (r resourceSource) clone() *resourceSource {
	r.r = nil
	return &r
}

func (r *resourceSource) LangIndex() int {
	return r.langIndex
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

func (r *resourceSource) GetIdentity() identity.Identity {
	if r.r != nil {
		return r.r.(identity.IdentityProvider).GetIdentity()
	}
	return r.path
}

func (r *resourceSource) ForEeachIdentity(f func(identity.Identity) bool) bool {
	return f(r.GetIdentity())
}

func (r *resourceSource) Path() string {
	return r.path.Path()
}

func (r *resourceSource) isContentNodeBranch() bool {
	return false
}

var _ contentNodeI = (*resourceSources)(nil)

type resourceSources []*resourceSource

func (n resourceSources) MarkStale() {
	for _, r := range n {
		if r != nil {
			r.MarkStale()
		}
	}
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

func (cfg contentMapConfig) getTaxonomyConfig(s string) (v viewName) {
	for _, n := range cfg.taxonomyConfig.views {
		if strings.HasPrefix(s, n.pluralTreeKey) {
			return n
		}
	}
	return
}

func (m *pageMap) insertPageWithLock(s string, p *pageState) (contentNodeI, contentNodeI, bool) {
	u, n, replaced := m.treePages.InsertIntoValuesDimensionWithLock(s, p)

	if replaced && !m.s.h.isRebuild() && m.s.conf.PrintPathWarnings {
		var messageDetail string
		if p1, ok := n.(*pageState); ok && p1.File() != nil {
			messageDetail = fmt.Sprintf(" file: %q", p1.File().Filename())
		}
		if p2, ok := u.(*pageState); ok && p2.File() != nil {
			messageDetail += fmt.Sprintf(" file: %q", p2.File().Filename())
		}

		m.s.Log.Warnf("Duplicate content path: %q%s", s, messageDetail)
	}

	return u, n, replaced
}

func (m *pageMap) insertResourceWithLock(s string, r contentNodeI) (contentNodeI, contentNodeI, bool) {
	u, n, replaced := m.treeResources.InsertIntoValuesDimensionWithLock(s, r)
	if replaced {
		m.handleDuplicateResourcePath(s, r, n)
	}
	return u, n, replaced
}

func (m *pageMap) insertResource(s string, r contentNodeI) (contentNodeI, contentNodeI, bool) {
	u, n, replaced := m.treeResources.InsertIntoValuesDimension(s, r)
	if replaced {
		m.handleDuplicateResourcePath(s, r, n)
	}
	return u, n, replaced
}

func (m *pageMap) handleDuplicateResourcePath(s string, updated, existing contentNodeI) {
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

	insertResource := func(fim hugofs.FileMetaInfo) error {
		resourceCount++
		pi := fi.Meta().PathInfo
		key := pi.Base()
		tree := m.treeResources

		commit := tree.Lock(true)
		defer commit()

		r := func() (hugio.ReadSeekCloser, error) {
			return fim.Meta().Open()
		}

		var rs *resourceSource
		if pi.IsContent() {
			// Create the page now as we need it at assembly time.
			// The other resources are created if needed.
			pageResource, pi, err := m.s.h.newPage(
				&pageMeta{
					f:        source.NewFileInfo(fim),
					pathInfo: pi,
					bundled:  true,
				},
			)
			if err != nil {
				return err
			}
			if pageResource == nil {
				// Disabled page.
				return nil
			}
			key = pi.Base()

			rs = &resourceSource{r: pageResource, langIndex: pageResource.s.languagei}
		} else {
			rs = &resourceSource{path: pi, opener: r, fi: fim, langIndex: fim.Meta().LangIndex}
		}

		_, _, _ = m.insertResource(key, rs)

		return nil
	}

	meta := fi.Meta()
	pi := meta.PathInfo

	switch pi.BundleType() {
	case paths.PathTypeFile, paths.PathTypeContentResource:
		m.s.Log.Trace(logg.StringFunc(
			func() string {
				return fmt.Sprintf("insert resource: %q", fi.Meta().Filename)
			},
		))
		if err := insertResource(fi); err != nil {
			addErr = err
			return
		}
	case paths.PathTypeContentData:
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

		pageCount++

		// A content file.
		p, pi, err := m.s.h.newPage(
			&pageMeta{
				f:        source.NewFileInfo(fi),
				pathInfo: pi,
				bundled:  false,
			},
		)
		if err != nil {
			addErr = err
			return
		}
		if p == nil {
			// Disabled page.
			return
		}

		m.insertPageWithLock(pi.Base(), p)

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

	s := m.s.h.resolveSite(fi.Meta().Lang)
	f := source.NewFileInfo(fi)
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
						TmplFinder: ss.TextTmpl(),
						TmplExec:   ss.Tmpl(),
					}
				},
				DependencyManager: s.Conf.NewIdentityManager("pagesfromdata"),
				Watching:          s.Conf.Watching(),
				HandlePage: func(pt *pagesfromdata.PagesFromTemplate, pc *pagemeta.PageConfig) error {
					s := pt.Site.(*Site)
					if err := pc.Compile(pt.GoTmplFi.Meta().PathInfo.Base(), true, "", s.Log, s.conf.MediaTypes.Config); err != nil {
						return err
					}

					ps, pi, err := h.newPage(
						&pageMeta{
							f: f,
							s: s,
							pageMetaParams: &pageMetaParams{
								pageConfig: pc,
							},
						},
					)
					if err != nil {
						return err
					}

					if ps == nil {
						// Disabled page.
						return nil
					}

					u, n, replaced := s.pageMap.insertPageWithLock(pi.Base(), ps)

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
					s := pt.Site.(*Site)
					if err := rc.Compile(
						pt.GoTmplFi.Meta().PathInfo.Base(),
						s.Conf.PathParser(),
						s.conf.MediaTypes.Config,
					); err != nil {
						return err
					}

					rs := &resourceSource{path: rc.PathInfo, rc: rc, opener: nil, fi: pt.GoTmplFi, langIndex: s.languagei}

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

	if !rebuild && bi.EnableAllLanguages {
		// Clone and insert the adapter for the other sites.
		for _, ss := range s.h.Sites {
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
