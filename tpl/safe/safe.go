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

	"github.com/spf13/cast"
)

// New returns a new instance of the safe-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "safe" namespace.
type Namespace struct{}

// CSS returns the string s as html/template CSS content.
func (ns *Namespace) CSS(s any) (template.CSS, error) {
	ss, err := cast.ToStringE(s)
	return template.CSS(ss), err
}

// HTML returns the string s as html/template HTML content.
func (ns *Namespace) HTML(s any) (template.HTML, error) {
	ss, err := cast.ToStringE(s)
	return template.HTML(ss), err
}

// HTMLAttr returns the string s as html/template HTMLAttr content.
func (ns *Namespace) HTMLAttr(s any) (template.HTMLAttr, error) {
	ss, err := cast.ToStringE(s)
	return template.HTMLAttr(ss), err
}

// JS returns the given string as a html/template JS content.
func (ns *Namespace) JS(s any) (template.JS, error) {
	ss, err := cast.ToStringE(s)
	return template.JS(ss), err
}

// JSStr returns the given string as a html/template JSStr content.
func (ns *Namespace) JSStr(s any) (template.JSStr, error) {
	ss, err := cast.ToStringE(s)
	return template.JSStr(ss), err
}

// URL returns the string s as html/template URL content.
func (ns *Namespace) URL(s any) (template.URL, error) {
	ss, err := cast.ToStringE(s)
	return template.URL(ss), err
}
