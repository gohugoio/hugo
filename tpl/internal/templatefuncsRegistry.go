// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package internal

import (
	"github.com/spf13/hugo/deps"
)

var TemplateFuncsNamespaceRegistry []func(d *deps.Deps) *TemplateFuncsNamespace

func AddTemplateFuncsNamespace(ns func(d *deps.Deps) *TemplateFuncsNamespace) {
	TemplateFuncsNamespaceRegistry = append(TemplateFuncsNamespaceRegistry, ns)
}

type TemplateFuncsNamespace struct {
	// The namespace name, "strings", "lang", etc.
	Name string

	// This is the method receiver.
	Context interface{}

	// Any template funcs aliases. This is mainly motivated by keeping
	// backwards compability, but some new template funcs may also make
	// sense to give short and snappy aliases.
	// Note that these aliases are global and will be merged, so the last
	// key will win.
	Aliases map[string]interface{}

	// A slice of input/expected examples.
	// We keep it a the namespace level for now, but may find a way to keep track
	// of the single template func, for documentation purposes.
	// Some of these, hopefully just a few, may depend on some test data to run.
	Examples [][2]string
}
