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

package page

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/modules"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/viper"

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
	return &testPage{
		params: make(map[string]interface{}),
		data:   make(map[string]interface{}),
		file:   file,
	}
}

func newTestPathSpec() *helpers.PathSpec {
	return newTestPathSpecFor(viper.New())
}

func newTestPathSpecFor(cfg config.Provider) *helpers.PathSpec {
	config.SetBaseTestDefaults(cfg)
	langs.LoadLanguageSettings(cfg, nil)
	mod, err := modules.CreateProjectModule(cfg)
	if err != nil {
		panic(err)
	}
	cfg.Set("allModules", modules.Modules{mod})
	fs := hugofs.NewMem(cfg)
	s, err := helpers.NewPathSpec(fs, cfg, nil)
	if err != nil {
		panic(err)
	}
	return s
}

type testPage struct {
	description string
	title       string
	linkTitle   string

	section string

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

	params map[string]interface{}
	data   map[string]interface{}

	file source.File
}

func (p *testPage) Aliases() []string {
	panic("not implemented")
}

func (p *testPage) AllTranslations() Pages {
	panic("not implemented")
}

func (p *testPage) AlternativeOutputFormats() OutputFormats {
	panic("not implemented")
}

func (p *testPage) Author() Author {
	return Author{}

}
func (p *testPage) Authors() AuthorList {
	return nil
}

func (p *testPage) BaseFileName() string {
	panic("not implemented")
}

func (p *testPage) BundleType() files.ContentClass {
	panic("not implemented")
}

func (p *testPage) Content() (interface{}, error) {
	panic("not implemented")
}

func (p *testPage) ContentBaseName() string {
	panic("not implemented")
}

func (p *testPage) CurrentSection() Page {
	panic("not implemented")
}

func (p *testPage) Data() interface{} {
	return p.data
}

func (p *testPage) Sitemap() config.Sitemap {
	return config.Sitemap{}
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
	panic("not implemented")
}

func (p *testPage) Draft() bool {
	panic("not implemented")
}

func (p *testPage) Eq(other interface{}) bool {
	return p == other
}

func (p *testPage) ExpiryDate() time.Time {
	return p.expiryDate
}

func (p *testPage) Ext() string {
	panic("not implemented")
}

func (p *testPage) Extension() string {
	panic("not implemented")
}

func (p *testPage) File() source.File {
	return p.file
}

func (p *testPage) FileInfo() hugofs.FileMetaInfo {
	panic("not implemented")
}

func (p *testPage) Filename() string {
	panic("not implemented")
}

func (p *testPage) FirstSection() Page {
	panic("not implemented")
}

func (p *testPage) FuzzyWordCount() int {
	return p.fuzzyWordCount
}

func (p *testPage) GetPage(ref string) (Page, error) {
	panic("not implemented")
}

func (p *testPage) GetParam(key string) interface{} {
	panic("not implemented")
}

func (p *testPage) GetTerms(taxonomy string) Pages {
	panic("not implemented")
}

func (p *testPage) GetRelatedDocsHandler() *RelatedDocsHandler {
	return relatedDocsHandler
}

func (p *testPage) GitInfo() *gitmap.GitInfo {
	return nil
}

func (p *testPage) HasMenuCurrent(menuID string, me *navigation.MenuEntry) bool {
	panic("not implemented")
}

func (p *testPage) HasShortcode(name string) bool {
	panic("not implemented")
}

func (p *testPage) Hugo() hugo.Info {
	panic("not implemented")
}

func (p *testPage) InSection(other interface{}) (bool, error) {
	panic("not implemented")
}

func (p *testPage) IsAncestor(other interface{}) (bool, error) {
	panic("not implemented")
}

func (p *testPage) IsDescendant(other interface{}) (bool, error) {
	panic("not implemented")
}

func (p *testPage) IsDraft() bool {
	return false
}

func (p *testPage) IsHome() bool {
	panic("not implemented")
}

func (p *testPage) IsMenuCurrent(menuID string, inme *navigation.MenuEntry) bool {
	panic("not implemented")
}

func (p *testPage) IsNode() bool {
	panic("not implemented")
}

func (p *testPage) IsPage() bool {
	panic("not implemented")
}

func (p *testPage) IsSection() bool {
	panic("not implemented")
}

func (p *testPage) IsTranslated() bool {
	panic("not implemented")
}

func (p *testPage) Keywords() []string {
	return nil
}

