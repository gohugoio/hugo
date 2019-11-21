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

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/navigation"
)

// Site represents a site in the build. This is currently a very narrow interface,
// but the actual implementation will be richer, see hugolib.SiteInfo.
type Site interface {
	Language() *langs.Language
	RegularPages() Pages
	Pages() Pages
	IsServer() bool
	ServerPort() int
	Title() string
	Sites() Sites
	Hugo() hugo.Info
	BaseURL() template.URL
	Taxonomies() interface{}
	LastChange() time.Time
	Menus() navigation.Menus
	Params() maps.Params
	Data() map[string]interface{}
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

func (t testSite) IsServer() bool {
	return false
}

func (t testSite) Language() *langs.Language {
	return t.l
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

func (t testSite) Taxonomies() interface{} {
	return nil
}

func (t testSite) BaseURL() template.URL {
	return ""
}

func (t testSite) Params() maps.Params {
	return nil
}

func (t testSite) Data() map[string]interface{} {
	return nil
}

// NewDummyHugoSite creates a new minimal test site.
func NewDummyHugoSite(cfg config.Provider) Site {
	return testSite{
		h: hugo.NewInfo(hugo.EnvironmentProduction),
		l: langs.NewLanguage("en", cfg),
	}
}
