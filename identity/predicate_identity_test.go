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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGlobIdentity(t *testing.T) {
	c := qt.New(t)

	gid := NewGlobIdentity("/a/b/*")

	c.Assert(isNotDependent(gid, StringIdentity("/a/b/c")), qt.IsFalse)
	c.Assert(isNotDependent(gid, StringIdentity("/a/c/d")), qt.IsTrue)
	c.Assert(isNotDependent(StringIdentity("/a/b/c"), gid), qt.IsTrue)
	c.Assert(isNotDependent(StringIdentity("/a/c/d"), gid), qt.IsTrue)
}

func isNotDependent(a, b Identity) bool {
	f := NewFinder(FinderConfig{})
	r := f.Contains(a, b, -1)
	return r == 0
}

func TestPredicateIdentity(t *testing.T) {
	c := qt.New(t)

	isDependent := func(id Identity) bool {
		return id.IdentifierBase() == "foo"
	}
	isDependency := func(id Identity) bool {
		return id.IdentifierBase() == "baz"
	}

	id := NewPredicateIdentity(isDependent, isDependency)

	c.Assert(id.IsProbablyDependent(StringIdentity("foo")), qt.IsTrue)
	c.Assert(id.IsProbablyDependent(StringIdentity("bar")), qt.IsFalse)
	c.Assert(id.IsProbablyDependent(id), qt.IsFalse)
	c.Assert(id.IsProbablyDependent(NewPredicateIdentity(isDependent, nil)), qt.IsFalse)
	c.Assert(id.IsProbablyDependency(StringIdentity("baz")), qt.IsTrue)
	c.Assert(id.IsProbablyDependency(StringIdentity("foo")), qt.IsFalse)
}
