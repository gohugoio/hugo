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
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/converter"
	xmaps "golang.org/x/exp/maps"

	"github.com/gohugoio/hugo/related"

	"github.com/gohugoio/hugo/source"

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

type pageMeta struct {
	kind     string // Page kind.
	term     string // Set for kind == KindTerm.
	singular string // Set for kind == KindTerm and kind == KindTaxonomy.

	resource.Staler
	pageMetaParams

	pageMetaFrontMatter

	// Set for standalone pages, e.g. robotsTXT.
	standaloneOutputFormat output.Format

	resourcePath string // Set for bundled pages; path relative to its bundle root.
	bundled      bool   // Set if this page is bundled inside another.

	pathInfo *paths.Path // Always set. This the canonical path to the Page.
	f        *source.File

	s *Site // The site this page belongs to.
}

// Prepare for a rebuild of the data passed in from front matter.
func (m *pageMeta) setMetaPostPrepareRebuild() {
	params := xmaps.Clone[map[string]any](m.paramsOriginal)
	m.pageMetaParams.params = params
	m.pageMetaFrontMatter = pageMetaFrontMatter{}
}

type pageMetaParams struct {
	setMetaPostCount          int
	setMetaPostCascadeChanged bool

	params  map[string]any                   // Params contains configuration defined in the params section of page frontmatter.
	cascade map[page.PageMatcher]maps.Params // cascade contains default configuration to be cascaded downwards.

	// These are only set in watch mode.
	datesOriginal   pageMetaDates
	paramsOriginal  map[string]any                   // contains the original params as defined in the front matter.
	cascadeOriginal map[page.PageMatcher]maps.Params // contains the original cascade as defined in the front matter.
}

// From page front matter.
type pageMetaFrontMatter struct {
	draft          bool // Only published when running with -D flag
	title          string
	linkTitle      string
	summary        string
	weight         int
	markup         string
	contentType    string // type in front matter.
	isCJKLanguage  bool   // whether the content is in a CJK language.
	layout         string
	aliases        []string
	description    string
	keywords       []string
	translationKey string // maps to translation(s) of this page.

	buildConfig             pagemeta.BuildConfig
	configuredOutputFormats output.Formats       // outputs defiend in front matter.
	pageMetaDates                                // The 4 front matter dates that Hugo cares about.
	resourcesMetadata       []map[string]any     // Raw front matter metadata that is going to be assigned to the page resources.
	sitemap                 config.SitemapConfig // Sitemap overrides from front matter.
	urlPaths                pagemeta.URLPath
}

func (m *pageMetaParams) init(preserveOringal bool) {
	if preserveOringal {
		m.paramsOriginal = xmaps.Clone[maps.Params](m.params)
		m.cascadeOriginal = xmaps.Clone[map[page.PageMatcher]maps.Params](m.cascade)
	}
}

func (p *pageMeta) Aliases() []string {
	return p.aliases
}

func (p *pageMeta) Author() page.Author {
	hugo.Deprecate(".Author", "Use taxonomies.", "v0.98.0")
	authors := p.Authors()

	for _, author := range authors {
		return author
	}
	return page.Author{}
}

func (p *pageMeta) Authors() page.AuthorList {
	hugo.Deprecate(".Author", "Use taxonomies.", "v0.112.0")
	return nil
}

func (p *pageMeta) BundleType() string {
	switch p.pathInfo.BundleType() {
	case paths.PathTypeLeaf:
		return "leaf"
	case paths.PathTypeBranch:
		return "branch"
	default:
		return ""
	}
}

func (p *pageMeta) Description() string {
	return p.description
}

func (p *pageMeta) Lang() string {
	return p.s.Lang()
}

func (p *pageMeta) Draft() bool {
	return p.draft
}

func (p *pageMeta) File() *source.File {
	return p.f
}

func (p *pageMeta) IsHome() bool {
	return p.Kind() == kinds.KindHome
}

