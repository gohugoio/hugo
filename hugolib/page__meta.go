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

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/related"

	"github.com/gohugoio/hugo/source"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/output"
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
	// page, home, section, taxonomy and taxonomyTerm.
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
	params map[string]interface{}

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
	resourcesMetadata []map[string]interface{}

	f source.File

	sections []string

	// Sitemap overrides from front matter.
	sitemap config.Sitemap

	s *Site

	renderingConfigOverrides map[string]interface{}
	contentConverterInit     sync.Once
	contentConverter         converter.Converter
}

func (p *pageMeta) Aliases() []string {
	return p.aliases
}

func (p *pageMeta) Author() page.Author {
	authors := p.Authors()

	for _, author := range authors {
		return author
	}
	return page.Author{}
}

func (p *pageMeta) Authors() page.AuthorList {
	authorKeys, ok := p.params["authors"]
	if !ok {
		return page.AuthorList{}
	}
	authors := authorKeys.([]string)
	if len(authors) < 1 || len(p.s.Info.Authors) < 1 {
		return page.AuthorList{}
	}

	al := make(page.AuthorList)
	for _, author := range authors {
		a, ok := p.s.Info.Authors[author]
		if ok {
			al[author] = a
		}
	}
	return al
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
	return p.Kind() == page.KindHome
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
	return p.Kind() == page.KindPage
}

// Param is a convenience method to do lookups in Page's and Site's Params map,
// in that order.
//
// This method is also implemented on SiteInfo.
// TODO(bep) interface
func (p *pageMeta) Param(key interface{}) (interface{}, error) {
	return resource.Param(p, p.s.Info.Params(), key)
}

func (p *pageMeta) Params() maps.Params {
	return p.params
}

func (p *pageMeta) Path() string {
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
	return p.Kind() == page.KindSection
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

func (p *pageMeta) Sitemap() config.Sitemap {
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
		b1.cascade = make(map[string]interface{})
	}
	if b2 != nil && b2.cascade != nil {
		for k, v := range b2.cascade {
			if _, found := b1.cascade[k]; !found {
				b1.cascade[k] = v
			}
		}
	}
}

