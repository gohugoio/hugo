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
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/langs"

	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/related"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/maps"
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
	// kind is the discriminator that identifies the different page types
	// in the different page collections. This can, as an example, be used
	// to to filter regular pages, find sections etc.
	// Kind will, for the pages available to the templates, be one of:
	// page, home, section, taxonomy and term.
	// It is of string type to make it easy to reason about in
	// the templates.
	kind string

	// This is a standalone page not part of any page collection. These
	// include sitemap, robotsTXT and similar. It will have no pageOutputs, but
	// a fixed pageOutput.
	standalone bool

	draft       bool // Only published when running with -D flag
	buildConfig pagemeta.BuildConfig

	bundleType files.ContentClass

	// Params contains configuration defined in the params section of page frontmatter.
	params map[string]any

	title     string
	linkTitle string

	summary string

	resourcePath string

	weight int

	markup      string
	contentType string

	// whether the content is in a CJK language.
	isCJKLanguage bool

	layout string

	aliases []string

	description string
	keywords    []string

	urlPaths pagemeta.URLPath

	resource.Dates

	// Set if this page is bundled inside another.
	bundled bool

	// A key that maps to translation(s) of this page. This value is fetched
	// from the page front matter.
	translationKey string

	// From front matter.
	configuredOutputFormats output.Formats

	// This is the raw front matter metadata that is going to be assigned to
	// the Resources above.
	resourcesMetadata []map[string]any

	f source.File

	sections []string

	// Sitemap overrides from front matter.
	sitemap config.SitemapConfig

	s *Site

	contentConverterInit sync.Once
	contentConverter     converter.Converter
}

func (p *pageMeta) Aliases() []string {
	return p.aliases
}

func (p *pageMeta) Author() page.Author {
	helpers.Deprecated(".Author", "Use taxonomies.", false)
	authors := p.Authors()

	for _, author := range authors {
		return author
	}
	return page.Author{}
}

func (p *pageMeta) Authors() page.AuthorList {
	helpers.Deprecated(".Authors", "Use taxonomies.", true)
	return nil
}