func (p *pageMeta) Keywords() []string {
	return p.keywords
}

func (p *pageMeta) Kind() string {
	return p.kind
}

func (p *pageMeta) Layout() string {
	return p.layout
}

func (p *pageMeta) LinkTitle() string {
	if p.linkTitle != "" {
		return p.linkTitle
	}

	return p.Title()
}

func (p *pageMeta) Name() string {
	if p.resourcePath != "" {
		return p.resourcePath
	}
	if p.kind == kinds.KindTerm {
		return p.pathInfo.Unmormalized().BaseNameNoIdentifier()
	}
	return p.Title()
}

func (p *pageMeta) IsNode() bool {
	return !p.IsPage()
}

func (p *pageMeta) IsPage() bool {
	return p.Kind() == kinds.KindPage
}

// Param is a convenience method to do lookups in Page's and Site's Params map,
// in that order.
//
// This method is also implemented on SiteInfo.
// TODO(bep) interface
func (p *pageMeta) Param(key any) (any, error) {
	return resource.Param(p, p.s.Params(), key)
}

func (p *pageMeta) Params() maps.Params {
	return p.params
}

func (p *pageMeta) Path() string {
	return p.pathInfo.Base()
}

func (p *pageMeta) PathInfo() *paths.Path {
	return p.pathInfo
}

// RelatedKeywords implements the related.Document interface needed for fast page searches.
func (p *pageMeta) RelatedKeywords(cfg related.IndexConfig) ([]related.Keyword, error) {
	v, err := p.Param(cfg.Name)
	if err != nil {
		return nil, err
	}

	return cfg.ToKeywords(v)
}

func (p *pageMeta) IsSection() bool {
	return p.Kind() == kinds.KindSection
}

func (p *pageMeta) Section() string {
	return p.pathInfo.Section()
}

func (p *pageMeta) Sitemap() config.SitemapConfig {
	return p.sitemap
}

func (p *pageMeta) Title() string {
	return p.title
}

const defaultContentType = "page"

func (p *pageMeta) Type() string {
	if p.contentType != "" {
		return p.contentType
	}

	if sect := p.Section(); sect != "" {
		return sect
	}

	return defaultContentType
}

func (p *pageMeta) Weight() int {
	return p.weight
}

func (ps *pageState) setMetaPre() error {
	pm := ps.m
	p := ps
	frontmatter := p.content.parseInfo.frontMatter
	watching := p.s.watching()

	if frontmatter != nil {
		// Needed for case insensitive fetching of params values
		maps.PrepareParams(frontmatter)
		pm.pageMetaParams.params = frontmatter
		if p.IsNode() {
			// Check for any cascade define on itself.
			if cv, found := frontmatter["cascade"]; found {
				var err error
				cascade, err := page.DecodeCascade(cv)
				if err != nil {
					return err
				}
				pm.pageMetaParams.cascade = cascade

			}
		}
	} else if pm.pageMetaParams.params == nil {
		pm.pageMetaParams.params = make(maps.Params)
	}

	pm.pageMetaParams.init(watching)

	return nil
}

