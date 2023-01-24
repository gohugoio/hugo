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
	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/navigation"
)

// Site represents a site. There can be multople sites in a multilingual setup.
type Site interface {
	// Returns the Language configured for this Site.
	Language() *langs.Language

	// Returns all the regular Pages in this Site.
	RegularPages() Pages

	// Returns all Pages in this Site.
	Pages() Pages

	// A shortcut to the home page.
	Home() Page

	// Returns true if we're running in a server.
	IsServer() bool

	// Returns the server port.
	ServerPort() int

	// Returns the configured title for this Site.
	Title() string

	// Returns all Sites for all languages.
	Sites() Sites

	// Returns Site currently rendering.
	Current() Site

	// Returns a struct with some information about the build.
	Hugo() hugo.Info

	// Returns the BaseURL for this Site.
	BaseURL() template.URL

	// Retuns a taxonomy map.
	Taxonomies() TaxonomyList

	// Returns the last modification date of the content.
	LastChange() time.Time

	// Returns the Menus for this site.
	Menus() navigation.Menus

	// Returns the Params configured for this site.
	Params() maps.Params

	// Returns a map of all the data inside /data.
	Data() map[string]any

	// Returns the identity of this site.
	GetIdentity() identity.Identity
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

type testSite struct {
	h hugo.Info
	l *langs.Language
}

func (t testSite) Hugo() hugo.Info {
	return t.h
}

func (t testSite) ServerPort() int {
	return 1313
}

func (testSite) LastChange() (t time.Time) {
	return
}

func (t testSite) Title() string {
	return "foo"
}

func (t testSite) Sites() Sites {
	return nil
}

func (t testSite) Current() Site {
	return t
}

func (t testSite) GetIdentity() identity.Identity {
	return identity.KeyValueIdentity{Key: "site", Value: t.l.Lang}
}

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

func (t testSite) RegularPages() Pages {
	return nil
}

func (t testSite) Menus() navigation.Menus {
	return nil
}

func (t testSite) Taxonomies() TaxonomyList {
	return nil
}

func (t testSite) BaseURL() template.URL {
	return ""
}

func (t testSite) Params() maps.Params {
	return nil
}

func (t testSite) Data() map[string]any {
	return nil
}

// NewDummyHugoSite creates a new minimal test site.
func NewDummyHugoSite(cfg config.Provider) Site {
	return testSite{
		h: hugo.NewInfo(hugo.EnvironmentProduction, nil),
		l: langs.NewLanguage("en", cfg),
	}
}
