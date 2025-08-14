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
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bep/logg"
	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/resources"
	xmaps "golang.org/x/exp/maps"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hdebug"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/cast"
)

var cjkRe = regexp.MustCompile(`\p{Han}|\p{Hangul}|\p{Hiragana}|\p{Katakana}`)

// Implement contentNodeI
func (m *pageMetaSource) GetIdentity() identity.Identity {
	return m.pathInfo
}

// ForEeachIdentityProvider calls cb for each Identity.
// If cb returns true, the iteration is terminated.
// The return value is whether the iteration was terminated.
func (m *pageMetaSource) ForEeachIdentity(cb func(id identity.Identity) bool) bool {
	panic("not implemented") // TODO: Implement
}

func (m *pageMetaSource) Path() string {
	return m.pathInfo.Base()
}

func (m *pageMetaSource) isContentNodeBranch() bool {
	panic("not implemented") // TODO: Implement
}

func (m *pageMetaSource) contentWeight() int {
	panic("not implemented") // TODO: Implement
}

func (m *pageMetaSource) sitesMatrix() sitesmatrix.VectorProvider {
	return m.pageConfigSource.SitesMatrix
}

func (m *pageMetaSource) resetBuildState() {
	panic("not implemented") // TODO: Implement
}

func (m *pageMetaSource) MarkStale() {
	panic("not implemented") // TODO: Implement
}

type pageMetaSource struct {
	pathInfo *paths.Path // Always set. This the canonical path to the Page. // TODO1 remove.
	f        *source.File
	pi       *contentParseInfo
	resource.Staler

	pageConfigSource *pagemeta.PageConfig

	resourcePath  string // Set for bundled pages; path relative to its bundle root.
	bundled       bool   // Set if this page is bundled inside another.
	noFrontMatter bool
	openSource    func() (hugio.ReadSeekCloser, error)

	initEarlyInit sync.Once
}

func (m *pageMetaSource) String() string {
	if m.f != nil {
		return fmt.Sprintf("pageMetaSource(%s, %s)", m.f.FileInfo().Meta().PathInfo, m.f.FileInfo().Meta().Filename)
	}
	return fmt.Sprintf("pageMetaSource(%s)", m.pathInfo)
}

type pageMeta struct {
	// Shared between all dimensions of this page.
	*pageMetaSource

	pageConfig *pagemeta.PageConfig

	// Per Page. TODO1, potentially share some.
	*pageMetaParams

	term     string // Set for kind == KindTerm.
	singular string // Set for kind == KindTerm and kind == KindTaxonomy.

	resource.Staler // TODO1 remove?

	// Set for standalone pages, e.g. robotsTXT.
	standaloneOutputFormat output.Format

	content *cachedContent // The source and the parsed page content.
}

func (m *pageMetaSource) initPathInfo(h *HugoSites) error {
	pcfg := m.pageConfigSource
	if pcfg.Path != "" {
		s := pcfg.Path
		// Paths from content adapters should never have any extension.
		if pcfg.IsFromContentAdapter || !paths.HasExt(s) {
			var (
				isBranch    bool
				isBranchSet bool
				ext         string = pcfg.ContentMediaType.FirstSuffix.Suffix
			)
			if pcfg.Kind != "" {
				isBranch = kinds.IsBranch(pcfg.Kind)
				isBranchSet = true
			}

			if !pcfg.IsFromContentAdapter {
				if m.pathInfo != nil {
					if !isBranchSet {
						isBranch = m.pathInfo.IsBranchBundle()
					}
					if m.pathInfo.Ext() != "" {
						ext = m.pathInfo.Ext()
					}
				} else if m.f != nil {
					pi := m.f.FileInfo().Meta().PathInfo
					if !isBranchSet {
						isBranch = pi.IsBranchBundle()
					}
					if pi.Ext() != "" {
						ext = pi.Ext()
					}
				}
			}
			if isBranch {
				s += "/_index." + ext
			} else {
				s += "/index." + ext
			}

		}
		m.pathInfo = h.Conf.PathParser().Parse(files.ComponentFolderContent, s)
	} else {
		if m.f != nil {
			m.pathInfo = m.f.FileInfo().Meta().PathInfo
		}

		if m.pathInfo == nil {
			panic(fmt.Sprintf("missing pathInfo in %v", m))
		}
	}
	return nil
}