func (ps *pageState) setMetaPost(cascade map[page.PageMatcher]maps.Params) error {
	ps.m.setMetaPostCount++
	var cascadeHashPre uint64
	if ps.m.setMetaPostCount > 1 {
		cascadeHashPre = identity.HashUint64(ps.m.cascade)
		ps.m.cascade = xmaps.Clone[map[page.PageMatcher]maps.Params](ps.m.cascadeOriginal)

	}

	// Apply cascades first so they can be overriden later.
	if cascade != nil {
		if ps.m.cascade != nil {
			for k, v := range cascade {
				vv, found := ps.m.cascade[k]
				if !found {
					ps.m.cascade[k] = v
				} else {
					// Merge
					for ck, cv := range v {
						if _, found := vv[ck]; !found {
							vv[ck] = cv
						}
					}
				}
			}
			cascade = ps.m.cascade
		} else {
			ps.m.cascade = cascade
		}
	}

	if cascade == nil {
		cascade = ps.m.cascade
	}

	if ps.m.setMetaPostCount > 1 {
		ps.m.setMetaPostCascadeChanged = cascadeHashPre != identity.HashUint64(ps.m.cascade)
		if !ps.m.setMetaPostCascadeChanged {
			// No changes, restore any value that may be changed by aggregation.
			ps.m.dates = ps.m.datesOriginal.dates
			return nil
		}
		ps.m.setMetaPostPrepareRebuild()

	}

	// Cascade is also applied to itself.
	for m, v := range cascade {
		if !m.Matches(ps) {
			continue
		}
		for kk, vv := range v {
			if _, found := ps.m.params[kk]; !found {
				ps.m.params[kk] = vv
			}
		}
	}

	if err := ps.setMetaPostParams(); err != nil {
		return err
	}

	if err := ps.m.applyDefaultValues(); err != nil {
		return err
	}

	// Store away any original values that may be changed from aggregation.
	ps.m.datesOriginal = ps.m.pageMetaDates

	return nil
}

