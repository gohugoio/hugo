// Copyright 2023 The Hugo Authors. All rights reserved.
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

package page

import (
	"context"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/navigation"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/related"

	"github.com/gohugoio/hugo/source"
)

var (
	_ resource.LengthProvider = (*testPage)(nil)
	_ Page                    = (*testPage)(nil)
)

var relatedDocsHandler = NewRelatedDocsHandler(related.DefaultConfig)

func newTestPage() *testPage {
	return newTestPageWithFile("/a/b/c.md")
}

func newTestPageWithFile(filename string) *testPage {
	filename = filepath.FromSlash(filename)
	file := source.NewTestFile(filename)

	l, err := langs.NewLanguage(
		"en",
		"en",
		"UTC",
		langs.LanguageConfig{
			LanguageName: "English",
		},
	)
	if err != nil {
		panic(err)
	}

	return &testPage{
		params: make(map[string]any),
		data:   make(map[string]any),
		file:   file,
		currentSection: &testPage{
			sectionEntries: []string{"a", "b", "c"},
		},
		site: testSite{l: l},
	}
}

type testPage struct {
	kind        string
	description string
	title       string
	linkTitle   string
	lang        string
	section     string
	site        testSite

	content string

	fuzzyWordCount int

	path string

	slug string

	// Dates
	date       time.Time
	lastMod    time.Time
	expiryDate time.Time
	pubDate    time.Time

	weight int

	params map[string]any
	data   map[string]any

	file source.File

	currentSection *testPage
	sectionEntries []string
}

func (p *testPage) Err() resource.ResourceError {
	return nil
}

func (p *testPage) Aliases() []string {
	panic("testpage: not implemented")
}

func (p *testPage) AllTranslations() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) AlternativeOutputFormats() OutputFormats {
	panic("testpage: not implemented")
}

func (p *testPage) Author() Author {
	return Author{}
}

func (p *testPage) Authors() AuthorList {
	return nil
}

func (p *testPage) BaseFileName() string {
	panic("testpage: not implemented")
}

func (p *testPage) BundleType() files.ContentClass {
	panic("testpage: not implemented")
}

func (p *testPage) Content(context.Context) (any, error) {
	panic("testpage: not implemented")
}

func (p *testPage) ContentBaseName() string {
	panic("testpage: not implemented")
}

func (p *testPage) CurrentSection() Page {
	return p.currentSection
}

func (p *testPage) Data() any {
	return p.data
}

func (p *testPage) Sitemap() config.SitemapConfig {
	return config.SitemapConfig{}
}

func (p *testPage) Layout() string {
	return ""
}

func (p *testPage) Date() time.Time {
	return p.date
}

func (p *testPage) Description() string {
	return ""
}

func (p *testPage) Dir() string {
	panic("testpage: not implemented")
}

func (p *testPage) Draft() bool {
	panic("testpage: not implemented")
}

func (p *testPage) Eq(other any) bool {
	return p == other
}

func (p *testPage) ExpiryDate() time.Time {
	return p.expiryDate
}

func (p *testPage) Ext() string {
	panic("testpage: not implemented")
}

func (p *testPage) Extension() string {
	panic("testpage: not implemented")
}

func (p *testPage) File() source.File {
	return p.file
}

func (p *testPage) FileInfo() hugofs.FileMetaInfo {
	panic("testpage: not implemented")
}

func (p *testPage) Filename() string {
	panic("testpage: not implemented")
}

func (p *testPage) FirstSection() Page {
	panic("testpage: not implemented")
}

func (p *testPage) FuzzyWordCount(context.Context) int {
	return p.fuzzyWordCount
}

func (p *testPage) GetPage(ref string) (Page, error) {
	panic("testpage: not implemented")
}

func (p *testPage) GetPageWithTemplateInfo(info tpl.Info, ref string) (Page, error) {
	panic("testpage: not implemented")
}

func (p *testPage) GetParam(key string) any {
	panic("testpage: not implemented")
}

