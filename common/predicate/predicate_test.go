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

// testIndexResolver implements IndexResolver for testing.
type testIndexResolver struct {
	values []string
}

func (r *testIndexResolver) ResolveIndex(name string) int {
	for i, v := range r.values {
		if v == name {
			return i
		}
	}
	return -1
}

func TestParseRangeOp(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		pattern   string
		wantOp    predicate.RangeOp
		wantValue string
		wantOK    bool
	}{
		{">= v2.0.0", predicate.RangeOpGte, "v2.0.0", true},
		{"<= v3.0.0", predicate.RangeOpLte, "v3.0.0", true},
		{"> v1.0.0", predicate.RangeOpGt, "v1.0.0", true},
		{"< v4.0.0", predicate.RangeOpLt, "v4.0.0", true},
		// Without space
		{">=v2.0.0", predicate.RangeOpGte, "v2.0.0", true},
		{"<=v3.0.0", predicate.RangeOpLte, "v3.0.0", true},
		{">v1.0.0", predicate.RangeOpGt, "v1.0.0", true},
		{"<v4.0.0", predicate.RangeOpLt, "v4.0.0", true},
		// Not range patterns
		{"v2.0.0", predicate.RangeOpNone, "", false},
		{"v2.*.*", predicate.RangeOpNone, "", false},
		{"! v2.0.0", predicate.RangeOpNone, "", false},
		{"**", predicate.RangeOpNone, "", false},
		{"", predicate.RangeOpNone, "", false},
	}

	for _, tt := range tests {
		c.Run(tt.pattern, func(c *qt.C) {
			op, value, ok := predicate.ParseRangeOp(tt.pattern)
			c.Assert(op, qt.Equals, tt.wantOp)
			c.Assert(value, qt.Equals, tt.wantValue)
			c.Assert(ok, qt.Equals, tt.wantOK)
		})
	}
}

func TestRangeMatcherMatchIndex(t *testing.T) {
	c := qt.New(t)

	// Test with index 2 as the target.
	// Lower index = greater value (0 > 1 > 2 > 3).
	// So "> index 2" matches indices < 2 (which are "greater" values).
	tests := []struct {
		name      string
		op        predicate.RangeOp
		targetIdx int
		testIdx   int
		want      bool
	}{
		// Greater than (> target): matches lower indices
		{"gt lower idx", predicate.RangeOpGt, 2, 1, true},
		{"gt equal", predicate.RangeOpGt, 2, 2, false},
		{"gt higher idx", predicate.RangeOpGt, 2, 3, false},
		// Greater than or equal (>= target): matches lower or equal indices
		{"gte lower idx", predicate.RangeOpGte, 2, 1, true},
		{"gte equal", predicate.RangeOpGte, 2, 2, true},
		{"gte higher idx", predicate.RangeOpGte, 2, 3, false},
		// Less than (< target): matches higher indices
		{"lt lower idx", predicate.RangeOpLt, 2, 1, false},
		{"lt equal", predicate.RangeOpLt, 2, 2, false},
		{"lt higher idx", predicate.RangeOpLt, 2, 3, true},
		// Less than or equal (<= target): matches higher or equal indices
		{"lte lower idx", predicate.RangeOpLte, 2, 1, false},
		{"lte equal", predicate.RangeOpLte, 2, 2, true},
		{"lte higher idx", predicate.RangeOpLte, 2, 3, true},
	}

	for _, tt := range tests {
		c.Run(tt.name, func(c *qt.C) {
			rm := predicate.RangeMatcher{Op: tt.op, Index: tt.targetIdx}
			c.Assert(rm.MatchIndex(tt.testIdx), qt.Equals, tt.want)
		})
	}
}

