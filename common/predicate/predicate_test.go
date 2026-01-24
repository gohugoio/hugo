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

package predicate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/predicate"
)

func TestPredicate(t *testing.T) {
	c := qt.New(t)

	n := func() predicate.PR[int] {
		var pr predicate.PR[int]
		return pr
	}

	var pr predicate.PR[int]
	p := pr.BoolFunc()
	c.Assert(p(1), qt.IsFalse)

	pr = n().Or(intP1).Or(intP2)
	p = pr.BoolFunc()
	c.Assert(p(1), qt.IsTrue)  // true || false
	c.Assert(p(2), qt.IsTrue)  // false || true
	c.Assert(p(3), qt.IsFalse) // false || false

	pr = pr.And(intP3)
	p = pr.BoolFunc()
	c.Assert(p(2), qt.IsFalse)  // true || true && false
	c.Assert(pr(10), qt.IsTrue) // true || true && true

	pr = pr.And(intP4)
	p = pr.BoolFunc()
	c.Assert(p(10), qt.IsTrue)  // true || true && true && true
	c.Assert(p(2), qt.IsFalse)  // true || true && false && false
	c.Assert(p(1), qt.IsFalse)  // true || false && false && false
	c.Assert(p(3), qt.IsFalse)  // false || false && false && false
	c.Assert(p(4), qt.IsFalse)  // false || false && false && false
	c.Assert(p(42), qt.IsFalse) // false || false && false && false

	pr = n().And(intP1).And(intP2).And(intP3).And(intP4)
	p = pr.BoolFunc()
	c.Assert(p(1), qt.IsFalse)
	c.Assert(p(2), qt.IsFalse)
	c.Assert(p(10), qt.IsTrue)

	pr = n().And(intP1).And(intP2).And(intP3).And(intP4)
	p = pr.BoolFunc()
	c.Assert(p(1), qt.IsFalse)
	c.Assert(p(2), qt.IsFalse)
	c.Assert(p(10), qt.IsTrue)

	pr = n().Or(intP1).Or(intP2).Or(intP3)
	p = pr.BoolFunc()
	c.Assert(p(1), qt.IsTrue)
	c.Assert(p(10), qt.IsTrue)
	c.Assert(p(4), qt.IsFalse)
}

func TestFilter(t *testing.T) {
	c := qt.New(t)

	var p predicate.PR[int]
	p = p.Or(intP1).Or(intP2)

	ints := []int{1, 2, 3, 4, 1, 6, 7, 8, 2}

	c.Assert(p.Filter(ints), qt.DeepEquals, []int{1, 2, 1, 2})
	c.Assert(ints, qt.DeepEquals, []int{1, 2, 1, 2, 1, 6, 7, 8, 2})
}

func TestFilterCopy(t *testing.T) {
	c := qt.New(t)

	var p predicate.PR[int]
	p = p.Or(intP1).Or(intP2)

	ints := []int{1, 2, 3, 4, 1, 6, 7, 8, 2}

	c.Assert(p.FilterCopy(ints), qt.DeepEquals, []int{1, 2, 1, 2})
	c.Assert(ints, qt.DeepEquals, []int{1, 2, 3, 4, 1, 6, 7, 8, 2})
}

var intP1 = func(i int) predicate.Match {
	if i == 10 {
		return predicate.True
	}
	return predicate.BoolMatch(i == 1)
}

var intP2 = func(i int) predicate.Match {
	if i == 10 {
		return predicate.True
	}
	return predicate.BoolMatch(i == 2)
}

var intP3 = func(i int) predicate.Match {
	if i == 10 {
		return predicate.True
	}
	return predicate.BoolMatch(i == 3)
}

var intP4 = func(i int) predicate.Match {
	if i == 10 {
		return predicate.True
	}
	return predicate.BoolMatch(i == 4)
}

func TestNewStringPredicateFromGlobs(t *testing.T) {
	c := qt.New(t)

	getGlob := func(pattern string) (glob.Glob, error) {
		return glob.Compile(pattern)
	}

	n := func(patterns ...string) predicate.P[string] {
		p, err := predicate.NewStringPredicateFromGlobs(patterns, getGlob)
		c.Assert(err, qt.IsNil)
		return p
	}

	m := n("a", "! ab*", "abc")
	c.Assert(m("a"), qt.IsTrue)
	c.Assert(m("ab"), qt.IsFalse)
	c.Assert(m("abc"), qt.IsFalse)

	m = n()
	c.Assert(m("anything"), qt.IsFalse)
}

