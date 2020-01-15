// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package templates provides template functions for working with templates.
package templates

import (
	"github.com/gohugoio/hugo/deps"
)

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "templates" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Exists returns whether the template with the given name exists.
// Note that this is the Unix-styled relative path including filename suffix,
// e.g. partials/header.html
func (ns *Namespace) Exists(name string) bool {
	_, found := ns.deps.Tmpl().Lookup(name)
	return found

}
