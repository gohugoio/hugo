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
	"github.com/gohugoio/hugo/hugofs/files"

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

// AlternativeOutputFormatsProvider provides alternative output formats for a
// Page.
type AlternativeOutputFormatsProvider interface {
	// AlternativeOutputFormats gives the alternative output formats for the
	// current output.
	// Note that we use the term "alternative" and not "alternate" here, as it
	// does not necessarily replace the other format, it is an alternative representation.
	AlternativeOutputFormats() OutputFormats
}

// AuthorProvider provides author information.
type AuthorProvider interface {
	Author() Author
	Authors() AuthorList
}

// ChildCareProvider provides accessors to child resources.
type ChildCareProvider interface {
	Pages() Pages

	// RegularPages returns a list of pages of kind 'Page'.
	// In Hugo 0.57 we changed the Pages method so it returns all page
	// kinds, even sections. If you want the old behaviour, you can
	// use RegularPages.
	RegularPages() Pages

	// RegularPagesRecursive returns all regular pages below the current
	// section.
	RegularPagesRecursive() Pages

	Resources() resource.Resources
}

// ContentProvider provides the content related values for a Page.
type ContentProvider interface {
	Content() (interface{}, error)
	Plain() string
	PlainWords() []string
	Summary() template.HTML
	Truncated() bool
	FuzzyWordCount() int
	WordCount() int
	ReadingTime() int
	Len() int
}

// FileProvider provides the source file.
type FileProvider interface {
	File() source.File
}

// GetPageProvider provides the GetPage method.
type GetPageProvider interface {
	// GetPage looks up a page for the given ref.
	//    {{ with .GetPage "blog" }}{{ .Title }}{{ end }}
	//
	// This will return nil when no page could be found, and will return
	// an error if the ref is ambiguous.
	GetPage(ref string) (Page, error)
}

// GitInfoProvider provides Git info.
type GitInfoProvider interface {
	GitInfo() *gitmap.GitInfo
}

// InSectionPositioner provides section navigation.
type InSectionPositioner interface {
	NextInSection() Page
	PrevInSection() Page
}

// InternalDependencies is considered an internal interface.
type InternalDependencies interface {
	GetRelatedDocsHandler() *RelatedDocsHandler
}

// OutputFormatsProvider provides the OutputFormats of a Page.
type OutputFormatsProvider interface {
	OutputFormats() OutputFormats
}

// Page is the core interface in Hugo.
type Page interface {
	ContentProvider
	TableOfContentsProvider
	PageWithoutContent
}

// PageMetaProvider provides page metadata, typically provided via front matter.
type PageMetaProvider interface {
	// The 4 page dates
	resource.Dated

	// Aliases forms the base for redirects generation.
	Aliases() []string

	// BundleType returns the bundle type: "leaf", "branch" or an empty string if it is none.
	// See https://gohugo.io/content-management/page-bundles/
	BundleType() files.ContentClass

	// A configured description.
	Description() string

	// Whether this is a draft. Will only be true if run with the --buildDrafts (-D) flag.
	Draft() bool

	// IsHome returns whether this is the home page.
	IsHome() bool

	// Configured keywords.
	Keywords() []string

	// The Page Kind. One of page, home, section, taxonomy, taxonomyTerm.
	Kind() string

	// The configured layout to use to render this page. Typically set in front matter.
	Layout() string

	// The title used for links.
	LinkTitle() string

	// IsNode returns whether this is an item of one of the list types in Hugo,
	// i.e. not a regular content
	IsNode() bool

	// IsPage returns whether this is a regular content
	IsPage() bool

	// Param looks for a param in Page and then in Site config.
	Param(key interface{}) (interface{}, error)

	// Path gets the relative path, including file name and extension if relevant,
	// to the source of this Page. It will be relative to any content root.
	Path() string

	// The slug, typically defined in front matter.
	Slug() string

	// This page's language code. Will be the same as the site's.
	Lang() string

	// IsSection returns whether this is a section
	IsSection() bool

	// Section returns the first path element below the content root.
	Section() string

	// Returns a slice of sections (directories if it's a file) to this
	// Page.
	SectionsEntries() []string

	// SectionsPath is SectionsEntries joined with a /.
	SectionsPath() string

	// Sitemap returns the sitemap configuration for this page.
	Sitemap() config.Sitemap

	// Type is a discriminator used to select layouts etc. It is typically set
	// in front matter, but will fall back to the root section.
	Type() string

	// The configured weight, used as the first sort value in the default
	// page sort if non-zero.
	Weight() int
}

// PageRenderProvider provides a way for a Page to render content.
type PageRenderProvider interface {
	Render(layout ...string) (template.HTML, error)
	RenderString(args ...interface{}) (template.HTML, error)
}

