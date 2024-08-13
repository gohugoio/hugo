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
	"context"
	"html/template"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
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
	// Deprecated: Use taxonomies instead.
	Author() Author
	// Deprecated: Use taxonomies instead.
	Authors() AuthorList
}

// ChildCareProvider provides accessors to child resources.
type ChildCareProvider interface {
	// Pages returns a list of pages of all kinds.
	Pages() Pages

	// RegularPages returns a list of pages of kind 'Page'.
	RegularPages() Pages

	// RegularPagesRecursive returns all regular pages below the current
	// section.
	RegularPagesRecursive() Pages

	// Resources returns a list of all resources.
	Resources() resource.Resources
}

type MarkupProvider interface {
	Markup(opts ...any) Markup
}

// ContentProvider provides the content related values for a Page.
type ContentProvider interface {
	Content(context.Context) (any, error)

	// ContentWithoutSummary returns the Page Content stripped of the summary.
	ContentWithoutSummary(ctx context.Context) (template.HTML, error)

	// Plain returns the Page Content stripped of HTML markup.
	Plain(context.Context) string

	// PlainWords returns a string slice from splitting Plain using https://pkg.go.dev/strings#Fields.
	PlainWords(context.Context) []string

	// Summary returns a generated summary of the content.
	// The breakpoint can be set manually by inserting a summary separator in the source file.
	Summary(context.Context) template.HTML

	// Truncated returns whether the Summary  is truncated or not.
	Truncated(context.Context) bool

	// FuzzyWordCount returns the approximate number of words in the content.
	FuzzyWordCount(context.Context) int

	// WordCount returns the number of words in the content.
	WordCount(context.Context) int

	// ReadingTime returns the reading time based on the length of plain text.
	ReadingTime(context.Context) int

	// Len returns the length of the content.
	// This is for internal use only.
	Len(context.Context) int
}

// ContentRenderer provides the content rendering methods for some content.
type ContentRenderer interface {
	// ParseAndRenderContent renders the given content.
	// For internal use only.
	ParseAndRenderContent(ctx context.Context, content []byte, enableTOC bool) (converter.ResultRender, error)
	// For internal use only.
	ParseContent(ctx context.Context, content []byte) (converter.ResultParse, bool, error)
	// For internal use only.
	RenderContent(ctx context.Context, content []byte, doc any) (converter.ResultRender, bool, error)
}

// FileProvider provides the source file.
type FileProvider interface {
	// File returns the source file for this Page,
	// or a zero File if this Page is not backed by a file.
	File() *source.File
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
	// GitInfo returns the Git info for this object.
	GitInfo() source.GitInfo
	// CodeOwners returns the code owners for this object.
	CodeOwners() []string
}

// InSectionPositioner provides section navigation.
type InSectionPositioner interface {
	// NextInSection returns the next page in the same section.
	NextInSection() Page
	// PrevInSection returns the previous page in the same section.
	PrevInSection() Page
}

// InternalDependencies is considered an internal interface.
type InternalDependencies interface {
	// GetRelatedDocsHandler is for internal use only.
	GetRelatedDocsHandler() *RelatedDocsHandler
}

// OutputFormatsProvider provides the OutputFormats of a Page.
type OutputFormatsProvider interface {
	// OutputFormats returns the OutputFormats for this Page.
	OutputFormats() OutputFormats
}

// PageProvider provides access to a Page.
// Implemented by shortcodes and others.
type PageProvider interface {
	Page() Page
}

// Page is the core interface in Hugo and what you get as the top level data context in your templates.
type Page interface {
	MarkupProvider
	ContentProvider
	TableOfContentsProvider
	PageWithoutContent
}

type PageFragment interface {
	resource.ResourceLinksProvider
	resource.ResourceNameTitleProvider
}

// PageMetaProvider provides page metadata, typically provided via front matter.
type PageMetaProvider interface {
	// The 4 page dates
	resource.Dated

	// Aliases forms the base for redirects generation.
	Aliases() []string

	// BundleType returns the bundle type: `leaf`, `branch` or an empty string.
	BundleType() string

	// A configured description.
	Description() string

	// Whether this is a draft. Will only be true if run with the --buildDrafts (-D) flag.
	Draft() bool

	// IsHome returns whether this is the home page.
	IsHome() bool

	// Configured keywords.
	Keywords() []string

	// The Page Kind. One of page, home, section, taxonomy, term.
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
	Param(key any) (any, error)

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

	// Sitemap returns the sitemap configuration for this page.
	// This is for internal use only.
	Sitemap() config.SitemapConfig

	// Type is a discriminator used to select layouts etc. It is typically set
	// in front matter, but will fall back to the root section.
	Type() string

	// The configured weight, used as the first sort value in the default
	// page sort if non-zero.
	Weight() int
}

// PageMetaInternalProvider provides internal page metadata.
type PageMetaInternalProvider interface {
	// This is for internal use only.
	PathInfo() *paths.Path
}

// PageRenderProvider provides a way for a Page to render content.
type PageRenderProvider interface {
	// Render renders the given layout with this Page as context.
	Render(ctx context.Context, layout ...string) (template.HTML, error)
	// RenderString renders the first value in args with the content renderer defined
	// for this Page.
	// It takes an optional map as a second argument:
	//
	// display (“inline”):
	// - inline or block. If inline (default), surrounding <p></p> on short snippets will be trimmed.
	// markup (defaults to the Page’s markup)
	RenderString(ctx context.Context, args ...any) (template.HTML, error)
}

