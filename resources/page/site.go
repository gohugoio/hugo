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
	Title() string
	Sites() Sites
	Hugo() hugo.Info
	BaseURL() template.URL
	Taxonomies() interface{}
	Menus() navigation.Menus
	Params() map[string]interface{}
	Data() map[string]interface{} // TODO(bep) page Data type
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