func TestNewIndexStringPredicateFromGlobsAndRanges(t *testing.T) {
	c := qt.New(t)

	// Simulate versions: v4.0.0=0, v3.0.0=1, v2.0.0=2, v1.0.0=3
	// Lower index = greater value.
	versions := []string{"v4.0.0", "v3.0.0", "v2.0.0", "v1.0.0"}
	getIndex := func(s string) int {
		for i, v := range versions {
			if v == s {
				return i
			}
		}
		return -1
	}
	getGlob := func(pattern string) (glob.Glob, error) {
		return glob.Compile(pattern)
	}

	n := func(patterns ...string) predicate.P[predicate.IndexString] {
		p, err := predicate.NewIndexStringPredicateFromGlobsAndRanges(patterns, getIndex, getGlob)
		c.Assert(err, qt.IsNil)
		return p
	}

	is := func(i int) predicate.IndexString {
		return predicate.IndexString{Index: i, String: versions[i]}
	}

	// Test >= v2.0.0 (index 2): should match indices <= 2
	m := n(">= v2.0.0")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsTrue)  // v3.0.0
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0

	// Test > v2.0.0 (index 2): should match indices < 2
	m = n("> v2.0.0")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsTrue)  // v3.0.0
	c.Assert(m(is(2)), qt.IsFalse) // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0

	// Test < v3.0.0 (index 1): should match indices > 1
	m = n("< v3.0.0")
	c.Assert(m(is(0)), qt.IsFalse) // v4.0.0
	c.Assert(m(is(1)), qt.IsFalse) // v3.0.0
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsTrue)  // v1.0.0

	// Test range: >= v2.0.0 AND <= v3.0.0
	m = n(">= v2.0.0", "<= v3.0.0")
	c.Assert(m(is(0)), qt.IsFalse) // v4.0.0 - too high
	c.Assert(m(is(1)), qt.IsTrue)  // v3.0.0
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0 - too low

	// Test glob pattern
	m = n("v2.*.*")
	c.Assert(m(is(0)), qt.IsFalse) // v4.0.0
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0

	// Test glob with negation
	m = n("v*.*.*", "! v3.*.*")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsFalse) // v3.0.0 - negated
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0

	// Test range with negation: >= v2.0.0 but not v3.0.0
	m = n(">= v2.0.0", "! v3.0.0")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsFalse) // v3.0.0 - negated
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0 - out of range

	// Test unknown value in range returns no match
	m = n(">= v99.0.0")
	c.Assert(m(is(0)), qt.IsFalse)
	c.Assert(m(is(3)), qt.IsFalse)

	// Test == v2.0.0: should only match v2.0.0
	m = n("== v2.0.0")
	c.Assert(m(is(0)), qt.IsFalse) // v4.0.0
	c.Assert(m(is(1)), qt.IsFalse) // v3.0.0
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0

	// Test != v2.0.0: should match everything except v2.0.0
	m = n("!= v2.0.0")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsTrue)  // v3.0.0
	c.Assert(m(is(2)), qt.IsFalse) // v2.0.0
	c.Assert(m(is(3)), qt.IsTrue)  // v1.0.0

	// Test != with range: >= v2.0.0 AND != v3.0.0
	m = n(">= v2.0.0", "!= v3.0.0")
	c.Assert(m(is(0)), qt.IsTrue)  // v4.0.0
	c.Assert(m(is(1)), qt.IsFalse) // v3.0.0 - excluded by !=
	c.Assert(m(is(2)), qt.IsTrue)  // v2.0.0
	c.Assert(m(is(3)), qt.IsFalse) // v1.0.0 - out of range
}

func BenchmarkPredicate(b *testing.B) {
	b.Run("and or no match", func(b *testing.B) {
		var p predicate.PR[int] = intP1
		p = p.And(intP2).Or(intP3)
		for b.Loop() {
			_ = p(3).OK()
		}
	})

	b.Run("and and no match", func(b *testing.B) {
		var p predicate.PR[int] = intP1
		p = p.And(intP2)
		for b.Loop() {
			_ = p(3).OK()
		}
	})

	b.Run("and and match", func(b *testing.B) {
		var p predicate.PR[int] = intP1
		p = p.And(intP2)
		for b.Loop() {
			_ = p(10).OK()
		}
	})

	b.Run("or or match", func(b *testing.B) {
		var p predicate.PR[int] = intP1
		p = p.Or(intP2).Or(intP3)
		for b.Loop() {
			_ = p(2).OK()
		}
	})
}
