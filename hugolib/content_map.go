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
	"unicode"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/resources/page"
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
	path   *paths.Path
	opener hugio.OpenReadSeekCloser
	fi     hugofs.FileMetaInfo

	r resource.Resource
}

func (r resourceSource) clone() *resourceSource {
	r.r = nil
	return &r
}

func (r *resourceSource) LangIndex() int {
	if r.r != nil && r.isPage() {
		return r.r.(*pageState).s.languagei
	}

	return r.fi.Meta().LangIndex
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

// TODO1 move this somewhere
var skipc = maps.NewCache[string, map[uint64]bool]()

func (m *pageMap) AddFi(fi hugofs.FileMetaInfo, whatChanged *whatChanged) error {
	if fi.IsDir() {
		return nil
	}

	rebuilding := m.s.h.buildCounter.Load() > 0

	insertResource := func(fim hugofs.FileMetaInfo) error {
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
			// Create the page now as we need it at assemembly time.
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

			rs = &resourceSource{r: pageResource}
		} else {
			rs = &resourceSource{path: pi, opener: r, fi: fim}
		}

		tree.InsertIntoValuesDimension(key, rs)

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
			return err
		}
	case paths.PathTypeContentData:
		m.s.Log.Trace(logg.StringFunc(
			func() string {
				return fmt.Sprintf("insert pages from data file: %q", fi.Meta().Filename)
			},
		))

		if err := func() error {
			// TODO1 validate no leading slash in v1.
			// TODO1 benchmark this vs a more fine grained lock.
			commit1, commit2 := m.treePages.Lock(true), m.treeResources.Lock(true)
			defer commit1()
			defer commit2()

			seen := make(map[uint64]bool)
			handle := func(p pagesfromdata.PageData, skip bool) error {
				seen[p.SourceHash] = true
				if skip {
					return nil
				}
				pc := p.PageConfig
				pc.Path = path.Join(pi.Base(), pc.Path)
				if err := pc.Compile(true); err != nil {
					return err
				}
				ps, pi, err := m.s.h.newPage(
					&pageMeta{
						f: source.NewFileInfo(fi),
						pageMetaParams: pageMetaParams{
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

				n, _, replaced := m.treePages.InsertIntoValuesDimension(pi.Base(), ps)

				if rebuilding && replaced {
					whatChanged.Add(n.GetIdentity())
				}
				return nil
			}

			skip, _ := skipc.Get(fi.Meta().Filename)
			if err := pagesfromdata.PagesFromJSONFile(fi, skip, handle); err != nil {
				return err
			}

			skipc.Set(fi.Meta().Filename, seen)

			return nil
		}(); err != nil {
			return err
		}
	default:
		m.s.Log.Trace(logg.StringFunc(
			func() string {
				return fmt.Sprintf("insert bundle: %q", fi.Meta().Filename)
			},
		))
		// A content file.
		p, pi, err := m.s.h.newPage(
			&pageMeta{
				f:        source.NewFileInfo(fi),
				pathInfo: pi,
				bundled:  false,
			},
		)
		if err != nil {
			return err
		}
		if p == nil {
			// Disabled page.
			return nil
		}

		m.treePages.InsertIntoValuesDimensionWithLock(pi.Base(), p)

	}
	return nil
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
