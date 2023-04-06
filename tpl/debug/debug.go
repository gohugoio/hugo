// Copyright 2020 The Hugo Authors. All rights reserved.
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

// Package debug provides template functions to help debugging templates.
package debug

import (
	"github.com/sanity-io/litter"

	"github.com/gohugoio/hugo/deps"
)

// New returns a new instance of the debug-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "debug" namespace.
type Namespace struct {
}

// Dump returns a object dump of val as a string.
// Note that not every value passed to Dump will print so nicely, but
// we'll improve on that.
//
// We recommend using the "go" Chroma lexer to format the output
// nicely.
//
// Also note that the output from Dump may change from Hugo version to the next,
// so don't depend on a specific output.
func (ns *Namespace) Dump(val any) string {
	return litter.Sdump(val)
}
