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
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/navigation"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	NopPage                 Page            = new(nopPage)
	NopContentRenderer      ContentRenderer = new(nopContentRenderer)
	NopCPageContentRenderer                 = struct {
		OutputFormatPageContentProvider
		ContentRenderer
	}{
		NopPage,
		NopContentRenderer,
	}
	NilPage *nopPage
)

// PageNop implements Page, but does nothing.
type nopPage int

var noOpPathInfo = paths.Parse(files.ComponentFolderContent, "no-op.md")

func (p *nopPage) Err() resource.ResourceError {
	return nil
}

func (p *nopPage) Aliases() []string {
	return nil
}

func (p *nopPage) Sitemap() config.SitemapConfig {
	return config.SitemapConfig{}
}

func (p *nopPage) Layout() string {
	return ""
}

func (p *nopPage) RSSLink() template.URL {
	return ""
}

func (p *nopPage) Author() Author {
	return Author{}
}

func (p *nopPage) Authors() AuthorList {
	return nil
}

func (p *nopPage) AllTranslations() Pages {
	return nil
}

func (p *nopPage) LanguagePrefix() string {
	return ""
}

func (p *nopPage) AlternativeOutputFormats() OutputFormats {
	return nil
}

func (p *nopPage) BaseFileName() string {
	return ""
}

func (p *nopPage) BundleType() string {
	return ""
}

func (p *nopPage) Content(context.Context) (any, error) {
	return "", nil
}

func (p *nopPage) ContentBaseName() string {
	return ""
}

func (p *nopPage) CurrentSection() Page {
	return nil
}

func (p *nopPage) Data() any {
	return nil
}

func (p *nopPage) Date() (t time.Time) {
	return
}

func (p *nopPage) Description() string {
	return ""
}

func (p *nopPage) RefFrom(argsm map[string]any, source any) (string, error) {
	return "", nil
}

func (p *nopPage) RelRefFrom(argsm map[string]any, source any) (string, error) {
	return "", nil
}

func (p *nopPage) Dir() string {
	return ""
}

func (p *nopPage) Draft() bool {
	return false
}

func (p *nopPage) Eq(other any) bool {
	return p == other
}

func (p *nopPage) ExpiryDate() (t time.Time) {
	return
}

func (p *nopPage) Ext() string {
	return ""
}

func (p *nopPage) Extension() string {
	return ""
}

func (p *nopPage) File() *source.File {
	return nil
}

func (p *nopPage) FileInfo() hugofs.FileMetaInfo {
	return nil
}

func (p *nopPage) Filename() string {
	return ""
}

func (p *nopPage) FirstSection() Page {
	return nil
}

func (p *nopPage) FuzzyWordCount(context.Context) int {
	return 0
}

func (p *nopPage) GetPage(ref string) (Page, error) {
	return nil, nil
}

func (p *nopPage) GetParam(key string) any {
	return nil
}

func (p *nopPage) GetTerms(taxonomy string) Pages {
	return nil
}

func (p *nopPage) GitInfo() source.GitInfo {
	return source.GitInfo{}
}

func (p *nopPage) CodeOwners() []string {
	return nil
}

func (p *nopPage) HasMenuCurrent(menuID string, me *navigation.MenuEntry) bool {
	return false
}

func (p *nopPage) HasShortcode(name string) bool {
	return false
}

func (p *nopPage) Hugo() (h hugo.HugoInfo) {
	return
}

func (p *nopPage) InSection(other any) bool {
	return false
}

func (p *nopPage) IsAncestor(other any) bool {
	return false
}

func (p *nopPage) IsDescendant(other any) bool {
	return false
}

func (p *nopPage) IsDraft() bool {
	return false
}

func (p *nopPage) IsHome() bool {
	return false
}

func (p *nopPage) IsMenuCurrent(menuID string, inme *navigation.MenuEntry) bool {
	return false
}

func (p *nopPage) IsNode() bool {
	return false
}

func (p *nopPage) IsPage() bool {
	return false
}

func (p *nopPage) IsSection() bool {
	return false
}

func (p *nopPage) IsTranslated() bool {
	return false
}

func (p *nopPage) Keywords() []string {
	return nil
}

func (p *nopPage) Kind() string {
	return ""
}

func (p *nopPage) Lang() string {
	return ""
}

func (p *nopPage) Language() *langs.Language {
	return nil
}

func (p *nopPage) Lastmod() (t time.Time) {
	return
}

func (p *nopPage) Len(context.Context) int {
	return 0
}

func (p *nopPage) LinkTitle() string {
	return ""
}

func (p *nopPage) LogicalName() string {
	return ""
}

func (p *nopPage) MediaType() (m media.Type) {
	return
}

func (p *nopPage) Menus() (m navigation.PageMenus) {
	return
}

func (p *nopPage) Name() string {
	return ""
}

