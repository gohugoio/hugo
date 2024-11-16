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
	"time"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config/privacy"
	"github.com/gohugoio/hugo/config/services"
	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/navigation"
)

// Site represents a site. There can be multiple sites in a multilingual setup.
type Site interface {
	// Returns the Language configured for this Site.
	Language() *langs.Language

	// Returns all the languages configured for all sites.
	Languages() langs.Languages

	GetPage(ref ...string) (Page, error)

	// AllPages returns all pages for all languages.
	AllPages() Pages

	// Returns all the regular Pages in this Site.
	RegularPages() Pages

	// Returns all Pages in this Site.
	Pages() Pages

	// Returns all the top level sections.
	Sections() Pages

	// A shortcut to the home
	Home() Page

	// Returns the server port.
	ServerPort() int

	// Returns the configured title for this Site.
	Title() string

	// Deprecated: Use .Language.LanguageCode instead.
	LanguageCode() string

	// Returns the configured copyright information for this Site.
	Copyright() string

	// Returns all Sites for all languages.
	Sites() Sites

	// Returns Site currently rendering.
	Current() Site

	// Returns a struct with some information about the build.
	Hugo() hugo.HugoInfo

	// Returns the BaseURL for this Site.
	BaseURL() string

	// Returns a taxonomy map.
	Taxonomies() TaxonomyList

	// Deprecated: Use .Lastmod instead.
	LastChange() time.Time

	// Returns the last modification date of the content.
	Lastmod() time.Time

	// Returns the Menus for this site.
	Menus() navigation.Menus

	// The main sections in the site.
	MainSections() []string

	// Returns the Params configured for this site.
	Params() maps.Params

	// Param is a convenience method to do lookups in Params.
	Param(key any) (any, error)

	// Returns a map of all the data inside /data.
	Data() map[string]any

	// Returns the site config.
	Config() SiteConfig

	// Deprecated: Use taxonomies instead.
	Author() map[string]interface{}

	// Deprecated: Use taxonomies instead.
	Authors() AuthorList

	// Deprecated: Use .Site.Params instead.
	Social() map[string]string

	// BuildDrafts is deprecated and will be removed in a future release.
	BuildDrafts() bool

	// Deprecated: Use hugo.IsMultilingual instead.
	IsMultiLingual() bool

	// LanguagePrefix returns the language prefix for this site.
	LanguagePrefix() string

	maps.StoreProvider

	// For internal use only.
	// This will panic if the site is not fully initialized.
	// This is typically used to inform the user in the content adapter templates,
	// as these are executed before all the page collections etc. are ready to use.
	CheckReady()
}

// Sites represents an ordered list of sites (languages).
type Sites []Site

// Deprecated: Use .Sites.Default instead.
func (s Sites) First() Site {
	hugo.Deprecate(".Sites.First", "Use .Sites.Default instead.", "v0.127.0")
	return s.Default()
}

// Default is a convenience method to get the site corresponding to the default
// content language.
func (s Sites) Default() Site {
	if len(s) == 0 {
		return nil
	}
	return s[0]
}

// Some additional interfaces implemented by siteWrapper that's not on Site.
var _ identity.ForEeachIdentityByNameProvider = (*siteWrapper)(nil)

type siteWrapper struct {
	s Site
}

func WrapSite(s Site) Site {
	if s == nil {
		panic("Site is nil")
	}
	return &siteWrapper{s: s}
}

func (s *siteWrapper) Key() string {
	return s.s.Language().Lang
}

// Deprecated: Use .Site.Params instead.
func (s *siteWrapper) Social() map[string]string {
	return s.s.Social()
}

// Deprecated: Use taxonomies instead.
func (s *siteWrapper) Author() map[string]interface{} {
	return s.s.Author()
}

// Deprecated: Use taxonomies instead.
func (s *siteWrapper) Authors() AuthorList {
	return s.s.Authors()
}

func (s *siteWrapper) GetPage(ref ...string) (Page, error) {
	return s.s.GetPage(ref...)
}

func (s *siteWrapper) Language() *langs.Language {
	return s.s.Language()
}

func (s *siteWrapper) Languages() langs.Languages {
	return s.s.Languages()
}

func (s *siteWrapper) AllPages() Pages {
	return s.s.AllPages()
}

func (s *siteWrapper) RegularPages() Pages {
	return s.s.RegularPages()
}

func (s *siteWrapper) Pages() Pages {
	return s.s.Pages()
}

func (s *siteWrapper) Sections() Pages {
	return s.s.Sections()
}

func (s *siteWrapper) Home() Page {
	return s.s.Home()
}

func (s *siteWrapper) ServerPort() int {
	return s.s.ServerPort()
}

func (s *siteWrapper) Title() string {
	return s.s.Title()
}

func (s *siteWrapper) LanguageCode() string {
	return s.s.LanguageCode()
}

func (s *siteWrapper) Copyright() string {
	return s.s.Copyright()
}

func (s *siteWrapper) Sites() Sites {
	return s.s.Sites()
}

