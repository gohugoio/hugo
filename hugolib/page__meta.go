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
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	xmaps "maps"

	"github.com/bep/logg"
	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hiter"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
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

func (m *pageMetaSource) GetIdentity() identity.Identity {
	return m.pathInfo
}

func (m *pageMetaSource) Path() string {
	return m.pathInfo.Base()
}

func (m *pageMetaSource) PathInfo() *paths.Path {
	return m.pathInfo
}

func (m *pageMetaSource) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	return f(sitesmatrix.Vector{}, m)
}

type pageMetaSource struct {
	pathInfo                      *paths.Path // Always set. This the canonical path to the Page.
	f                             *source.File
	pi                            *contentParseInfo
	contentAdapterSourceEntryHash uint64

	sitesMatrixBase     sitesmatrix.VectorIterator
	sitesMatrixBaseOnly bool

	// Set for standalone pages, e.g. robotsTXT.
	standaloneOutputFormat output.Format
	resources.AtomicStaler

	pageConfigSource *pagemeta.PageConfigEarly
	cascadeCompiled  *page.PageMatcherParamsConfigs

	resourcePath  string // Set for bundled pages; path relative to its bundle root.
	bundled       bool   // Set if this page is bundled inside another.
	noFrontMatter bool
	openSource    func() (hugio.ReadSeekCloser, error)
}

func (m *pageMetaSource) String() string {
	if m.f != nil {
		return fmt.Sprintf("pageMetaSource(%s, %s)", m.f.FileInfo().Meta().PathInfo, m.f.FileInfo().Meta().Filename)
	}
	return fmt.Sprintf("pageMetaSource(%s)", m.pathInfo)
}

func (m *pageMetaSource) nodeCategoryPage() {
	// Marker method.
}

func (m *pageMetaSource) nodeCategorySingle() {
	// Marker method.
}

func (m *pageMetaSource) nodeSourceEntryID() any {
	if m.f != nil && !m.f.IsContentAdapter() {
		return m.f.FileInfo().Meta().Filename
	}
	return m.contentAdapterSourceEntryHash
}

func (m *pageMetaSource) setCascadeFromMap(frontmatter map[string]any, defaultSitesMatrix sitesmatrix.VectorStore, configuredDimensions *sitesmatrix.ConfiguredDimensions, logger loggers.Logger) error {
	const (
		pageMetaKeyCascade = "cascade"
	)

	// Check for any cascade define on itself.
	if cv, found := frontmatter[pageMetaKeyCascade]; found {
		var err error
		cascade, err := page.DecodeCascadeConfig(cv)
		if err != nil {
			return err
		}
		if err := cascade.InitConfig(logger, defaultSitesMatrix, configuredDimensions); err != nil {
			return err
		}

		m.cascadeCompiled = cascade
	}
	return nil
}

var _ contentNodeCascadeProvider = (*pageState)(nil)

type pageMeta struct {
	// Shared between all dimensions of this page.
	*pageMetaSource

	pageConfig *pagemeta.PageConfigLate

	term     string // Set for kind == KindTerm.
	singular string // Set for kind == KindTerm and kind == KindTaxonomy.

	content *cachedContent // The source and the parsed page content.

	datesOriginal pagemeta.Dates // Original dates for rebuilds.
}

func (p *pageState) getCascade() *page.PageMatcherParamsConfigs {
	return p.m.cascadeCompiled
}

func (m *pageMetaSource) initEarly(h *HugoSites, cascades *page.PageMatcherParamsConfigs) error {
	if err := m.doInitEarly(h, cascades); err != nil {
		return m.wrapError(err, h.SourceFs)
	}
	return nil
}