// initEarly may be called before we know which Site(s) this belongs to.
func (m *pageMetaSource) initEarly(h *HugoSites, sitesMatrixFile sitesmatrix.VectorStore) error {
	var initErr error
	m.initEarlyInit.Do(func() {
		initErr = func() error {
			if m.pathInfo == nil {
				if err := m.initPathInfo(h); err != nil {
					return err
				}
			}
			if m.pageConfigSource == nil {
				m.pageConfigSource = &pagemeta.PageConfig{}
			}

			if m.Staler == nil {
				m.Staler = &resources.AtomicStaler{}
			}

			sid := pageSourceIDCounter.Add(1)

			// Read front matter.
			if err := m.parseFrontMatter(h, sid); err != nil {
				return err
			}

			if m.pi.frontMatter != nil {
				if err := m.pageConfigSource.SetMetaPreFromMap(m.pi.frontMatter, h.Log, h.Conf); err != nil {
					return err
				}
			} else {
				m.pageConfigSource.Params = make(maps.Params)
			}

			var fim *hugofs.FileMeta
			if m.f != nil {
				fim = m.f.FileInfo().Meta()
			}
			if err := m.pageConfigSource.CompileEearly(h.Conf, fim, sitesMatrixFile); err != nil {
				return err
			}

			return nil
		}()
	})
	return initErr
}

func (m *pageMeta) initLate(s *Site) error {
	if m.pageMetaParams == nil {
		m.pageMetaParams = &pageMetaParams{}
	}

	var sitesMatrixFile sitesmatrix.VectorStore
	if m.f != nil {
		sitesMatrixFile = m.f.FileInfo().Meta().SitesMatrix
	}
	if err := m.initEarly(s.h, sitesMatrixFile); err != nil {
		return err
	}

	if m.pageConfig == nil {
		if len(s.h.sitesVersionsRolesMap) > 0 {
			m.pageConfig = pagemeta.ClonePageConfigForSite(m.pageConfigSource)
		} else {
			m.pageConfig = m.pageConfigSource
		}
	}

	// Remove me. TODO1
	if m.Staler == nil {
		m.Staler = &resources.AtomicStaler{}
	}

	h := s.h
	var tc viewName

	if m.pageConfig.Kind == "" {
		// Resolve page kind.
		m.pageConfig.Kind = kinds.KindSection
		if m.pathInfo.Base() == "/" {
			m.pageConfig.Kind = kinds.KindHome
		} else if m.pathInfo.IsBranchBundle() {
			// A section, taxonomy or term.
			tc = s.pageMap.cfg.getTaxonomyConfig(m.Path())
			if !tc.IsZero() {
				// Either a taxonomy or a term.
				if tc.pluralTreeKey == m.Path() {
					m.pageConfig.Kind = kinds.KindTaxonomy
				} else {
					m.pageConfig.Kind = kinds.KindTerm
				}
			}
		} else if m.f != nil {
			m.pageConfig.Kind = kinds.KindPage
		}
	}

	if m.pageConfig.Kind == kinds.KindTerm || m.pageConfig.Kind == kinds.KindTaxonomy {
		if tc.IsZero() {
			tc = s.pageMap.cfg.getTaxonomyConfig(m.Path())
		}
		if tc.IsZero() {
			return fmt.Errorf("no taxonomy configuration found for %q", m.Path())
		}
		m.singular = tc.singular

		if m.pageConfig.Kind == kinds.KindTerm {
			m.term = paths.TrimLeading(strings.TrimPrefix(m.pathInfo.Unnormalized().Base(), tc.pluralTreeKey))
		}
	}

	m.initPageMetaParams(h.Conf.Watching())

	// TODO1 can we do this here?
	/*if m.pageConfig.Kind == kinds.KindPage && !m.s.conf.IsKindEnabled(m.pageConfig.Kind) {
		return nil,, nil
	}*/

	return nil
}