func (p *pageMeta) BundleType() files.ContentClass {
	return p.bundleType
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

func (p *pageMeta) File() source.File {
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
	if !p.File().IsZero() {
		const example = `
  {{ $path := "" }}
  {{ with .File }}
	{{ $path = .Path }}
  {{ else }}
	{{ $path = .Path }}
  {{ end }}
`
		helpers.Deprecated(".Path when the page is backed by a file", "We plan to use Path for a canonical source path and you probably want to check the source is a file. To get the current behaviour, you can use a construct similar to the one below:\n"+example, false)

	}

	return p.Pathc()
}

// This is just a bridge method, use Path in templates.
func (p *pageMeta) Pathc() string {
	if !p.File().IsZero() {
		return p.File().Path()
	}
	return p.SectionsPath()
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
	if p.IsHome() {
		return ""
	}

	if p.IsNode() {
		if len(p.sections) == 0 {
			// May be a sitemap or similar.
			return ""
		}
		return p.sections[0]
	}

	if !p.File().IsZero() {
		return p.File().Section()
	}

	panic("invalid page state")
}

func (p *pageMeta) SectionsEntries() []string {
	return p.sections
}

func (p *pageMeta) SectionsPath() string {
	return path.Join(p.SectionsEntries()...)
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

func (pm *pageMeta) mergeBucketCascades(b1, b2 *pagesMapBucket) {
	if b1.cascade == nil {
		b1.cascade = make(map[page.PageMatcher]maps.Params)
	}

	if b2 != nil && b2.cascade != nil {
		for k, v := range b2.cascade {

			vv, found := b1.cascade[k]
			if !found {
				b1.cascade[k] = v
			} else {
				// Merge
				for ck, cv := range v {
					if _, found := vv[ck]; !found {
						vv[ck] = cv
					}
				}
			}
		}
	}
}

func (pm *pageMeta) setMetadata(parentBucket *pagesMapBucket, p *pageState, frontmatter map[string]any) error {
	pm.params = make(maps.Params)

	if frontmatter == nil && (parentBucket == nil || parentBucket.cascade == nil) {
		return nil
	}

	if frontmatter != nil {
		// Needed for case insensitive fetching of params values
		maps.PrepareParams(frontmatter)
		if p.bucket != nil {
			// Check for any cascade define on itself.
			if cv, found := frontmatter["cascade"]; found {
				var err error
				p.bucket.cascade, err = page.DecodeCascade(cv)
				if err != nil {
					return err
				}
			}
		}
	} else {
		frontmatter = make(map[string]any)
	}

	var cascade map[page.PageMatcher]maps.Params

	if p.bucket != nil {
		if parentBucket != nil {
			// Merge missing keys from parent into this.
			pm.mergeBucketCascades(p.bucket, parentBucket)
		}
		cascade = p.bucket.cascade
	} else if parentBucket != nil {
		cascade = parentBucket.cascade
	}

	for m, v := range cascade {
		if !m.Matches(p) {
			continue
		}
		for kk, vv := range v {
			if _, found := frontmatter[kk]; !found {
				frontmatter[kk] = vv
			}
		}
	}

	var mtime time.Time
	var contentBaseName string
	if !p.File().IsZero() {
		contentBaseName = p.File().ContentBaseName()
		if p.File().FileInfo() != nil {
			mtime = p.File().FileInfo().ModTime()
		}
	}

	var gitAuthorDate time.Time
	if !p.gitInfo.IsZero() {
		gitAuthorDate = p.gitInfo.AuthorDate
	}

	descriptor := &pagemeta.FrontMatterDescriptor{
		Frontmatter:   frontmatter,
		Params:        pm.params,
		Dates:         &pm.Dates,
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

	pm.buildConfig, err = pagemeta.DecodeBuildConfig(frontmatter["_build"])
	if err != nil {
		return err
	}

	var sitemapSet bool

	var draft, published, isCJKLanguage *bool
	for k, v := range frontmatter {
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
			lang := p.s.GetLanguagePrefix()
			if lang != "" && !strings.HasPrefix(url, "/") && strings.HasPrefix(url, lang+"/") {
				if strings.HasPrefix(hugo.CurrentVersion.String(), "0.55") {
					// We added support for page relative URLs in Hugo 0.55 and
					// this may get its language path added twice.
					// TODO(bep) eventually remove this.
					p.s.Log.Warnf(`Front matter in %q with the url %q with no leading / has what looks like the language prefix added. In Hugo 0.55 we added support for page relative URLs in front matter, no language prefix needed. Check the URL and consider to either add a leading / or remove the language prefix.`, p.pathOrTitle(), url)
				}
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
	} else if p.s.conf.HasCJKLanguage && p.source.parsed != nil {
		if cjkRe.Match(p.source.parsed.Input()) {
			pm.isCJKLanguage = true
		} else {
			pm.isCJKLanguage = false
		}
	}

	pm.params["iscjklanguage"] = p.m.isCJKLanguage

	return nil
}

func (p *pageMeta) noListAlways() bool {
	return p.buildConfig.List != pagemeta.Always
}

func (p *pageMeta) getListFilter(local bool) contentTreeNodeCallback {
	return newContentTreeFilter(func(n *contentNode) bool {
		if n == nil {
			return true
		}

		var shouldList bool
		switch n.p.m.buildConfig.List {
		case pagemeta.Always:
			shouldList = true
		case pagemeta.Never:
			shouldList = false
		case pagemeta.ListLocally:
			shouldList = local
		}

		return !shouldList
	})
}

func (p *pageMeta) noRender() bool {
	return p.buildConfig.Render != pagemeta.Always
}

func (p *pageMeta) noLink() bool {
	return p.buildConfig.Render == pagemeta.Never
}

func (p *pageMeta) applyDefaultValues(n *contentNode) error {
	if p.buildConfig.IsZero() {
		p.buildConfig, _ = pagemeta.DecodeBuildConfig(nil)
	}

	if !p.s.isEnabled(p.Kind()) {
		(&p.buildConfig).Disable()
	}

	if p.markup == "" {
		if !p.File().IsZero() {
			// Fall back to file extension
			p.markup = p.s.ContentSpec.ResolveMarkup(p.File().Ext())
		}
		if p.markup == "" {
			p.markup = "markdown"
		}
	}

	if p.title == "" && p.f.IsZero() {
		switch p.Kind() {
		case kinds.KindHome:
			p.title = p.s.Title()
		case kinds.KindSection:
			var sectionName string
			if n != nil {
				sectionName = n.rootSection()
			} else {
				sectionName = p.sections[0]
			}

			sectionName = helpers.FirstUpper(sectionName)
			if p.s.conf.PluralizeListTitles {
				p.title = flect.Pluralize(sectionName)
			} else {
				p.title = sectionName
			}
		case kinds.KindTerm:
			// TODO(bep) improve
			key := p.sections[len(p.sections)-1]
			p.title = strings.Replace(p.s.conf.C.CreateTitle(key), "-", " ", -1)
		case kinds.KindTaxonomy:
			p.title = p.s.conf.C.CreateTitle(p.sections[0])
		case kinds.Kind404:
			p.title = "404 Page not found"

		}
	}

	if p.IsNode() {
		p.bundleType = files.ContentClassBranch
	} else {
		source := p.File()
		if fi, ok := source.(*fileInfo); ok {
			class := fi.FileInfo().Meta().Classifier
			switch class {
			case files.ContentClassBranch, files.ContentClassLeaf:
				p.bundleType = class
			}
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
	if !p.f.IsZero() {
		id = p.f.UniqueID()
		filename = p.f.Filename()
		path = p.f.Path()
	} else {
		path = p.Pathc()
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