func (m *pageMetaSource) doInitEarly(h *HugoSites, cascades *page.PageMatcherParamsConfigs) error {
	if err := m.initFrontMatter(h); err != nil {
		return err
	}

	if m.pageConfigSource.Kind == "" {
		// Resolve page kind.
		m.pageConfigSource.Kind = kinds.KindSection
		if m.pathInfo.Base() == "" {
			m.pageConfigSource.Kind = kinds.KindHome
		} else if m.pathInfo.IsBranchBundle() {
			// A section, taxonomy or term.
			tc := h.getFirstTaxonomyConfig(m.Path())
			if !tc.IsZero() {
				// Either a taxonomy or a term.
				if tc.pluralTreeKey == m.Path() {
					m.pageConfigSource.Kind = kinds.KindTaxonomy
				} else {
					m.pageConfigSource.Kind = kinds.KindTerm
				}
			}
		} else if m.f != nil {
			m.pageConfigSource.Kind = kinds.KindPage
		}
	}

	sitesMatrixBase := m.sitesMatrixBase
	if sitesMatrixBase == nil && m.f != nil {
		sitesMatrixBase = m.f.FileInfo().Meta().SitesMatrix
	}

	if m.sitesMatrixBaseOnly && sitesMatrixBase == nil {
		panic("sitesMatrixBaseOnly set, but no base sites matrix")
	}

	var fim *hugofs.FileMeta
	if m.f != nil {
		fim = m.f.FileInfo().Meta()
	}
	if err := m.pageConfigSource.CompileEarly(m.pathInfo, cascades, h.Conf, fim, sitesMatrixBase, m.sitesMatrixBaseOnly); err != nil {
		return err
	}

	if cnh.isBranchNode(m) && m.pageConfigSource.Frontmatter != nil {
		if err := m.setCascadeFromMap(m.pageConfigSource.Frontmatter, m.pageConfigSource.SitesMatrix, h.Conf.ConfiguredDimensions(), h.Log); err != nil {
			return err
		}
	}
	return nil
}

func (m *pageMetaSource) initFrontMatter(h *HugoSites) error {
	if m.pathInfo == nil {
		if err := m.initPathInfo(h); err != nil {
			return err
		}
	}

	if m.pageConfigSource == nil {
		m.pageConfigSource = &pagemeta.PageConfigEarly{}
	}

	// Read front matter.
	if err := m.parseFrontMatter(h, pageSourceIDCounter.Add(1)); err != nil {
		return err
	}
	var ext string
	if m.f != nil {
		if !m.f.IsContentAdapter() {
			ext = m.f.Ext()
		}
	}

	if err := m.pageConfigSource.SetMetaPreFromMap(ext, m.pi.frontMatter, h.Log, h.Conf); err != nil {
		return err
	}

	return nil
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

func (m *pageMeta) initLate(s *Site) error {
	var tc viewName

	if m.pageConfigSource.Kind == kinds.KindTerm || m.pageConfigSource.Kind == kinds.KindTaxonomy {
		if tc.IsZero() {
			tc = s.pageMap.cfg.getTaxonomyConfig(m.Path())
		}
		if tc.IsZero() {
			return fmt.Errorf("no taxonomy configuration found for %q", m.Path())
		}
		m.singular = tc.singular

		if m.pageConfigSource.Kind == kinds.KindTerm {
			m.term = paths.TrimLeading(strings.TrimPrefix(m.pathInfo.Unnormalized().Base(), tc.pluralTreeKey))
		}
	}

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
			f:          source.NewFileInfo(fi),
			pathInfo:   fi.Meta().PathInfo,
			openSource: openSource,
		}

		return p, nil
	}()
	if err != nil {
		return nil, hugofs.AddFileInfoToError(err, fi, h.SourceFs)
	}

	return p, err
}

func (h *HugoSites) newPageMetaSourceForContentAdapter(fi hugofs.FileMetaInfo, sitesMatrixBase sitesmatrix.VectorIterator, pc *pagemeta.PageConfigEarly) (*pageMetaSource, error) {
	p := &pageMetaSource{
		f:                             source.NewFileInfo(fi),
		noFrontMatter:                 true,
		pageConfigSource:              pc,
		contentAdapterSourceEntryHash: pc.SourceEntryHash,
		sitesMatrixBase:               sitesMatrixBase,
		openSource:                    pc.Content.ValueAsOpenReadSeekCloser(),
	}

	return p, p.initPathInfo(h)
}