// PageWithoutContent is the Page without any of the content methods.
type PageWithoutContent interface {
	RawContentProvider
	resource.Resource
	PageMetaProvider
	resource.LanguageProvider

	// For pages backed by a file.
	FileProvider

	GitInfoProvider

	// Output formats
	OutputFormatsProvider
	AlternativeOutputFormatsProvider

	// Tree navigation
	ChildCareProvider
	TreeProvider

	// Horizontal navigation
	InSectionPositioner
	PageRenderProvider
	PaginatorProvider
	Positioner
	navigation.PageMenusProvider

	// TODO(bep)
	AuthorProvider

	// Page lookups/refs
	GetPageProvider
	RefProvider

	resource.TranslationKeyProvider
	TranslationsProvider

	SitesProvider

	// Helper methods
	ShortcodeInfoProvider
	compare.Eqer
	maps.Scratcher
	RelatedKeywordsProvider

	// GetTerms gets the terms of a given taxonomy,
	// e.g. GetTerms("categories")
	GetTerms(taxonomy string) Pages

	DeprecatedWarningPageMethods
}

// Positioner provides next/prev navigation.
type Positioner interface {
	Next() Page
	Prev() Page

	// Deprecated: Use Prev. Will be removed in Hugo 0.57
	PrevPage() Page

	// Deprecated: Use Next. Will be removed in Hugo 0.57
	NextPage() Page
}

// RawContentProvider provides the raw, unprocessed content of the page.
type RawContentProvider interface {
	RawContent() string
}

// RefProvider provides the methods needed to create reflinks to pages.
type RefProvider interface {
	Ref(argsm map[string]interface{}) (string, error)
	RefFrom(argsm map[string]interface{}, source interface{}) (string, error)
	RelRef(argsm map[string]interface{}) (string, error)
	RelRefFrom(argsm map[string]interface{}, source interface{}) (string, error)
}

// RelatedKeywordsProvider allows a Page to be indexed.
type RelatedKeywordsProvider interface {
	// Make it indexable as a related.Document
	RelatedKeywords(cfg related.IndexConfig) ([]related.Keyword, error)
}

// ShortcodeInfoProvider provides info about the shortcodes in a Page.
type ShortcodeInfoProvider interface {
	// HasShortcode return whether the page has a shortcode with the given name.
	// This method is mainly motivated with the Hugo Docs site's need for a list
	// of pages with the `todo` shortcode in it.
	HasShortcode(name string) bool
}

// SitesProvider provide accessors to get sites.
type SitesProvider interface {
	Site() Site
	Sites() Sites
}

// TableOfContentsProvider provides the table of contents for a Page.
type TableOfContentsProvider interface {
	TableOfContents() template.HTML
}

// TranslationsProvider provides access to any translations.
type TranslationsProvider interface {

	// IsTranslated returns whether this content file is translated to
	// other language(s).
	IsTranslated() bool

	// AllTranslations returns all translations, including the current Page.
	AllTranslations() Pages

	// Translations returns the translations excluding the current Page.
	Translations() Pages
}

// TreeProvider provides section tree navigation.
type TreeProvider interface {

	// IsAncestor returns whether the current page is an ancestor of the given
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsAncestor(other interface{}) (bool, error)

	// CurrentSection returns the page's current section or the page itself if home or a section.
	// Note that this will return nil for pages that is not regular, home or section pages.
	CurrentSection() Page

	// IsDescendant returns whether the current page is a descendant of the given
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsDescendant(other interface{}) (bool, error)

	// FirstSection returns the section on level 1 below home, e.g. "/docs".
	// For the home page, this will return itself.
	FirstSection() Page

	// InSection returns whether the given page is in the current section.
	// Note that this will always return false for pages that are
	// not either regular, home or section pages.
	InSection(other interface{}) (bool, error)

	// Parent returns a section's parent section or a page's section.
	// To get a section's subsections, see Page's Sections method.
	Parent() Page

	// Sections returns this section's subsections, if any.
	// Note that for non-sections, this method will always return an empty list.
	Sections() Pages

	// Page returns a reference to the Page itself, kept here mostly
	// for legacy reasons.
	Page() Page
}

// DeprecatedWarningPageMethods lists deprecated Page methods that will trigger
// a WARNING if invoked.
// This was added in Hugo 0.55.
type DeprecatedWarningPageMethods interface {
	source.FileWithoutOverlap
	DeprecatedWarningPageMethods1
}

type DeprecatedWarningPageMethods1 interface {
	IsDraft() bool
	Hugo() hugo.Info
	LanguagePrefix() string
	GetParam(key string) interface{}
	RSSLink() template.URL
	URL() string
}

// Move here to trigger ERROR instead of WARNING.
// TODO(bep) create wrappers and put into the Page once it has some methods.
type DeprecatedErrorPageMethods interface {
}