func (p *testPage) GetTerms(taxonomy string) Pages {
	panic("testpage: not implemented")
}

func (p *testPage) GetRelatedDocsHandler() *RelatedDocsHandler {
	return relatedDocsHandler
}

func (p *testPage) GitInfo() source.GitInfo {
	return source.GitInfo{}
}

func (p *testPage) CodeOwners() []string {
	return nil
}

func (p *testPage) HasMenuCurrent(menuID string, me *navigation.MenuEntry) bool {
	panic("testpage: not implemented")
}

func (p *testPage) HasShortcode(name string) bool {
	panic("testpage: not implemented")
}

func (p *testPage) Hugo() hugo.HugoInfo {
	panic("testpage: not implemented")
}

func (p *testPage) InSection(other any) (bool, error) {
	panic("testpage: not implemented")
}

func (p *testPage) IsAncestor(other any) (bool, error) {
	panic("testpage: not implemented")
}

func (p *testPage) IsDescendant(other any) (bool, error) {
	panic("testpage: not implemented")
}

func (p *testPage) IsDraft() bool {
	return false
}

func (p *testPage) IsHome() bool {
	panic("testpage: not implemented")
}

func (p *testPage) IsMenuCurrent(menuID string, inme *navigation.MenuEntry) bool {
	panic("testpage: not implemented")
}

func (p *testPage) IsNode() bool {
	panic("testpage: not implemented")
}

func (p *testPage) IsPage() bool {
	panic("testpage: not implemented")
}

func (p *testPage) IsSection() bool {
	panic("testpage: not implemented")
}

func (p *testPage) IsTranslated() bool {
	panic("testpage: not implemented")
}

func (p *testPage) Keywords() []string {
	return nil
}

func (p *testPage) Kind() string {
	return p.kind
}

func (p *testPage) Lang() string {
	return p.lang
}

func (p *testPage) Language() *langs.Language {
	panic("testpage: not implemented")
}

func (p *testPage) LanguagePrefix() string {
	return ""
}

func (p *testPage) Fragments(context.Context) *tableofcontents.Fragments {
	return nil
}

func (p *testPage) HeadingsFiltered(context.Context) tableofcontents.Headings {
	return nil
}

func (p *testPage) Lastmod() time.Time {
	return p.lastMod
}

func (p *testPage) Len(context.Context) int {
	return len(p.content)
}

func (p *testPage) LinkTitle() string {
	if p.linkTitle == "" {
		if p.title == "" {
			return p.path
		}
		return p.title
	}
	return p.linkTitle
}

func (p *testPage) LogicalName() string {
	panic("testpage: not implemented")
}

func (p *testPage) MediaType() media.Type {
	panic("testpage: not implemented")
}

func (p *testPage) Menus() navigation.PageMenus {
	return navigation.PageMenus{}
}

func (p *testPage) Name() string {
	panic("testpage: not implemented")
}

func (p *testPage) Next() Page {
	panic("testpage: not implemented")
}

func (p *testPage) NextInSection() Page {
	return nil
}

func (p *testPage) NextPage() Page {
	return nil
}

func (p *testPage) OutputFormats() OutputFormats {
	panic("testpage: not implemented")
}

func (p *testPage) Pages() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) RegularPages() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) RegularPagesRecursive() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) Paginate(seq any, options ...any) (*Pager, error) {
	return nil, nil
}

func (p *testPage) Paginator(options ...any) (*Pager, error) {
	return nil, nil
}

func (p *testPage) Param(key any) (any, error) {
	return resource.Param(p, nil, key)
}

func (p *testPage) Params() maps.Params {
	return p.params
}

func (p *testPage) Page() Page {
	return p
}

func (p *testPage) Parent() Page {
	panic("testpage: not implemented")
}

func (p *testPage) Ancestors() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) Path() string {
	return p.path
}

func (p *testPage) Pathc() string {
	return p.path
}

func (p *testPage) Permalink() string {
	panic("testpage: not implemented")
}

func (p *testPage) Plain(context.Context) string {
	panic("testpage: not implemented")
}

