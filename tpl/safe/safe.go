// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package safe provides template functions for escaping untrusted content or
// encapsulating trusted content.
package safe

import (
	"html/template"

	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
)

// New returns a new instance of the safe-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "safe" namespace.
type Namespace struct{}

// CSS returns a given string as html/template CSS content.
func (ns *Namespace) CSS(a interface{}) (template.CSS, error) {
	s, err := cast.ToStringE(a)
	return template.CSS(s), err
}

// HTML returns a given string as html/template HTML content.
func (ns *Namespace) HTML(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	return template.HTML(s), err
}

// HTMLAttr returns a given string as html/template HTMLAttr content.
func (ns *Namespace) HTMLAttr(a interface{}) (template.HTMLAttr, error) {
	s, err := cast.ToStringE(a)
	return template.HTMLAttr(s), err
}

// JS returns the given string as a html/template JS content.
func (ns *Namespace) JS(a interface{}) (template.JS, error) {
	s, err := cast.ToStringE(a)
	return template.JS(s), err
}

// JSStr returns the given string as a html/template JSStr content.
func (ns *Namespace) JSStr(a interface{}) (template.JSStr, error) {
	s, err := cast.ToStringE(a)
	return template.JSStr(s), err
}

// URL returns a given string as html/template URL content.
func (ns *Namespace) URL(a interface{}) (template.URL, error) {
	s, err := cast.ToStringE(a)
	return template.URL(s), err
}

// SanitizeURL returns a given string as html/template URL content.
func (ns *Namespace) SanitizeURL(a interface{}) (string, error) {
	s, err := cast.ToStringE(a)
	return helpers.SanitizeURL(s), err
}