func (s *siteWrapper) Current() Site {
	return s.s.Current()
}

func (s *siteWrapper) Config() SiteConfig {
	return s.s.Config()
}

func (s *siteWrapper) Hugo() hugo.HugoInfo {
	return s.s.Hugo()
}

func (s *siteWrapper) BaseURL() string {
	return s.s.BaseURL()
}

func (s *siteWrapper) Taxonomies() TaxonomyList {
	return s.s.Taxonomies()
}

// Deprecated: Use .Site.Lastmod instead.
func (s *siteWrapper) LastChange() time.Time {
	return s.s.LastChange()
}

func (s *siteWrapper) Lastmod() time.Time {
	return s.s.Lastmod()
}

func (s *siteWrapper) Menus() navigation.Menus {
	return s.s.Menus()
}

func (s *siteWrapper) MainSections() []string {
	return s.s.MainSections()
}

func (s *siteWrapper) Params() maps.Params {
	return s.s.Params()
}

func (s *siteWrapper) Param(key any) (any, error) {
	return s.s.Param(key)
}

func (s *siteWrapper) Data() map[string]any {
	return s.s.Data()
}

func (s *siteWrapper) BuildDrafts() bool {
	return s.s.BuildDrafts()
}

// Deprecated: Use hugo.IsMultilingual instead.
func (s *siteWrapper) IsMultiLingual() bool {
	return s.s.IsMultiLingual()
}

func (s *siteWrapper) LanguagePrefix() string {
	return s.s.LanguagePrefix()
}

func (s *siteWrapper) Store() *maps.Scratch {
	return s.s.Store()
}

// For internal use only.
func (s *siteWrapper) ForEeachIdentityByName(name string, f func(identity.Identity) bool) {
	s.s.(identity.ForEeachIdentityByNameProvider).ForEeachIdentityByName(name, f)
}

// For internal use only.
func (s *siteWrapper) CheckReady() {
	s.s.CheckReady()
}

type testSite struct {
	h hugo.HugoInfo
	l *langs.Language
}

// Deprecated: Use taxonomies instead.
func (s testSite) Author() map[string]interface{} {
	return nil
}

// Deprecated: Use taxonomies instead.
func (s testSite) Authors() AuthorList {
	return AuthorList{}
}

// Deprecated: Use .Site.Params instead.
func (s testSite) Social() map[string]string {
	return make(map[string]string)
}

func (t testSite) Hugo() hugo.HugoInfo {
	return t.h
}

func (t testSite) ServerPort() int {
	return 1313
}

// Deprecated: Use .Site.Lastmod instead.
func (testSite) LastChange() (t time.Time) {
	return
}

func (testSite) Lastmod() (t time.Time) {
	return
}

func (t testSite) Title() string {
	return "foo"
}

func (t testSite) LanguageCode() string {
	return t.l.Lang
}

func (t testSite) Copyright() string {
	return ""
}

func (t testSite) Sites() Sites {
	return nil
}

func (t testSite) Sections() Pages {
	return nil
}

func (t testSite) GetPage(ref ...string) (Page, error) {
	return nil, nil
}

func (t testSite) Current() Site {
	return t
}

func (s testSite) LanguagePrefix() string {
	return ""
}

func (t testSite) Languages() langs.Languages {
	return nil
}

func (t testSite) MainSections() []string {
	return nil
}

func (t testSite) Language() *langs.Language {
	return t.l
}

func (t testSite) Home() Page {
	return nil
}

func (t testSite) Pages() Pages {
	return nil
}

func (t testSite) AllPages() Pages {
	return nil
}

func (t testSite) RegularPages() Pages {
	return nil
}

func (t testSite) Menus() navigation.Menus {
	return nil
}

func (t testSite) Taxonomies() TaxonomyList {
	return nil
}

func (t testSite) BaseURL() string {
	return ""
}

func (t testSite) Params() maps.Params {
	return nil
}

func (t testSite) Data() map[string]any {
	return nil
}

func (s testSite) Config() SiteConfig {
	return SiteConfig{}
}

func (s testSite) BuildDrafts() bool {
	return false
}

// Deprecated: Use hugo.IsMultilingual instead.
func (s testSite) IsMultiLingual() bool {
	return false
}

func (s testSite) Param(key any) (any, error) {
	return nil, nil
}

func (s testSite) Store() *maps.Scratch {
	return maps.NewScratch()
}

func (s testSite) CheckReady() {
}

// NewDummyHugoSite creates a new minimal test site.
func NewDummyHugoSite(conf config.AllProvider) Site {
	return testSite{
		h: hugo.NewInfo(conf, nil),
		l: &langs.Language{
			Lang: "en",
		},
	}
}

// SiteConfig holds the config in site.Config.
type SiteConfig struct {
	// This contains all privacy related settings that can be used to
	// make the YouTube template etc. GDPR compliant.
	Privacy privacy.Config

	// Services contains config for services such as Google Analytics etc.
	Services services.Config
}