func (p *nopPage) Next() Page {
	return nil
}

func (p *nopPage) OutputFormats() OutputFormats {
	return nil
}

func (p *nopPage) Pages() Pages {
	return nil
}

func (p *nopPage) RegularPages() Pages {
	return nil
}

func (p *nopPage) RegularPagesRecursive() Pages {
	return nil
}

func (p *nopPage) Paginate(seq any, options ...any) (*Pager, error) {
	return nil, nil
}

func (p *nopPage) Paginator(options ...any) (*Pager, error) {
	return nil, nil
}

func (p *nopPage) Param(key any) (any, error) {
	return nil, nil
}

func (p *nopPage) Params() maps.Params {
	return nil
}

func (p *nopPage) Page() Page {
	return p
}

func (p *nopPage) Parent() Page {
	return nil
}

func (p *nopPage) Ancestors() Pages {
	return nil
}

func (p *nopPage) Path() string {
	return ""
}

func (p *nopPage) PathInfo() *paths.Path {
	return noOpPathInfo
}

func (p *nopPage) Permalink() string {
	return ""
}

func (p *nopPage) Plain(context.Context) string {
	return ""
}

func (p *nopPage) PlainWords(context.Context) []string {
	return nil
}

func (p *nopPage) Prev() Page {
	return nil
}

func (p *nopPage) PublishDate() (t time.Time) {
	return
}

func (p *nopPage) PrevInSection() Page {
	return nil
}

func (p *nopPage) NextInSection() Page {
	return nil
}

func (p *nopPage) PrevPage() Page {
	return nil
}

func (p *nopPage) NextPage() Page {
	return nil
}

func (p *nopPage) RawContent() string {
	return ""
}

func (p *nopPage) RenderShortcodes(ctx context.Context) (template.HTML, error) {
	return "", nil
}

func (p *nopPage) ReadingTime(context.Context) int {
	return 0
}

func (p *nopPage) Ref(argsm map[string]any) (string, error) {
	return "", nil
}

func (p *nopPage) RelPermalink() string {
	return ""
}

func (p *nopPage) RelRef(argsm map[string]any) (string, error) {
	return "", nil
}

func (p *nopPage) Render(ctx context.Context, layout ...string) (template.HTML, error) {
	return "", nil
}

func (p *nopPage) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	return "", nil
}

func (p *nopPage) ResourceType() string {
	return ""
}

func (p *nopPage) Resources() resource.Resources {
	return nil
}

func (p *nopPage) Scratch() *maps.Scratch {
	return nil
}

func (p *nopPage) Store() *maps.Scratch {
	return nil
}

func (p *nopPage) RelatedKeywords(cfg related.IndexConfig) ([]related.Keyword, error) {
	return nil, nil
}

func (p *nopPage) Section() string {
	return ""
}

func (p *nopPage) Sections() Pages {
	return nil
}

func (p *nopPage) SectionsEntries() []string {
	return nil
}

func (p *nopPage) SectionsPath() string {
	return ""
}

func (p *nopPage) Site() Site {
	return nil
}

func (p *nopPage) Sites() Sites {
	return nil
}

func (p *nopPage) Slug() string {
	return ""
}

func (p *nopPage) String() string {
	return "nopPage"
}

func (p *nopPage) Summary(context.Context) template.HTML {
	return ""
}

func (p *nopPage) TableOfContents(context.Context) template.HTML {
	return ""
}

func (p *nopPage) Title() string {
	return ""
}

func (p *nopPage) TranslationBaseName() string {
	return ""
}

func (p *nopPage) TranslationKey() string {
	return ""
}

func (p *nopPage) Translations() Pages {
	return nil
}

func (p *nopPage) Truncated(context.Context) bool {
	return false
}

func (p *nopPage) Type() string {
	return ""
}

func (p *nopPage) URL() string {
	return ""
}

func (p *nopPage) UniqueID() string {
	return ""
}

func (p *nopPage) Weight() int {
	return 0
}

func (p *nopPage) WordCount(context.Context) int {
	return 0
}

func (p *nopPage) Fragments(context.Context) *tableofcontents.Fragments {
	return nil
}

func (p *nopPage) HeadingsFiltered(context.Context) tableofcontents.Headings {
	return nil
}

type nopContentRenderer int

func (r *nopContentRenderer) ParseAndRenderContent(ctx context.Context, content []byte, renderTOC bool) (converter.ResultRender, error) {
	b := &bytes.Buffer{}
	return b, nil
}

func (r *nopContentRenderer) ParseContent(ctx context.Context, content []byte) (converter.ResultParse, bool, error) {
	return nil, false, nil
}

func (r *nopContentRenderer) RenderContent(ctx context.Context, content []byte, doc any) (converter.ResultRender, bool, error) {
	return nil, false, nil
}