func (p *testPage) Kind() string {
	panic("not implemented")
}

func (p *testPage) Lang() string {
	panic("not implemented")
}

func (p *testPage) Language() *langs.Language {
	panic("not implemented")
}

func (p *testPage) LanguagePrefix() string {
	return ""
}

func (p *testPage) Lastmod() time.Time {
	return p.lastMod
}

func (p *testPage) Len() int {
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
	panic("not implemented")
}

func (p *testPage) MediaType() media.Type {
	panic("not implemented")
}

func (p *testPage) Menus() navigation.PageMenus {
	return navigation.PageMenus{}
}

func (p *testPage) Name() string {
	panic("not implemented")
}

func (p *testPage) Next() Page {
	panic("not implemented")
}

func (p *testPage) NextInSection() Page {
	return nil
}

func (p *testPage) NextPage() Page {
	return nil
}

func (p *testPage) OutputFormats() OutputFormats {
	panic("not implemented")
}

func (p *testPage) Pages() Pages {
	panic("not implemented")
}

func (p *testPage) RegularPages() Pages {
	panic("not implemented")
}

func (p *testPage) RegularPagesRecursive() Pages {
	panic("not implemented")
}

func (p *testPage) Paginate(seq interface{}, options ...interface{}) (*Pager, error) {
	return nil, nil
}

func (p *testPage) Paginator(options ...interface{}) (*Pager, error) {
	return nil, nil
}

func (p *testPage) Param(key interface{}) (interface{}, error) {
	return resource.Param(p, nil, key)
}

func (p *testPage) Params() maps.Params {
	return p.params
}

func (p *testPage) Page() Page {
	return p
}

func (p *testPage) Parent() Page {
	panic("not implemented")
}

func (p *testPage) Path() string {
	return p.path
}

func (p *testPage) Permalink() string {
	panic("not implemented")
}

func (p *testPage) Plain() string {
	panic("not implemented")
}

func (p *testPage) PlainWords() []string {
	panic("not implemented")
}

func (p *testPage) Prev() Page {
	panic("not implemented")
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
	panic("not implemented")
}

func (p *testPage) ReadingTime() int {
	panic("not implemented")
}

func (p *testPage) Ref(argsm map[string]interface{}) (string, error) {
	panic("not implemented")
}

func (p *testPage) RefFrom(argsm map[string]interface{}, source interface{}) (string, error) {
	return "", nil
}

func (p *testPage) RelPermalink() string {
	panic("not implemented")
}

func (p *testPage) RelRef(argsm map[string]interface{}) (string, error) {
	panic("not implemented")
}

func (p *testPage) RelRefFrom(argsm map[string]interface{}, source interface{}) (string, error) {
	return "", nil
}

func (p *testPage) Render(layout ...string) (template.HTML, error) {
	panic("not implemented")
}

func (p *testPage) RenderString(args ...interface{}) (template.HTML, error) {
	panic("not implemented")
}

func (p *testPage) ResourceType() string {
	panic("not implemented")
}

func (p *testPage) Resources() resource.Resources {
	panic("not implemented")
}

func (p *testPage) Scratch() *maps.Scratch {
	panic("not implemented")
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
	panic("not implemented")
}

func (p *testPage) SectionsEntries() []string {
	panic("not implemented")
}

func (p *testPage) SectionsPath() string {
	panic("not implemented")
}

func (p *testPage) Site() Site {
	panic("not implemented")
}

func (p *testPage) Sites() Sites {
	panic("not implemented")
}

func (p *testPage) Slug() string {
	return p.slug
}

func (p *testPage) String() string {
	return p.path
}

func (p *testPage) Summary() template.HTML {
	panic("not implemented")
}

func (p *testPage) TableOfContents() template.HTML {
	panic("not implemented")
}

func (p *testPage) Title() string {
	return p.title
}

func (p *testPage) TranslationBaseName() string {
	panic("not implemented")
}

func (p *testPage) TranslationKey() string {
	return p.path
}

func (p *testPage) Translations() Pages {
	panic("not implemented")
}

func (p *testPage) Truncated() bool {
	panic("not implemented")
}

func (p *testPage) Type() string {
	return p.section
}

func (p *testPage) URL() string {
	return ""
}

func (p *testPage) UniqueID() string {
	panic("not implemented")
}

func (p *testPage) Weight() int {
	return p.weight
}

func (p *testPage) WordCount() int {
	panic("not implemented")
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