// bookmark1
func (h *HugoSites) newPageMetaSourceFromFile(fi hugofs.FileMetaInfo) (*pageMetaSource, error) {
	p, err := func() (*pageMetaSource, error) {
		meta := fi.Meta()
		openSource := func() (hugio.ReadSeekCloser, error) {
			r, err := meta.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file %q: %w", meta.Filename, err)
			}
			return r, nil
		}

		p := &pageMetaSource{
			f:                source.NewFileInfo(fi),
			pathInfo:         fi.Meta().PathInfo,
			openSource:       openSource,
			Staler:           &resources.AtomicStaler{},
			pageConfigSource: &pagemeta.PageConfig{},
		}

		return p, p.initEarly(h, fi.Meta().SitesMatrix)
	}()
	if err != nil {
		return nil, hugofs.AddFileInfoToError(err, fi, h.SourceFs)
	}

	return p, err
}

func (h *HugoSites) newPageMetaSourceForContentAdapter(fi hugofs.FileMetaInfo, sitesMatrixFile sitesmatrix.VectorStore, pc *pagemeta.PageConfig) (*pageMetaSource, error) {
	p := &pageMetaSource{
		f:                source.NewFileInfo(fi),
		noFrontMatter:    true,
		Staler:           &resources.AtomicStaler{},
		pageConfigSource: pc,
		openSource:       pc.Content.ValueAsOpenReadSeekCloser(),
	}

	return p, p.initEarly(h, sitesMatrixFile)
}

func (s *Site) newPageFromPageMetasource(ms *pageMetaSource) (*pageState, error) {
	m, err := s.newPageMetaFromPageMetasource(ms)
	if err != nil {
		return nil, err
	}
	return s.newPageNew(m)
}

func (s *Site) newPageMetaFromPageMetasource(ms *pageMetaSource) (*pageMeta, error) {
	m := &pageMeta{
		pageMetaSource: ms,
		Staler:         &resources.AtomicStaler{},
		pageMetaParams: &pageMetaParams{},
	}

	var sitesMatrixFile sitesmatrix.VectorStore
	if ms.f != nil {
		sitesMatrixFile = ms.f.FileInfo().Meta().SitesMatrix
	}

	return m, m.initEarly(s.h, sitesMatrixFile)
}

// Prepare for a rebuild of the data passed in from front matter.
func (m *pageMeta) setMetaPostPrepareRebuild() {
	params := xmaps.Clone(m.paramsOriginal)
	m.pageConfig = pagemeta.ClonePageConfigForRebuild(m.pageConfig, params)
}

type pageMetaParams struct {
	setMetaPostCount          int
	setMetaPostCascadeChanged bool

	// These are only set in watch mode.
	datesOriginal   pagemeta.Dates
	paramsOriginal  map[string]any                                                // contains the original params as defined in the front matter.
	cascadeOriginal *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig] // contains the original cascade as defined in the front matter.
}

func (m *pageMeta) initPageMetaParams(preserveOriginal bool) {
	if preserveOriginal {
		if m.pageConfig.IsFromContentAdapter {
			m.paramsOriginal = xmaps.Clone(m.pageConfig.ContentAdapterData)
		} else {
			m.paramsOriginal = xmaps.Clone(m.pageConfig.Params)
		}
		m.cascadeOriginal = m.pageConfig.CascadeCompiled.Clone()
	}
}

func (m *pageMeta) Aliases() []string {
	return m.pageConfig.Aliases
}

func (m *pageMeta) BundleType() string {
	switch m.pathInfo.Type() {
	case paths.TypeLeaf:
		return "leaf"
	case paths.TypeBranch:
		return "branch"
	default:
		return ""
	}
}

func (m *pageMeta) Date() time.Time {
	return m.pageConfig.Dates.Date
}

func (m *pageMeta) PublishDate() time.Time {
	return m.pageConfig.Dates.PublishDate
}

func (m *pageMeta) Lastmod() time.Time {
	return m.pageConfig.Dates.Lastmod
}

func (m *pageMeta) ExpiryDate() time.Time {
	return m.pageConfig.Dates.ExpiryDate
}

func (m *pageMeta) Description() string {
	return m.pageConfig.Description
}

func (m *pageMeta) Draft() bool {
	return m.pageConfig.Draft
}

func (m *pageMeta) File() *source.File {
	return m.f
}

func (m *pageMeta) IsHome() bool {
	return m.Kind() == kinds.KindHome
}

func (m *pageMeta) Keywords() []string {
	return m.pageConfig.Keywords
}

func (m *pageMeta) Kind() string {
	return m.pageConfig.Kind
}

func (m *pageMeta) Layout() string {
	return m.pageConfig.Layout
}