func (pm *pageMeta) setMetadata(parentBucket *pagesMapBucket, p *pageState, frontmatter map[string]interface{}) error {
	pm.params = make(maps.Params)

	if frontmatter == nil && (parentBucket == nil || parentBucket.cascade == nil) {
		return nil
	}

	if frontmatter != nil {
		// Needed for case insensitive fetching of params values
		maps.ToLower(frontmatter)
		if p.bucket != nil {
			// Check for any cascade define on itself.
			if cv, found := frontmatter["cascade"]; found {
				p.bucket.cascade = maps.ToStringMap(cv)
			}
		}
	} else {
		frontmatter = make(map[string]interface{})
	}

	var cascade map[string]interface{}

	if p.bucket != nil {
		if parentBucket != nil {
			// Merge missing keys from parent into this.
			pm.mergeBucketCascades(p.bucket, parentBucket)
		}
		cascade = p.bucket.cascade
	} else if parentBucket != nil {
		cascade = parentBucket.cascade
	}

	for k, v := range cascade {
		if _, found := frontmatter[k]; !found {
			frontmatter[k] = v
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
	if p.gitInfo != nil {
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
	}

	// Handle the date separately
	// TODO(bep) we need to "do more" in this area so this can be split up and
	// more easily tested without the Page, but the coupling is strong.
	err := pm.s.frontmatterHandler.HandleDates(descriptor)
	if err != nil {
		p.s.Log.ERROR.Printf("Failed to handle dates for page %q: %s", p.pathOrTitle(), err)
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
					p.s.Log.WARN.Printf(`Front matter in %q with the url %q with no leading / has what looks like the language prefix added. In Hugo 0.55 we added support for page relative URLs in front matter, no language prefix needed. Check the URL and consider to either add a leading / or remove the language prefix.`, p.pathOrTitle(), url)

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
				pm.buildConfig.Render = false
			}
		case "outputs":
			o := cast.ToStringSlice(v)
			if len(o) > 0 {
				// Output formats are exlicitly set in front matter, use those.
				outFormats, err := p.s.outputFormatsConfig.GetByNames(o...)

				if err != nil {
					p.s.Log.ERROR.Printf("Failed to resolve output formats: %s", err)
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
			p.m.sitemap = config.DecodeSitemap(p.s.siteCfg.sitemap, maps.ToStringMap(v))
			pm.params[loki] = p.m.sitemap
			sitemapSet = true
		case "iscjklanguage":
			isCJKLanguage = new(bool)
			*isCJKLanguage = cast.ToBool(v)
		case "translationkey":
			pm.translationKey = cast.ToString(v)
			pm.params[loki] = pm.translationKey
		case "resources":
			var resources []map[string]interface{}
			handled := true

			switch vv := v.(type) {
			case []map[interface{}]interface{}:
				for _, vvv := range vv {
					resources = append(resources, maps.ToStringMap(vvv))
				}
			case []map[string]interface{}:
				resources = append(resources, vv...)
			case []interface{}:
				for _, vvv := range vv {
					switch vvvv := vvv.(type) {
					case map[interface{}]interface{}:
						resources = append(resources, maps.ToStringMap(vvvv))
					case map[string]interface{}:
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
			case bool:
				pm.params[loki] = vv
			case string:
				pm.params[loki] = vv
			case int64, int32, int16, int8, int:
				pm.params[loki] = vv
			case float64, float32:
				pm.params[loki] = vv
			case time.Time:
				pm.params[loki] = vv
			default: // handle array of strings as well
				switch vvv := vv.(type) {
				case []interface{}:
					if len(vvv) > 0 {
						switch vvv[0].(type) {
						case map[interface{}]interface{}: // Proper parsing structured array from YAML based FrontMatter
							pm.params[loki] = vvv
						case map[string]interface{}: // Proper parsing structured array from JSON based FrontMatter
							pm.params[loki] = vvv
						case []interface{}:
							pm.params[loki] = vvv
						default:
							a := make([]string, len(vvv))
							for i, u := range vvv {
								a[i] = cast.ToString(u)
							}

							pm.params[loki] = a
						}
					} else {
						pm.params[loki] = []string{}
					}
				default:
					pm.params[loki] = vv
				}
			}
		}
	}

	if !sitemapSet {
		pm.sitemap = p.s.siteCfg.sitemap
	}

	pm.markup = p.s.ContentSpec.ResolveMarkup(pm.markup)

	if draft != nil && published != nil {
		pm.draft = *draft
		p.m.s.Log.WARN.Printf("page %q has both draft and published settings in its frontmatter. Using draft.", p.File().Filename())
	} else if draft != nil {
		pm.draft = *draft
	} else if published != nil {
		pm.draft = !*published
	}
	pm.params["draft"] = pm.draft

	if isCJKLanguage != nil {
		pm.isCJKLanguage = *isCJKLanguage
	} else if p.s.siteCfg.hasCJKLanguage && p.source.parsed != nil {
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
	return !p.buildConfig.Render
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
		case page.KindHome:
			p.title = p.s.Info.title
		case page.KindSection:
			var sectionName string
			if n != nil {
				sectionName = n.rootSection()
			} else {
				sectionName = p.sections[0]
			}

			sectionName = helpers.FirstUpper(sectionName)
			if p.s.Cfg.GetBool("pluralizeListTitles") {
				p.title = inflect.Pluralize(sectionName)
			} else {
				p.title = sectionName
			}
		case page.KindTaxonomy:
			// TODO(bep) improve
			key := p.sections[len(p.sections)-1]
			p.title = strings.Replace(p.s.titleFunc(key), "-", " ", -1)
		case page.KindTaxonomyTerm:
			p.title = p.s.titleFunc(p.sections[0])
		case kind404:
			p.title = "404 Page not found"

		}
	}

	if p.IsNode() {
		p.bundleType = files.ContentClassBranch
	} else {
		source := p.File()
		if fi, ok := source.(*fileInfo); ok {
			class := fi.FileInfo().Meta().Classifier()
			switch class {
			case files.ContentClassBranch, files.ContentClassLeaf:
				p.bundleType = class
			}
		}
	}

	if !p.f.IsZero() {
		var renderingConfigOverrides map[string]interface{}
		bfParam := getParamToLower(p, "blackfriday")
		if bfParam != nil {
			renderingConfigOverrides = maps.ToStringMap(bfParam)
		}

		p.renderingConfigOverrides = renderingConfigOverrides

	}

	return nil

}

func (p *pageMeta) newContentConverter(ps *pageState, markup string, renderingConfigOverrides map[string]interface{}) (converter.Converter, error) {
	if ps == nil {
		panic("no Page provided")
	}
	cp := p.s.ContentSpec.Converters.Get(markup)
	if cp == nil {
		return converter.NopConverter, errors.Errorf("no content renderer found for markup %q", p.markup)
	}

	cpp, err := cp.New(
		converter.DocumentContext{
			Document:        newPageForRenderHook(ps),
			DocumentID:      p.f.UniqueID(),
			DocumentName:    p.f.Path(),
			ConfigOverrides: renderingConfigOverrides,
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

	return m.s.outputFormats[m.Kind()]
}

func (p *pageMeta) Slug() string {
	return p.urlPaths.Slug
}

func getParam(m resource.ResourceParamsProvider, key string, stringToLower bool) interface{} {
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

func getParamToLower(m resource.ResourceParamsProvider, key string) interface{} {
	return getParam(m, key, true)
}
