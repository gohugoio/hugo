// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package provides ways to identify values in Hugo. Used for dependency tracking etc.
package identity

import (
	"github.com/gobwas/glob"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
)

var _ Identity = &GlobIdentity{}

// NewGlobIdentity creates a new GlobIdentity.
func NewGlobIdentity(pattern string) *GlobIdentity {
	glob, err := hglob.GetGlob(pattern)
	if err != nil {
		panic(err)
	}

	return &GlobIdentity{
		pattern: pattern,
		glob:    glob,
	}
}

// GlobIdentity is an Identity that matches other identities based on a glob pattern.
type GlobIdentity struct {
	pattern string
	glob    glob.Glob
}

func (id *GlobIdentity) IdentifierBase() string {
	return id.pattern
}

func (id *GlobIdentity) IsProbablyDependent(other Identity) bool {
	return id.glob.Match(other.IdentifierBase())
}