// PageWithoutContent is the Page without any of the content methods.
type PageWithoutContent interface {
	RawContentProvider
	RenderShortcodesProvider
	resource.Resource
	PageMetaProvider
	PageMetaInternalProvider
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

	// Scratch returns a Scratch that can be used to store temporary state.
	// Note that this Scratch gets reset on server rebuilds. See Store() for a variant that survives.
	maps.Scratcher

	// Store returns a Scratch that can be used to store temporary state.
	// In contrast to Scratch(), this Scratch is not reset on server rebuilds.
	Store() *maps.Scratch

	RelatedKeywordsProvider

	// GetTerms gets the terms of a given taxonomy,
	// e.g. GetTerms("categories")
	GetTerms(taxonomy string) Pages

	// HeadingsFiltered returns the headings for this page when a filter is set.
	// This is currently only triggered with the Related content feature
	// and the "fragments" type of index.
	HeadingsFiltered(context.Context) tableofcontents.Headings
}

// Positioner provides next/prev navigation.
type Positioner interface {
	// Next points up to the next regular page (sorted by Hugo’s default sort).
	Next() Page
	// Prev points down to the previous regular page (sorted by Hugo’s default sort).
	Prev() Page

	// Deprecated: Use Prev. Will be removed in Hugo 0.57
	PrevPage() Page

	// Deprecated: Use Next. Will be removed in Hugo 0.57
	NextPage() Page
}

// RawContentProvider provides the raw, unprocessed content of the page.
type RawContentProvider interface {
	// RawContent returns the raw, unprocessed content of the page excluding any front matter.
	RawContent() string
}

type RenderShortcodesProvider interface {
	// RenderShortcodes returns RawContent with any shortcodes rendered.
	RenderShortcodes(context.Context) (template.HTML, error)
}

// RefProvider provides the methods needed to create reflinks to pages.
type RefProvider interface {
	// Ref returns an absolute URl to a page.
	Ref(argsm map[string]any) (string, error)

	// RefFrom is for internal use only.
	RefFrom(argsm map[string]any, source any) (string, error)

	// RelRef returns a relative URL to a page.
	RelRef(argsm map[string]any) (string, error)

	// RelRefFrom is for internal use only.
	RelRefFrom(argsm map[string]any, source any) (string, error)
}

// RelatedKeywordsProvider allows a Page to be indexed.
type RelatedKeywordsProvider interface {
	// Make it indexable as a related.Document
	// RelatedKeywords is meant for internal usage only.
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
	// Site returns the current site.
	Site() Site
	// Sites returns all sites.
	Sites() Sites
}

// TableOfContentsProvider provides the table of contents for a Page.
type TableOfContentsProvider interface {
	// TableOfContents returns the table of contents for the page rendered as HTML.
	TableOfContents(context.Context) template.HTML

	// Fragments returns the fragments for this page.
	Fragments(context.Context) *tableofcontents.Fragments
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
	// IsAncestor returns whether the current page is an ancestor of other.
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsAncestor(other any) bool

	// CurrentSection returns the page's current section or the page itself if home or a section.
	// Note that this will return nil for pages that is not regular, home or section pages.
	CurrentSection() Page

	// IsDescendant returns whether the current page is a descendant of other.
	// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
	IsDescendant(other any) bool

	// FirstSection returns the section on level 1 below home, e.g. "/docs".
	// For the home page, this will return itself.
	FirstSection() Page

	// InSection returns whether other is in the current section.
	// Note that this will always return false for pages that are
	// not either regular, home or section pages.
	InSection(other any) bool

	// Parent returns a section's parent section or a page's section.
	// To get a section's subsections, see Page's Sections method.
	Parent() Page

	// Ancestors returns the ancestors of each page
	Ancestors() Pages

	// Sections returns this section's subsections, if any.
	// Note that for non-sections, this method will always return an empty list.
	Sections() Pages

	// Page returns a reference to the Page itself, kept here mostly
	// for legacy reasons.
	Page() Page

	// Returns a slice of sections (directories if it's a file) to this
	// Page.
	SectionsEntries() []string

	// SectionsPath is SectionsEntries joined with a /.
	SectionsPath() string
}

// PageWithContext is a Page with a context.Context.
type PageWithContext struct {
	Page
	Ctx context.Context
}

func (p PageWithContext) Content() (any, error) {
	return p.Page.Content(p.Ctx)
}

func (p PageWithContext) Plain() string {
	return p.Page.Plain(p.Ctx)
}

func (p PageWithContext) PlainWords() []string {
	return p.Page.PlainWords(p.Ctx)
}

func (p PageWithContext) Summary() template.HTML {
	return p.Page.Summary(p.Ctx)
}

func (p PageWithContext) Truncated() bool {
	return p.Page.Truncated(p.Ctx)
}

func (p PageWithContext) FuzzyWordCount() int {
	return p.Page.FuzzyWordCount(p.Ctx)
}

func (p PageWithContext) WordCount() int {
	return p.Page.WordCount(p.Ctx)
}

func (p PageWithContext) ReadingTime() int {
	return p.Page.ReadingTime(p.Ctx)
}

func (p PageWithContext) Len() int {
	return p.Page.Len(p.Ctx)
}
