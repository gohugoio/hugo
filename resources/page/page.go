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

// Package page contains the core interfaces and types for the Page resource,
// a core component in Hugo.
package page

import (
	"html/template"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/compare"

	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/source"
)

// Clear clears any global package state.
func Clear() error {
	spc.clear()
	return nil
}

type AuthorProvider interface {
	Author() Author
	Authors() AuthorList
}

type ChildCareProvider interface {
	Pages() Pages
	Resources() resource.Resources
}

type ContentProvider interface {
	Content() (interface{}, error)
	FuzzyWordCount() int
	Len() int
	Plain() string
	PlainWords() []string
	ReadingTime() int
	Summary() template.HTML
	TableOfContents() template.HTML
	Truncated() bool
	WordCount() int
}

type FileProvider interface {
	File() source.File
}

type GetPageProvider interface {
	// GetPage looks up a page for the given ref.
	//    {{ with .GetPage "blog" }}{{ .Title }}{{ end }}
	//
	// This will return nil when no page could be found, and will return
	// an error if the ref is ambiguous.
	GetPage(ref string) (Page, error)
}

type GitInfoProvider interface {
	GitInfo() *gitmap.GitInfo
}

type InSectionPositioner interface {
	NextInSection() Page
	PrevInSection() Page
}

// InternalDependencies is considered an internal interface.
type InternalDependencies interface {
	GetRelatedDocsHandler() *RelatedDocsHandler
}

type OutputFormatsProvider interface {
	OutputFormats() OutputFormats
}

type AlternativeOutputFormatsProvider interface {
	// AlternativeOutputFormats gives the alternative output formats for the
	// current output.
	// Note that we use the term "alternative" and not "alternate" here, as it
	// does not necessarily replace the other format, it is an alternative representation.
	AlternativeOutputFormats() OutputFormats
}

type Page interface {
	ContentProvider
	PageWithoutContent
}

// Page metadata, typically provided via front matter.
type PageMetaProvider interface {
	resource.Dated

	Aliases() []string

	// BundleType returns the bundle type: "leaf", "branch" or an empty string if it is none.
	// See https://gohugo.io/content-management/page-bundles/
	BundleType() string

	Draft() bool

	// IsHome returns whether this is the home
	IsHome() bool

	// IsNode returns whether this is an item of one of the list types in Hugo,
	// i.e. not a regular content
	IsNode() bool

	// IsPage returns whether this is a regular content
	IsPage() bool

	// IsSection returns whether this is a section
	IsSection() bool

	// The layout to use to render this page. Typically set in front matter.
	Layout() string

	Description() string

	Keywords() []string

	Kind() string
	LinkTitle() string
	Param(key interface{}) (interface{}, error)
	Path() string
	Slug() string

	// Section returns the first path element below the content root.
	Section() string

	// TODO(bep) page name
	SectionsEntries() []string
	SectionsPath() string

	Sitemap() config.Sitemap

	Type() string
	Weight() int
}

type PageRenderProvider interface {
	Render(layout ...string) template.HTML
}

type ShortcodeInfoProvider interface {
	// HasShortcode return whether the page has a shortcode with the given name.
	// This method is mainly motivated with the Hugo Docs site's need for a list
	// of pages with the `todo` shortcode in it.
	HasShortcode(name string) bool
}

type PageWithoutContent interface {
	AuthorProvider
	ChildCareProvider
	FileProvider
	GetPageProvider
	InSectionPositioner
	OutputFormatsProvider
	AlternativeOutputFormatsProvider
	PageMetaProvider
	PageRenderProvider
	PaginatorProvider
	Positioner
	RawContentProvider
	RefProvider
	SitesProvider
	TODOProvider
	TranslationsProvider
	TreeProvider
	compare.Eqer
	maps.Scratcher
	ShortcodeInfoProvider
	navigation.PageMenusProvider
	resource.LanguageProvider
	resource.Resource
	resource.TranslationKeyProvider
}

type Positioner interface {
	Next() Page

	NextPage() Page
	Prev() Page

	// TODO(bep) deprecate these 2
	PrevPage() Page
}

type RawContentProvider interface {
	RawContent() string
}

type RefProvider interface {
	Ref(argsm map[string]interface{}) (string, error)
	RefFrom(argsm map[string]interface{}, source interface{}) (string, error)
	RelRef(argsm map[string]interface{}) (string, error)
	RelRefFrom(argsm map[string]interface{}, source interface{}) (string, error)
}

type SitesProvider interface {
	Site() Site
	Sites() Sites
}

type TODOProvider interface {
	pageAddons3

	// Make it indexable as a related.Document
	SearchKeywords(cfg related.IndexConfig) ([]related.Keyword, error)

	// See deprecated file Section() string

	SourceRef() string
}

//
// TranslationProvider provides translated versions of a Page.
type TranslationProvider interface {
}

type TranslationsProvider interface {

	// AllTranslations returns all translations, including the current Page.
	AllTranslations() Pages

	// IsTranslated returns whether this content file is translated to
	// other language(s).
	IsTranslated() bool

	// Translations returns the translations excluding the current Page.
	Translations() Pages
}

type TreeProvider interface {

	// CurrentSection returns the page's current section or the page itself if home or a section.
	// Note that this will return nil for pages that is not regular, home or section pages.
	CurrentSection() Page

	// FirstSection returns the section on level 1 below home, e.g. "/docs".
	// For the home page, this will return itself.
	FirstSection() Page

	// InSection returns whether the given page is in the current section.
	// Note that this will always return false for pages that are
	// not either regular, home or section pages.
	InSection(other interface{}) (bool, error)

	// IsAncestor returns whether the current page is an ancestor of the given
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsAncestor(other interface{}) (bool, error)

	// IsDescendant returns whether the current page is a descendant of the given
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsDescendant(other interface{}) (bool, error)

	// Parent returns a section's parent section or a page's section.
	// To get a section's subsections, see Page's Sections method.
	Parent() Page

	// Sections returns this section's subsections, if any.
	// Note that for non-sections, this method will always return an empty list.
	Sections() Pages
}

type deprecatedPageMethods interface {
	source.FileWithoutOverlap

	Hugo() hugo.Info // use hugo

	IsDraft() bool // => Draft

	LanguagePrefix() string // Use Site.LanguagePrefix

	RSSLink() template.URL // `Use the Output Format's link, e.g. something like: {{ with .OutputFormats.Get "RSS" }}{{ . RelPermalink }}{{ end }}`

	URL() string // => .Permalink / .RelPermalink

}

type pageAddons3 interface {

	// TODO(bep) page consider what to do.

	deprecatedPageMethods

	// TODO(bep) page remove/deprecate (use Param)
	GetParam(key string) interface{}
}