func (p *pageState) setMetaPostParams() error {
	pm := p.m
	var mtime time.Time
	var contentBaseName string
	if p.File() != nil {
		contentBaseName = p.File().ContentBaseName()
		if p.File().FileInfo() != nil {
			mtime = p.File().FileInfo().ModTime()
		}
	}

	var gitAuthorDate time.Time
	if !p.gitInfo.IsZero() {
		gitAuthorDate = p.gitInfo.AuthorDate
	}

	pm.pageMetaDates = pageMetaDates{}
	pm.urlPaths = pagemeta.URLPath{}

	descriptor := &pagemeta.FrontMatterDescriptor{
		Params:        pm.params,
		Dates:         &pm.pageMetaDates.dates,
		PageURLs:      &pm.urlPaths,
		BaseFilename:  contentBaseName,
		ModTime:       mtime,
		GitAuthorDate: gitAuthorDate,
		Location:      langs.GetLocation(pm.s.Language()),
	}

	// Handle the date separately
	// TODO(bep) we need to "do more" in this area so this can be split up and
	// more easily tested without the Page, but the coupling is strong.
	err := pm.s.frontmatterHandler.HandleDates(descriptor)
	if err != nil {
		p.s.Log.Errorf("Failed to handle dates for page %q: %s", p.pathOrTitle(), err)
	}

	pm.buildConfig, err = pagemeta.DecodeBuildConfig(pm.params["_build"])
	if err != nil {
		return err
	}

	var sitemapSet bool

	var draft, published, isCJKLanguage *bool
	for k, v := range pm.params {
		loki := strings.ToLower(k)

		if loki == "published" { // Intentionally undocumented
			vv, err := cast.ToBoolE(v)
			if err == nil {
				published = &vv
			}
			// published may also be a date
			continue
		}

		if pm.s.frontmatterHandler.IsDateKey(loki) {
			continue
		}

		switch loki {
		case "title":
			pm.title = cast.ToString(v)
			pm.params[loki] = pm.title
		case "linktitle":
			pm.linkTitle = cast.ToString(v)
			pm.params[loki] = pm.linkTitle
		case "summary":
			pm.summary = cast.ToString(v)
			pm.params[loki] = pm.summary
		case "description":
			pm.description = cast.ToString(v)
			pm.params[loki] = pm.description
		case "slug":
			// Don't start or end with a -
			pm.urlPaths.Slug = strings.Trim(cast.ToString(v), "-")
			pm.params[loki] = pm.Slug()
		case "url":
			url := cast.ToString(v)
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("URLs with protocol (http*) not supported: %q. In page %q", url, p.pathOrTitle())
			}
			pm.urlPaths.URL = url
			pm.params[loki] = url
		case "type":
			pm.contentType = cast.ToString(v)
			pm.params[loki] = pm.contentType
		case "keywords":
			pm.keywords = cast.ToStringSlice(v)
			pm.params[loki] = pm.keywords
		case "headless":
			// Legacy setting for leaf bundles.
			// This is since Hugo 0.63 handled in a more general way for all
			// pages.
			isHeadless := cast.ToBool(v)
			pm.params[loki] = isHeadless
			if p.File().TranslationBaseName() == "index" && isHeadless {
				pm.buildConfig.List = pagemeta.Never
				pm.buildConfig.Render = pagemeta.Never
			}
		case "outputs":
			o := cast.ToStringSlice(v)
			// lower case names:
			for i, s := range o {
				o[i] = strings.ToLower(s)
			}
			if len(o) > 0 {
				// Output formats are explicitly set in front matter, use those.
				outFormats, err := p.s.conf.OutputFormats.Config.GetByNames(o...)
				if err != nil {
					p.s.Log.Errorf("Failed to resolve output formats: %s", err)
				} else {
					pm.configuredOutputFormats = outFormats
					pm.params[loki] = outFormats
				}
			}
		case "draft":
			draft = new(bool)
			*draft = cast.ToBool(v)
		case "layout":
			pm.layout = cast.ToString(v)
			pm.params[loki] = pm.layout
		case "markup":
			pm.markup = cast.ToString(v)
			pm.params[loki] = pm.markup
		case "weight":
			pm.weight = cast.ToInt(v)
			pm.params[loki] = pm.weight
		case "aliases":
			pm.aliases = cast.ToStringSlice(v)
			for i, alias := range pm.aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("http* aliases not supported: %q", alias)
				}
				pm.aliases[i] = filepath.ToSlash(alias)
			}
			pm.params[loki] = pm.aliases
		case "sitemap":
			p.m.sitemap, err = config.DecodeSitemap(p.s.conf.Sitemap, maps.ToStringMap(v))
			if err != nil {
				return fmt.Errorf("failed to decode sitemap config in front matter: %s", err)
			}
			pm.params[loki] = p.m.sitemap
			sitemapSet = true
		case "iscjklanguage":
			isCJKLanguage = new(bool)
			*isCJKLanguage = cast.ToBool(v)
		case "translationkey":
			pm.translationKey = cast.ToString(v)
			pm.params[loki] = pm.translationKey
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
				pm.params[loki] = resources
				pm.resourcesMetadata = resources
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
						pm.params[loki] = a
					} else {
						pm.params[loki] = vv
					}
				} else {
					pm.params[loki] = []string{}
				}

			default:
				pm.params[loki] = vv
			}
		}
	}

	if !sitemapSet {
		pm.sitemap = p.s.conf.Sitemap
	}

	pm.markup = p.s.ContentSpec.ResolveMarkup(pm.markup)

	if draft != nil && published != nil {
		pm.draft = *draft
		p.m.s.Log.Warnf("page %q has both draft and published settings in its frontmatter. Using draft.", p.File().Filename())
	} else if draft != nil {
		pm.draft = *draft
	} else if published != nil {
		pm.draft = !*published
	}
	pm.params["draft"] = pm.draft

	if isCJKLanguage != nil {
		pm.isCJKLanguage = *isCJKLanguage
	} else if p.s.conf.HasCJKLanguage && p.content.openSource != nil {
		if cjkRe.Match(p.content.mustSource()) {
			pm.isCJKLanguage = true
		} else {
			pm.isCJKLanguage = false
		}
	}

	pm.params["iscjklanguage"] = p.m.isCJKLanguage

	return nil
}

// shouldList returns whether this page should be included in the list of pages.
// glogal indicates site.Pages etc.
func (p *pageMeta) shouldList(global bool) bool {
	if p.isStandalone() {
		// Never list 404, sitemap and similar.
		return false
	}

	switch p.buildConfig.List {
	case pagemeta.Always:
		return true
	case pagemeta.Never:
		return false
	case pagemeta.ListLocally:
		return !global
	}
	return false
}