func (m *pageMeta) LinkTitle() string {
	if m.pageConfig.LinkTitle != "" {
		return m.pageConfig.LinkTitle
	}

	return m.Title()
}

func (m *pageMeta) Name() string {
	if m.resourcePath != "" {
		return m.resourcePath
	}
	if m.pageConfig.Kind == kinds.KindTerm {
		return m.pathInfo.Unnormalized().BaseNameNoIdentifier()
	}
	return m.Title()
}

func (m *pageMeta) IsNode() bool {
	return !m.IsPage()
}

func (m *pageMeta) IsPage() bool {
	return m.Kind() == kinds.KindPage
}

func (m *pageMeta) Params() maps.Params {
	return m.pageConfig.Params
}

func (m *pageMeta) Path() string {
	return m.pathInfo.Base()
}

func (m *pageMeta) PathInfo() *paths.Path {
	return m.pathInfo
}

func (m *pageMeta) IsSection() bool {
	return m.Kind() == kinds.KindSection
}

func (m *pageMeta) Section() string {
	return m.pathInfo.Section()
}

func (m *pageMeta) Sitemap() config.SitemapConfig {
	return m.pageConfig.Sitemap
}

func (m *pageMeta) Title() string {
	return m.pageConfig.Title
}

const defaultContentType = "page"

func (m *pageMeta) Type() string {
	if m.pageConfig.Type != "" {
		return m.pageConfig.Type
	}

	if sect := m.Section(); sect != "" {
		return sect
	}

	return defaultContentType
}

func (m *pageMeta) Weight() int {
	return m.pageConfig.Weight
}

func (ps *pageState) setMetaPost(cascade *maps.Ordered[page.PageMatcher, page.PageMatcherParamsConfig]) error {
	hdebug.AssertNotNil(ps.m.pageMetaParams)
	ps.m.setMetaPostCount++
	var cascadeHashPre uint64
	if ps.m.setMetaPostCount > 1 {
		cascadeHashPre = hashing.HashUint64(ps.m.pageConfig.CascadeCompiled)
		ps.m.pageConfig.CascadeCompiled = ps.m.cascadeOriginal.Clone()
	}

	// Apply cascades first so they can be overridden later.
	if cascade != nil {
		if ps.m.pageConfig.CascadeCompiled != nil {
			cascade.Range(func(k page.PageMatcher, v page.PageMatcherParamsConfig) bool {
				vv, found := ps.m.pageConfig.CascadeCompiled.Get(k)
				if !found {
					ps.m.pageConfig.CascadeCompiled.Set(k, v)
				} else {
					// Merge
					for ck, cv := range v.Params {
						if _, found := vv.Params[ck]; !found {
							vv.Params[ck] = cv
						}
					}
					for ck, cv := range v.Fields {
						if _, found := vv.Fields[ck]; !found {
							vv.Fields[ck] = cv
						}
					}
				}
				return true
			})
			cascade = ps.m.pageConfig.CascadeCompiled
		} else {
			ps.m.pageConfig.CascadeCompiled = cascade
		}
	}

	if cascade == nil {
		cascade = ps.m.pageConfig.CascadeCompiled
	}

	if ps.m.setMetaPostCount > 1 {
		ps.m.setMetaPostCascadeChanged = cascadeHashPre != hashing.HashUint64(ps.m.pageConfig.CascadeCompiled)
		if !ps.m.setMetaPostCascadeChanged {

			// No changes, restore any value that may be changed by aggregation.
			ps.m.pageConfig.Dates = ps.m.datesOriginal
			return nil
		}
		ps.m.setMetaPostPrepareRebuild()

	}

	// Cascade is also applied to itself.
	cascade.Range(func(k page.PageMatcher, v page.PageMatcherParamsConfig) bool {
		if !k.Matches(ps) {
			return true
		}
		for kk, vv := range v.Params {
			if _, found := ps.m.pageConfig.Params[kk]; !found {
				ps.m.pageConfig.Params[kk] = vv
			}
		}

		for kk, vv := range v.Fields {
			if ps.m.pageConfig.IsFromContentAdapter {
				if _, found := ps.m.pageConfig.ContentAdapterData[kk]; !found {
					ps.m.pageConfig.ContentAdapterData[kk] = vv
				}
			} else {
				if _, found := ps.m.pageConfig.Params[kk]; !found {
					ps.m.pageConfig.Params[kk] = vv
				}
			}
		}
		return true
	})

	if err := ps.setMetaPostParams(); err != nil {
		return err
	}

	if err := ps.m.applyDefaultValues(ps.s); err != nil {
		return err
	}

	// Store away any original values that may be changed from aggregation.
	ps.m.datesOriginal = ps.m.pageConfig.Dates

	return nil
}