func (s *Site) newPageFromPageMetasource(ms *pageMetaSource, cascades *page.PageMatcherParamsConfigs) (*pageState, error) {
	m, err := s.newPageMetaFromPageMetasource(ms)
	if err != nil {
		return nil, err
	}
	p, err := s.newPageFromPageMeta(m, cascades)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Site) newPageMetaFromPageMetasource(ms *pageMetaSource) (*pageMeta, error) {
	m := &pageMeta{
		pageMetaSource: ms,
		pageConfig: &pagemeta.PageConfigLate{
			Params: make(maps.Params),
		},
	}
	return m, nil
}

// Prepare for a rebuild of the data passed in from front matter.
func (m *pageMeta) prepareRebuild() {
	// Restore orignal date values so we can repeat aggregation.
	m.pageConfig.Dates = m.datesOriginal
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

func (m *pageMetaSource) Kind() string {
	return m.pageConfigSource.Kind
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
	if m.Kind() == kinds.KindTerm {
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
	if m.IsHome() {
		// Base returns an empty string for the home page,
		// which is correct as the key to the page tree.
		return "/"
	}
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

func (ps *pageState) setMetaPost(cascades *page.PageMatcherParamsConfigs) error {
	if ps.m.pageConfigSource.Frontmatter != nil {
		if ps.m.pageConfigSource.IsFromContentAdapter {
			ps.m.pageConfig.ContentAdapterData = xmaps.Clone(ps.m.pageConfigSource.Frontmatter)
		} else {
			ps.m.pageConfig.Params = xmaps.Clone(ps.m.pageConfigSource.Frontmatter)
		}
	}

	if ps.m.pageConfig.Params == nil {
		ps.m.pageConfig.Params = make(maps.Params)
	}
	if ps.m.pageConfigSource.IsFromContentAdapter && ps.m.pageConfig.ContentAdapterData == nil {
		ps.m.pageConfig.ContentAdapterData = make(maps.Params)
	}

	// Cascade defined on itself has higher priority than inherited ones.
	allCascades := hiter.Concat(ps.m.cascadeCompiled.All(), cascades.All())

	for v := range allCascades {
		if !v.Target.Match(ps.Kind(), ps.Path(), ps.s.Conf.Environment(), ps.s.siteVector) {
			continue
		}

		for kk, vv := range v.Params {
			if _, found := ps.m.pageConfig.Params[kk]; !found {
				ps.m.pageConfig.Params[kk] = vv
			}
		}
		for kk, vv := range v.Fields {
			if ps.m.pageConfigSource.IsFromContentAdapter {
				if _, found := ps.m.pageConfig.ContentAdapterData[kk]; !found {
					ps.m.pageConfig.ContentAdapterData[kk] = vv
				}
			} else {
				if _, found := ps.m.pageConfig.Params[kk]; !found {
					ps.m.pageConfig.Params[kk] = vv
				}
			}
		}
	}

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
	var isContentAdapter bool
	if ps.File() != nil {
		isContentAdapter = ps.File().IsContentAdapter()
		contentBaseName = ps.File().ContentBaseName()
		if ps.File().FileInfo() != nil {
			mtime = ps.File().FileInfo().ModTime()
		}
	}

	var gitAuthorDate time.Time
	if ps.gitInfo != nil {
		gitAuthorDate = ps.gitInfo.AuthorDate
	}

	descriptor := &pagemeta.FrontMatterDescriptor{
		PageConfigEarly: pm.pageConfigSource,
		PageConfigLate:  pm.pageConfig,
		BaseFilename:    contentBaseName,
		ModTime:         mtime,
		GitAuthorDate:   gitAuthorDate,
		Location:        langs.GetLocation(ps.s.Language()),
		PathOrTitle:     ps.pathOrTitle(),
	}

	if isContentAdapter {
		if err := pm.pageConfig.Compile(pm.pageConfigSource, ps.s.Log, ps.s.conf.OutputFormats.Config); err != nil {
			return err
		}
	}

	pcfg := pm.pageConfig

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

	var sitemapSet bool

	var buildConfig any
	var isNewBuildKeyword bool
	if v, ok := pcfg.Params["_build"]; ok {
		hugo.Deprecate("The \"_build\" front matter key", "Use \"build\" instead. See https://gohugo.io/content-management/build-options.", "0.145.0")
		buildConfig = v
	} else {
		buildConfig = pcfg.Params["build"]
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
			pcfg.Params[loki] = pcfg.Title
		case "linktitle":
			pcfg.LinkTitle = cast.ToString(v)
			pcfg.Params[loki] = pcfg.LinkTitle
		case "summary":
			pcfg.Summary = cast.ToString(v)
			pcfg.Params[loki] = pcfg.Summary
		case "description":
			pcfg.Description = cast.ToString(v)
			pcfg.Params[loki] = pcfg.Description
		case "slug":
			// Don't start or end with a -
			pcfg.Slug = strings.Trim(cast.ToString(v), "-")
			pcfg.Params[loki] = pm.Slug()
		case "url":
			url := cast.ToString(v)
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("URLs with protocol (http*) not supported: %q. In page %q", url, ps.pathOrTitle())
			}
			pcfg.URL = url
			pcfg.Params[loki] = url
		case "type":
			pcfg.Type = cast.ToString(v)
			pcfg.Params[loki] = pcfg.Type
		case "keywords":
			pcfg.Keywords = cast.ToStringSlice(v)
			pcfg.Params[loki] = pcfg.Keywords
		case "headless":
			// Legacy setting for leaf bundles.
			// This is since Hugo 0.63 handled in a more general way for all
			// pages.
			isHeadless := cast.ToBool(v)
			pcfg.Params[loki] = isHeadless
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
			pcfg.Params[loki] = pcfg.Layout
		case "weight":
			pcfg.Weight = cast.ToInt(v)
			pcfg.Params[loki] = pcfg.Weight
		case "aliases":
			pcfg.Aliases = cast.ToStringSlice(v)
			for i, alias := range pcfg.Aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("http* aliases not supported: %q", alias)
				}
				pcfg.Aliases[i] = filepath.ToSlash(alias)
			}
			pcfg.Params[loki] = pcfg.Aliases
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
			pcfg.Params[loki] = pcfg.TranslationKey
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
						pcfg.Params[loki] = a
					} else {
						pcfg.Params[loki] = vv
					}
				} else {
					pcfg.Params[loki] = []string{}
				}

			default:
				pcfg.Params[loki] = vv
			}
		}
	}

	for k, v := range userParams {
		pcfg.Params[strings.ToLower(k)] = v
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
	pcfg.Params["draft"] = pcfg.Draft

	if isCJKLanguage != nil {
		pcfg.IsCJKLanguage = *isCJKLanguage
	} else if ps.s.conf.HasCJKLanguage && ps.m.content.pi.openSource != nil {
		if cjkRe.Match(ps.m.content.mustSource()) {
			pcfg.IsCJKLanguage = true
		} else {
			pcfg.IsCJKLanguage = false
		}
	}

	pcfg.Params["iscjklanguage"] = pcfg.IsCJKLanguage

	if err := pm.pageConfigSource.Init(false); err != nil {
		return err
	}

	if err := pcfg.Init(); err != nil {
		return err
	}

	if err := pcfg.Compile(pm.pageConfigSource, ps.s.Log, ps.s.conf.OutputFormats.Config); err != nil {
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

	return m.Kind() == kinds.KindHome || m.Kind() == kinds.KindSection || m.Kind() == kinds.KindPage
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

	if m.pageConfigSource.Content.Markup == "" {
		if m.File() != nil {
			// Fall back to file extension
			m.pageConfigSource.Content.Markup = s.ContentSpec.ResolveMarkup(m.File().Ext())
		}
		if m.pageConfigSource.Content.Markup == "" {
			m.pageConfigSource.Content.Markup = "markdown"
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

	var filename string
	var path string
	if m.f != nil {
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
			DocumentID:     hashing.XxHashFromStringHexEncoded(ps.Path()),
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

func getParamToLower(m resource.ResourceParamsProvider, key string) any {
	return getParam(m, key, true)
}
