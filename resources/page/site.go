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
	"html/template"
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

	// Returns true if we're running in a server.
	// Deprecated: use hugo.IsServer instead
	IsServer() bool

	// Returns the server port.
	ServerPort() int

	// Returns the configured title for this Site.
	Title() string

	// Returns the configured language code for this Site.
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

	// Returns the last modification date of the content.
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

	// Author is deprecated and will be removed in a future release.
	Author() map[string]interface{}

	// Authors is deprecated and will be removed in a future release.
	Authors() AuthorList

	// Returns the social links for this site.
	Social() map[string]string

	// Deprecated: Use Config().Services.GoogleAnalytics instead.
	GoogleAnalytics() string

	// Deprecated: Use Config().Privacy.Disqus instead.
	DisqusShortname() string

	// BuildDrafts is deprecated and will be removed in a future release.
	BuildDrafts() bool

	// IsMultiLingual reports whether this site is configured with more than one language.
	IsMultiLingual() bool

	// LanguagePrefix returns the language prefix for this site.
	LanguagePrefix() string

	// Deprecated. Use site.Home.OutputFormats.Get "rss" instead.
	RSSLink() template.URL
}

// Sites represents an ordered list of sites (languages).
type Sites []Site

// First is a convenience method to get the first Site, i.e. the main language.
func (s Sites) First() Site {
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

func (s *siteWrapper) Social() map[string]string {
	return s.s.Social()
}

func (s *siteWrapper) Author() map[string]interface{} {
	return s.s.Author()
}

func (s *siteWrapper) Authors() AuthorList {
	return AuthorList{}
}

// Deprecated: Use .Site.Config.Services.GoogleAnalytics.ID instead
func (s *siteWrapper) GoogleAnalytics() string {
	return s.s.GoogleAnalytics()
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

// Deprecated: use hugo.IsServer instead
func (s *siteWrapper) IsServer() bool {
	return s.s.IsServer()
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

func (s *siteWrapper) LastChange() time.Time {
	hugo.Deprecate(".Site.LastChange", "Use .Site.Lastmod instead.", "v0.123.0")
	return s.s.Lastmod()
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

func (s *siteWrapper) IsMultiLingual() bool {
	return s.s.IsMultiLingual()
}

// Deprecated: Use .Site.Config.Services.Disqus.Shortname instead
func (s *siteWrapper) DisqusShortname() string {
	return s.s.DisqusShortname()
}

func (s *siteWrapper) LanguagePrefix() string {
	return s.s.LanguagePrefix()
}

func (s *siteWrapper) RSSLink() template.URL {
	return s.s.RSSLink()
}

// For internal use only.
func (s *siteWrapper) ForEeachIdentityByName(name string, f func(identity.Identity) bool) {
	s.s.(identity.ForEeachIdentityByNameProvider).ForEeachIdentityByName(name, f)
}

type testSite struct {
	h hugo.HugoInfo
	l *langs.Language
}

func (s testSite) Author() map[string]interface{} {
	return nil
}

func (s testSite) Authors() AuthorList {
	return AuthorList{}
}

func (s testSite) Social() map[string]string {
	return make(map[string]string)
}

func (t testSite) Hugo() hugo.HugoInfo {
	return t.h
}

func (t testSite) ServerPort() int {
	return 1313
}

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

// Deprecated: Use .Site.Config.Services.GoogleAnalytics.ID instead
func (t testSite) GoogleAnalytics() string {
	return ""
}

func (t testSite) MainSections() []string {
	return nil
}

// Deprecated: use hugo.IsServer instead
func (t testSite) IsServer() bool {
	return false
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

// Deprecated: Use .Site.Config.Services.Disqus.Shortname instead
func (testSite) DisqusShortname() string {
	return ""
}

func (s testSite) BuildDrafts() bool {
	return false
}

func (s testSite) IsMultiLingual() bool {
	return false
}

func (s testSite) Param(key any) (any, error) {
	return nil, nil
}

func (s testSite) RSSLink() template.URL {
	return ""
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
