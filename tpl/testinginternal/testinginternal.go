// Copyright 2026 The Hugo Authors. All rights reserved.
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

// Package testinginternal provides template functions for internal use.
// //  These should not be used in user templates, as they are not guaranteed to be stable or even useful.
package testinginternal

import (
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/resource"
)

// New returns a new instance of the testinginternal-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	ns := &Namespace{}

	return ns
}

// Namespace provides template functions for the "testinginternal" namespace.
type Namespace struct{}

// HashString wraps the core hashing func used for e.g. calculating resource transformation keys.
func (ns *Namespace) HashString(args ...any) string {
	return hashing.HashString(args...)
}

func (ns *Namespace) NewCachedResourceGetter(args ...any) resource.ResourceGetter {
	return resource.NewCachedResourceGetter(args...)
}