func (ps *pageState) setMetaPostParams() error {
	pm := ps.m
	var mtime time.Time
	var contentBaseName string
	var ext string
	var isContentAdapter bool
	if ps.File() != nil {
		isContentAdapter = ps.File().IsContentAdapter()
		contentBaseName = ps.File().ContentBaseName()
		if ps.File().FileInfo() != nil {
			mtime = ps.File().FileInfo().ModTime()
		}
		if !isContentAdapter {
			ext = ps.File().Ext()
		}
	}

	var gitAuthorDate time.Time
	if ps.gitInfo != nil {
		gitAuthorDate = ps.gitInfo.AuthorDate
	}

	descriptor := &pagemeta.FrontMatterDescriptor{
		PageConfig:    pm.pageConfig,
		BaseFilename:  contentBaseName,
		ModTime:       mtime,
		GitAuthorDate: gitAuthorDate,
		Location:      langs.GetLocation(ps.s.Language()),
		PathOrTitle:   ps.pathOrTitle(),
	}

	if isContentAdapter {
		if err := pm.pageConfig.Compile(ext, ps.s.Log, ps.s.conf.OutputFormats.Config, ps.s.conf.MediaTypes.Config); err != nil {
			return err
		}
	}

	// Handle the date separately
	// TODO(bep) we need to "do more" in this area so this can be split up and
	// more easily tested without the Page, but the coupling is strong.
	err := ps.s.frontmatterHandler.HandleDates(descriptor)
	if err != nil {
		ps.s.Log.Errorf("Failed to handle dates for page %q: %s", ps.pathOrTitle(), err)
	}

	if isContentAdapter {
		// Done.
		return nil
	}

	var buildConfig any
	var isNewBuildKeyword bool
	if v, ok := pm.pageConfig.Params["_build"]; ok {
		hugo.Deprecate("The \"_build\" front matter key", "Use \"build\" instead. See https://gohugo.io/content-management/build-options.", "0.145.0")
		buildConfig = v
	} else {
		buildConfig = pm.pageConfig.Params["build"]
		isNewBuildKeyword = true
	}
	pm.pageConfig.Build, err = pagemeta.DecodeBuildConfig(buildConfig)
	if err != nil {
		var msgDetail string
		if isNewBuildKeyword {
			msgDetail = `. We renamed the _build keyword to build in Hugo 0.123.0. We recommend putting user defined params in the params section, e.g.:
---
title: "My Title"
params:
  build: "My Build"
---
´

`
		}
		return fmt.Errorf("failed to decode build config in front matter: %s%s", err, msgDetail)
	}

	var sitemapSet bool

	pcfg := pm.pageConfig
	params := pcfg.Params
	if params == nil {
		panic("params not set for " + ps.Path())
	}

	var draft, published, isCJKLanguage *bool
	var userParams map[string]any
	for k, v := range pcfg.Params {
		loki := strings.ToLower(k)

		if loki == "params" {
			vv, err := maps.ToStringMapE(v)
			if err != nil {
				return err
			}
			userParams = vv
			delete(pcfg.Params, k)
			continue
		}

		if loki == "published" { // Intentionally undocumented
			vv, err := cast.ToBoolE(v)
			if err == nil {
				published = &vv
			}
			// published may also be a date
			continue
		}

		if ps.s.frontmatterHandler.IsDateKey(loki) {
			continue
		}

		if loki == "path" || loki == "kind" || loki == "lang" {
			// See issue 12484.
			hugo.DeprecateLevelMin(loki+" in front matter", "", "v0.144.0", logg.LevelWarn)
		}

		switch loki {
		case "title":
			pcfg.Title = cast.ToString(v)
			params[loki] = pcfg.Title
		case "linktitle":
			pcfg.LinkTitle = cast.ToString(v)
			params[loki] = pcfg.LinkTitle
		case "summary":
			pcfg.Summary = cast.ToString(v)
			params[loki] = pcfg.Summary
		case "description":
			pcfg.Description = cast.ToString(v)
			params[loki] = pcfg.Description
		case "slug":
			// Don't start or end with a -
			pcfg.Slug = strings.Trim(cast.ToString(v), "-")
			params[loki] = pm.Slug()
		case "url":
			url := cast.ToString(v)
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("URLs with protocol (http*) not supported: %q. In page %q", url, ps.pathOrTitle())
			}
			pcfg.URL = url
			params[loki] = url
		case "type":
			pcfg.Type = cast.ToString(v)
			params[loki] = pcfg.Type
		case "keywords":
			pcfg.Keywords = cast.ToStringSlice(v)
			params[loki] = pcfg.Keywords
		case "headless":
			// Legacy setting for leaf bundles.
			// This is since Hugo 0.63 handled in a more general way for all
			// pages.
			isHeadless := cast.ToBool(v)
			params[loki] = isHeadless
			if isHeadless {
				pm.pageConfig.Build.List = pagemeta.Never
				pm.pageConfig.Build.Render = pagemeta.Never
			}
		case "outputs":
			o := cast.ToStringSlice(v)
			// lower case names:
			for i, s := range o {
				o[i] = strings.ToLower(s)
			}
			pm.pageConfig.Outputs = o
		case "draft":
			draft = new(bool)
			*draft = cast.ToBool(v)
		case "layout":
			pcfg.Layout = cast.ToString(v)
			params[loki] = pcfg.Layout
		case "markup":
			pcfg.Content.Markup = cast.ToString(v)
			params[loki] = pcfg.Content.Markup
		case "weight":
			pcfg.Weight = cast.ToInt(v)
			params[loki] = pcfg.Weight
		case "aliases":
			pcfg.Aliases = cast.ToStringSlice(v)
			for i, alias := range pcfg.Aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("http* aliases not supported: %q", alias)
				}
				pcfg.Aliases[i] = filepath.ToSlash(alias)
			}
			params[loki] = pcfg.Aliases
		case "sitemap":
			pcfg.Sitemap, err = config.DecodeSitemap(ps.s.conf.Sitemap, maps.ToStringMap(v))
			if err != nil {
				return fmt.Errorf("failed to decode sitemap config in front matter: %s", err)
			}
			sitemapSet = true
		case "iscjklanguage":
			isCJKLanguage = new(bool)
			*isCJKLanguage = cast.ToBool(v)
		case "translationkey":
			pcfg.TranslationKey = cast.ToString(v)
			params[loki] = pcfg.TranslationKey
		case "resources":
			var resources []map[string]any
			handled := true

			switch vv := v.(type) {
			case []map[any]any:
				for _, vvv := range vv {
					resources = append(resources, maps.ToStringMap(vvv))
				}
			case []map[string]any:
				resources = append(resources, vv...)
			case []any:
				for _, vvv := range vv {
					switch vvvv := vvv.(type) {
					case map[any]any:
						resources = append(resources, maps.ToStringMap(vvvv))
					case map[string]any:
						resources = append(resources, vvvv)
					}
				}
			default:
				handled = false
			}

			if handled {
				pcfg.ResourcesMeta = resources
				break
			}
			fallthrough
		default:
			// If not one of the explicit values, store in Params
			switch vv := v.(type) {
			case []any:
				if len(vv) > 0 {
					allStrings := true
					for _, vvv := range vv {
						if _, ok := vvv.(string); !ok {
							allStrings = false
							break
						}
					}
					if allStrings {
						// We need tags, keywords etc. to be []string, not []interface{}.
						a := make([]string, len(vv))
						for i, u := range vv {
							a[i] = cast.ToString(u)
						}
						params[loki] = a
					} else {
						params[loki] = vv
					}
				} else {
					params[loki] = []string{}
				}

			default:
				params[loki] = vv
			}
		}
	}

	for k, v := range userParams {
		params[strings.ToLower(k)] = v
	}

	if !sitemapSet {
		pcfg.Sitemap = ps.s.conf.Sitemap
	}

	if draft != nil && published != nil {
		pcfg.Draft = *draft
		ps.s.Log.Warnf("page %q has both draft and published settings in its frontmatter. Using draft.", ps.File().Filename())
	} else if draft != nil {
		pcfg.Draft = *draft
	} else if published != nil {
		pcfg.Draft = !*published
	}
	params["draft"] = pcfg.Draft

	if isCJKLanguage != nil {
		pcfg.IsCJKLanguage = *isCJKLanguage
	} else if ps.s.conf.HasCJKLanguage && ps.m.content.pi.openSource != nil {
		if cjkRe.Match(ps.m.content.mustSource()) {
			pcfg.IsCJKLanguage = true
		} else {
			pcfg.IsCJKLanguage = false
		}
	}

	params["iscjklanguage"] = pcfg.IsCJKLanguage

	if err := pcfg.Init(false); err != nil {
		return err
	}

	if err := pcfg.Compile(ext, ps.s.Log, ps.s.conf.OutputFormats.Config, ps.s.conf.MediaTypes.Config); err != nil {
		return err
	}

	return nil
}