func (p *pageMeta) shouldListAny() bool {
	return p.shouldList(true) || p.shouldList(false)
}

func (p *pageMeta) isStandalone() bool {
	return !p.standaloneOutputFormat.IsZero()
}

func (p *pageMeta) shouldBeCheckedForMenuDefinitions() bool {
	if !p.shouldList(false) {
		return false
	}

	return p.kind == kinds.KindHome || p.kind == kinds.KindSection || p.kind == kinds.KindPage
}

func (p *pageMeta) noRender() bool {
	return p.buildConfig.Render != pagemeta.Always
}

func (p *pageMeta) noLink() bool {
	return p.buildConfig.Render == pagemeta.Never
}

func (p *pageMeta) applyDefaultValues() error {
	if p.buildConfig.IsZero() {
		p.buildConfig, _ = pagemeta.DecodeBuildConfig(nil)
	}

	if !p.s.conf.IsKindEnabled(p.Kind()) {
		(&p.buildConfig).Disable()
	}

	if p.markup == "" {
		if p.File() != nil {
			// Fall back to file extension
			p.markup = p.s.ContentSpec.ResolveMarkup(p.File().Ext())
		}
		if p.markup == "" {
			p.markup = "markdown"
		}
	}

	if p.title == "" && p.f == nil {
		switch p.Kind() {
		case kinds.KindHome:
			p.title = p.s.Title()
		case kinds.KindSection:
			sectionName := p.pathInfo.Unmormalized().BaseNameNoIdentifier()
			if p.s.conf.PluralizeListTitles {
				sectionName = flect.Pluralize(sectionName)
			}
			p.title = p.s.conf.C.CreateTitle(sectionName)
		case kinds.KindTerm:
			if p.term != "" {
				p.title = p.s.conf.C.CreateTitle(p.term)
			} else {
				panic("term not set")
			}
		case kinds.KindTaxonomy:
			p.title = strings.Replace(p.s.conf.C.CreateTitle(p.pathInfo.Unmormalized().BaseNameNoIdentifier()), "-", " ", -1)
		case kinds.KindStatus404:
			p.title = "404 Page not found"
		}
	}

	return nil
}

func (p *pageMeta) newContentConverter(ps *pageState, markup string) (converter.Converter, error) {
	if ps == nil {
		panic("no Page provided")
	}
	cp := p.s.ContentSpec.Converters.Get(markup)
	if cp == nil {
		return converter.NopConverter, fmt.Errorf("no content renderer found for markup %q", markup)
	}

	var id string
	var filename string
	var path string
	if p.f != nil {
		id = p.f.UniqueID()
		filename = p.f.Filename()
		path = p.f.Path()
	} else {
		path = p.Path()
	}

	cpp, err := cp.New(
		converter.DocumentContext{
			Document:     newPageForRenderHook(ps),
			DocumentID:   id,
			DocumentName: path,
			Filename:     filename,
		},
	)
	if err != nil {
		return converter.NopConverter, err
	}

	return cpp, nil
}

// The output formats this page will be rendered to.
func (m *pageMeta) outputFormats() output.Formats {
	if len(m.configuredOutputFormats) > 0 {
		return m.configuredOutputFormats
	}
	return m.s.conf.C.KindOutputFormats[m.Kind()]
}

func (p *pageMeta) Slug() string {
	return p.urlPaths.Slug
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

type pageMetaDates struct {
	dates resource.Dates
}

func (d *pageMetaDates) Date() time.Time {
	return d.dates.Date()
}

func (d *pageMetaDates) Lastmod() time.Time {
	return d.dates.Lastmod()
}

func (d *pageMetaDates) PublishDate() time.Time {
	return d.dates.PublishDate()
}

func (d *pageMetaDates) ExpiryDate() time.Time {
	return d.dates.ExpiryDate()
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
			outputFormatsForPage = ps.m.outputFormats()
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