func (p *testPage) PlainWords(context.Context) []string {
	panic("testpage: not implemented")
}

func (p *testPage) Prev() Page {
	panic("testpage: not implemented")
}

func (p *testPage) PrevInSection() Page {
	return nil
}

func (p *testPage) PrevPage() Page {
	return nil
}

func (p *testPage) PublishDate() time.Time {
	return p.pubDate
}

func (p *testPage) RSSLink() template.URL {
	return ""
}

func (p *testPage) RawContent() string {
	panic("testpage: not implemented")
}

func (p *testPage) RenderShortcodes(context.Context) (template.HTML, error) {
	panic("testpage: not implemented")
}

func (p *testPage) ReadingTime(context.Context) int {
	panic("testpage: not implemented")
}

func (p *testPage) Ref(argsm map[string]any) (string, error) {
	panic("testpage: not implemented")
}

func (p *testPage) RefFrom(argsm map[string]any, source any) (string, error) {
	return "", nil
}

func (p *testPage) RelPermalink() string {
	panic("testpage: not implemented")
}

func (p *testPage) RelRef(argsm map[string]any) (string, error) {
	panic("testpage: not implemented")
}

func (p *testPage) RelRefFrom(argsm map[string]any, source any) (string, error) {
	return "", nil
}

func (p *testPage) Render(ctx context.Context, layout ...string) (template.HTML, error) {
	panic("testpage: not implemented")
}

func (p *testPage) RenderString(ctx context.Context, args ...any) (template.HTML, error) {
	panic("testpage: not implemented")
}

func (p *testPage) ResourceType() string {
	panic("testpage: not implemented")
}

func (p *testPage) Resources() resource.Resources {
	panic("testpage: not implemented")
}

func (p *testPage) Scratch() *maps.Scratch {
	panic("testpage: not implemented")
}

func (p *testPage) Store() *maps.Scratch {
	panic("testpage: not implemented")
}

func (p *testPage) RelatedKeywords(cfg related.IndexConfig) ([]related.Keyword, error) {
	v, err := p.Param(cfg.Name)
	if err != nil {
		return nil, err
	}

	return cfg.ToKeywords(v)
}

func (p *testPage) Section() string {
	return p.section
}

func (p *testPage) Sections() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) SectionsEntries() []string {
	return p.sectionEntries
}

func (p *testPage) SectionsPath() string {
	return path.Join(p.sectionEntries...)
}

func (p *testPage) Site() Site {
	return p.site
}

func (p *testPage) Sites() Sites {
	panic("testpage: not implemented")
}

func (p *testPage) Slug() string {
	return p.slug
}

func (p *testPage) String() string {
	return p.path
}

func (p *testPage) Summary(context.Context) template.HTML {
	panic("testpage: not implemented")
}

func (p *testPage) TableOfContents(context.Context) template.HTML {
	panic("testpage: not implemented")
}

func (p *testPage) Title() string {
	return p.title
}

func (p *testPage) TranslationBaseName() string {
	panic("testpage: not implemented")
}

func (p *testPage) TranslationKey() string {
	return p.path
}

func (p *testPage) Translations() Pages {
	panic("testpage: not implemented")
}

func (p *testPage) Truncated(context.Context) bool {
	panic("testpage: not implemented")
}

func (p *testPage) Type() string {
	return p.section
}

func (p *testPage) URL() string {
	return ""
}

func (p *testPage) UniqueID() string {
	panic("testpage: not implemented")
}

func (p *testPage) Weight() int {
	return p.weight
}

func (p *testPage) WordCount(context.Context) int {
	panic("testpage: not implemented")
}

func (p *testPage) GetIdentity() identity.Identity {
	panic("testpage: not implemented")
}

func createTestPages(num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		m := &testPage{
			path:           fmt.Sprintf("/x/y/z/p%d.md", i),
			weight:         5,
			fuzzyWordCount: i + 2, // magic
		}

		if i%2 == 0 {
			m.weight = 10
		}
		pages[i] = m

	}

	return pages
}