func TestParsePatterns(t *testing.T) {
	c := qt.New(t)

	// Values indexed: v1.0.0=0, v2.0.0=1, v3.0.0=2, v4.0.0=3
	// Lower index = greater value (0 > 1 > 2 > 3)
	resolver := &testIndexResolver{values: []string{"v1.0.0", "v2.0.0", "v3.0.0", "v4.0.0"}}

	c.Run("mixed patterns", func(c *qt.C) {
		patterns := []string{">= v2.0.0", "! v3.0.0", "<= v4.0.0", "v1.*.*"}
		globs, ranges, err := predicate.ParsePatterns(patterns, resolver)
		c.Assert(err, qt.IsNil)
		c.Assert(globs, qt.DeepEquals, []string{"! v3.0.0", "v1.*.*"})
		c.Assert(len(ranges), qt.Equals, 2)
		c.Assert(ranges[0].Op, qt.Equals, predicate.RangeOpGte)
		c.Assert(ranges[0].Index, qt.Equals, 1) // v2.0.0 -> index 1
		c.Assert(ranges[1].Op, qt.Equals, predicate.RangeOpLte)
		c.Assert(ranges[1].Index, qt.Equals, 3) // v4.0.0 -> index 3
	})

	c.Run("only globs", func(c *qt.C) {
		patterns := []string{"v1.*.*", "! v2.0.0"}
		globs, ranges, err := predicate.ParsePatterns(patterns, resolver)
		c.Assert(err, qt.IsNil)
		c.Assert(globs, qt.DeepEquals, []string{"v1.*.*", "! v2.0.0"})
		c.Assert(len(ranges), qt.Equals, 0)
	})

	c.Run("only ranges", func(c *qt.C) {
		patterns := []string{">= v2.0.0", "<= v3.0.0"}
		globs, ranges, err := predicate.ParsePatterns(patterns, resolver)
		c.Assert(err, qt.IsNil)
		c.Assert(len(globs), qt.Equals, 0)
		c.Assert(len(ranges), qt.Equals, 2)
	})

	c.Run("unknown value", func(c *qt.C) {
		patterns := []string{">= v5.0.0"}
		_, _, err := predicate.ParsePatterns(patterns, resolver)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.Error(), qt.Contains, "unknown value")
	})
}

func TestCombineRangeMatchers(t *testing.T) {
	c := qt.New(t)

	// Lower index = greater value (0 > 1 > 2 > 3).

	c.Run("empty ranges", func(c *qt.C) {
		filter := predicate.CombineRangeMatchers(nil)
		c.Assert(filter(0), qt.IsTrue)
		c.Assert(filter(100), qt.IsTrue)
	})

	c.Run("single range gte", func(c *qt.C) {
		// >= value at index 2 means indices <= 2 (lower indices are "greater" values)
		ranges := []predicate.RangeMatcher{{Op: predicate.RangeOpGte, Index: 2}}
		filter := predicate.CombineRangeMatchers(ranges)
		c.Assert(filter(0), qt.IsTrue)  // index 0 <= 2, so >= value at 2
		c.Assert(filter(1), qt.IsTrue)  // index 1 <= 2
		c.Assert(filter(2), qt.IsTrue)  // index 2 <= 2
		c.Assert(filter(3), qt.IsFalse) // index 3 > 2, so < value at 2
	})

	c.Run("range between", func(c *qt.C) {
		// >= value at index 3 AND <= value at index 1
		// means: indices <= 3 AND indices >= 1
		// which is indices 1, 2, 3
		ranges := []predicate.RangeMatcher{
			{Op: predicate.RangeOpGte, Index: 3},
			{Op: predicate.RangeOpLte, Index: 1},
		}
		filter := predicate.CombineRangeMatchers(ranges)
		c.Assert(filter(0), qt.IsFalse) // index 0 < 1, so > value at 1
		c.Assert(filter(1), qt.IsTrue)
		c.Assert(filter(2), qt.IsTrue)
		c.Assert(filter(3), qt.IsTrue)
		c.Assert(filter(4), qt.IsFalse) // index 4 > 3, so < value at 3
	})
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