// shouldList returns whether this page should be included in the list of pages.
// global indicates site.Pages etc.
func (m *pageMeta) shouldList(global bool) bool {
	if m.isStandalone() {
		// Never list 404, sitemap and similar.
		return false
	}

	switch m.pageConfig.Build.List {
	case pagemeta.Always:
		return true
	case pagemeta.Never:
		return false
	case pagemeta.ListLocally:
		return !global
	}
	return false
}

func (m *pageMeta) shouldListAny() bool {
	return m.shouldList(true) || m.shouldList(false)
}

func (m *pageMeta) isStandalone() bool {
	return !m.standaloneOutputFormat.IsZero()
}

func (m *pageMeta) shouldBeCheckedForMenuDefinitions() bool {
	if !m.shouldList(false) {
		return false
	}

	return m.pageConfig.Kind == kinds.KindHome || m.pageConfig.Kind == kinds.KindSection || m.pageConfig.Kind == kinds.KindPage
}

func (m *pageMeta) noRender() bool {
	return m.pageConfig.Build.Render != pagemeta.Always
}

func (m *pageMeta) noLink() bool {
	return m.pageConfig.Build.Render == pagemeta.Never
}

func (m *pageMeta) applyDefaultValues(s *Site) error {
	if m.pageConfig.Build.IsZero() {
		m.pageConfig.Build, _ = pagemeta.DecodeBuildConfig(nil)
	}

	if !s.conf.IsKindEnabled(m.Kind()) {
		(&m.pageConfig.Build).Disable()
	}

	if m.pageConfig.Content.Markup == "" {
		if m.File() != nil {
			// Fall back to file extension
			m.pageConfig.Content.Markup = s.ContentSpec.ResolveMarkup(m.File().Ext())
		}
		if m.pageConfig.Content.Markup == "" {
			m.pageConfig.Content.Markup = "markdown"
		}
	}

	if m.pageConfig.Title == "" && m.f == nil {
		switch m.Kind() {
		case kinds.KindHome:
			m.pageConfig.Title = s.Title()
		case kinds.KindSection:
			sectionName := m.pathInfo.Unnormalized().BaseNameNoIdentifier()
			if s.conf.PluralizeListTitles {
				sectionName = flect.Pluralize(sectionName)
			}
			if s.conf.CapitalizeListTitles {
				sectionName = s.conf.C.CreateTitle(sectionName)
			}
			m.pageConfig.Title = sectionName
		case kinds.KindTerm:
			if m.term != "" {
				if s.conf.CapitalizeListTitles {
					m.pageConfig.Title = s.conf.C.CreateTitle(m.term)
				} else {
					m.pageConfig.Title = m.term
				}
			} else {
				panic("term not set")
			}
		case kinds.KindTaxonomy:
			if s.conf.CapitalizeListTitles {
				m.pageConfig.Title = strings.Replace(s.conf.C.CreateTitle(m.pathInfo.Unnormalized().BaseNameNoIdentifier()), "-", " ", -1)
			} else {
				m.pageConfig.Title = strings.Replace(m.pathInfo.Unnormalized().BaseNameNoIdentifier(), "-", " ", -1)
			}
		case kinds.KindStatus404:
			m.pageConfig.Title = "404 Page not found"
		}
	}

	return nil
}

