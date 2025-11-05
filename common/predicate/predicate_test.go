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