func (m *pageMeta) newContentConverter(ps *pageState, markup string) (converter.Converter, error) {
	if ps == nil {
		panic("no Page provided")
	}
	cp := ps.s.ContentSpec.Converters.Get(markup)
	if cp == nil {
		return converter.NopConverter, fmt.Errorf("no content renderer found for markup %q, page: %s", markup, ps.getPageInfoForError())
	}

	var id string
	var filename string
	var path string
	if m.f != nil {
		id = m.f.UniqueID()
		filename = m.f.Filename()
		path = m.f.Path()
	} else {
		path = m.Path()
	}

	doc := newPageForRenderHook(ps)

	documentLookup := func(id uint64) any {
		if id == ps.pid {
			// This prevents infinite recursion in some cases.
			return doc
		}
		if v, ok := ps.pageOutput.pco.otherOutputs.Get(id); ok {
			return v.po.p
		}
		return nil
	}

	cpp, err := cp.New(
		converter.DocumentContext{
			Document:       doc,
			DocumentLookup: documentLookup,
			DocumentID:     id,
			DocumentName:   path,
			Filename:       filename,
		},
	)
	if err != nil {
		return converter.NopConverter, err
	}

	return cpp, nil
}

// The output formats this page will be rendered to.
func (ps *pageState) outputFormats() output.Formats {
	if len(ps.m.pageConfig.ConfiguredOutputFormats) > 0 {
		return ps.m.pageConfig.ConfiguredOutputFormats
	}
	return ps.s.conf.C.KindOutputFormats[ps.Kind()]
}

func (m *pageMeta) Slug() string {
	return m.pageConfig.Slug
}

func getParam(m resource.ResourceParamsProvider, key string, stringToLower bool) any {
	v := m.Params()[strings.ToLower(key)]

	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		if stringToLower {
			return strings.ToLower(val)
		}
		return val
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return val
	case []string:
		if stringToLower {
			return helpers.SliceToLower(val)
		}
		return v
	default:
		return v
	}
}

// Implement contentNodeI.
// Note that pageMeta is just a temporary contentNode. It will be replaced in the tree with a *pageState.
// TODO1 make some partial interfaces.
func (m *pageMeta) GetIdentity() identity.Identity {
	panic("not supported")
}

func (m *pageMeta) ForEeachIdentity(cb func(id identity.Identity) bool) bool {
	panic("not supported")
}

func (m *pageMeta) isContentNodeBranch() bool {
	panic("not supported")
}

func (m *pageMeta) contentWeight() int {
	if m.f == nil {
		return 0
	}
	return m.f.FileInfo().Meta().Weight
}

func (m *pageMeta) matchDirectOrInDelegees(_ sitesmatrix.Vector) (contentNode, sitesmatrix.Vector) {
	panic("not implemented") // TODO: Implement
}

func (m *pageMeta) sitesMatrix() sitesmatrix.VectorProvider {
	return m.sitesMatrix()
}

func (m *pageMeta) resetBuildState() {
	panic("not supported")
}

func (m *pageMeta) MarkStale() {
	// panic("not supported")
}

func getParamToLower(m resource.ResourceParamsProvider, key string) any {
	return getParam(m, key, true)
}

func (ps *pageState) initLazyProviders() error {
	ps.init.Add(func(ctx context.Context) (any, error) {
		pp, err := newPagePaths(ps)
		if err != nil {
			return nil, err
		}

		var outputFormatsForPage output.Formats
		var renderFormats output.Formats

		if ps.m.standaloneOutputFormat.IsZero() {
			outputFormatsForPage = ps.outputFormats()
			renderFormats = ps.s.h.renderFormats
		} else {
			// One of the fixed output format pages, e.g. 404.
			outputFormatsForPage = output.Formats{ps.m.standaloneOutputFormat}
			renderFormats = outputFormatsForPage
		}

		// Prepare output formats for all sites.
		// We do this even if this page does not get rendered on
		// its own. It may be referenced via one of the site collections etc.
		// it will then need an output format.
		ps.pageOutputs = make([]*pageOutput, len(renderFormats))
		created := make(map[string]*pageOutput)
		shouldRenderPage := !ps.m.noRender()

		for i, f := range renderFormats {

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
				contentProvider, err := newPageContentOutput(po)
				if err != nil {
					return nil, err
				}
				po.setContentProvider(contentProvider)
			}

			ps.pageOutputs[i] = po
			created[f.Name] = po

		}

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return nil
}
