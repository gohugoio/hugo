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

package compare

import (
	"math"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/spf13/cast"
)

type T struct {
	NonEmptyInterfaceNil      I
	NonEmptyInterfaceTypedNil I
}

type I interface {
	Foo() string
}

func (t *T) Foo() string {
	return "foo"
}

var testT = &T{
	NonEmptyInterfaceTypedNil: (*T)(nil),
}

type (
	tstEqerType1 string
	tstEqerType2 string
)

func (t tstEqerType2) Eq(other any) bool {
	return cast.ToString(t) == cast.ToString(other)
}

func (t tstEqerType2) String() string {
	return string(t)
}

func (t tstEqerType1) Eq(other any) bool {
	return cast.ToString(t) == cast.ToString(other)
}

func (t tstEqerType1) String() string {
	return string(t)
}

type stringType string

type tstCompareType int

const (
	tstEq tstCompareType = iota
	tstNe
	tstGt
	tstGe
	tstLt
	tstLe
)

func tstIsEq(tp tstCompareType) bool { return tp == tstEq || tp == tstGe || tp == tstLe }
func tstIsGt(tp tstCompareType) bool { return tp == tstGt || tp == tstGe }
func tstIsLt(tp tstCompareType) bool { return tp == tstLt || tp == tstLe }

func TestDefaultFunc(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	then := time.Now()
	now := time.Now()
	ns := New(time.UTC, false)

	for i, test := range []struct {
		dflt   any
		given  any
		expect any
	}{
		{true, false, false},
		{"5", 0, "5"},

		{"test1", "set", "set"},
		{"test2", "", "test2"},
		{"test3", nil, "test3"},

		{[2]int{10, 20}, [2]int{1, 2}, [2]int{1, 2}},
		{[2]int{10, 20}, [0]int{}, [2]int{10, 20}},
		{[2]int{100, 200}, nil, [2]int{100, 200}},

		{[]string{"one"}, []string{"uno"}, []string{"uno"}},
		{[]string{"two"}, []string{}, []string{"two"}},
		{[]string{"three"}, nil, []string{"three"}},

		{map[string]int{"one": 1}, map[string]int{"uno": 1}, map[string]int{"uno": 1}},
		{map[string]int{"one": 1}, map[string]int{}, map[string]int{"one": 1}},
		{map[string]int{"two": 2}, nil, map[string]int{"two": 2}},

		{10, 1, 1},
		{10, 0, 10},
		{20, nil, 20},

		{float32(10), float32(1), float32(1)},
		{float32(10), 0, float32(10)},
		{float32(20), nil, float32(20)},

		{complex(2, -2), complex(1, -1), complex(1, -1)},
		{complex(2, -2), complex(0, 0), complex(2, -2)},
		{complex(3, -3), nil, complex(3, -3)},

		{struct{ f string }{f: "one"}, struct{}{}, struct{}{}},
		{struct{ f string }{f: "two"}, nil, struct{ f string }{f: "two"}},

		{then, now, now},
		{then, time.Time{}, then},
	} {

		eq := qt.CmpEquals(hqt.DeepAllowUnexported(test.dflt))

		errMsg := qt.Commentf("[%d] %v", i, test)

		result, err := ns.Default(test.dflt, test.given)

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, eq, test.expect, errMsg)
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	n := New(time.UTC, false)

	twoEq := func(a, b any) bool {
		return n.Eq(a, b)
	}

	twoGt := func(a, b any) bool {
		return n.Gt(a, b)
	}

	twoLt := func(a, b any) bool {
		return n.Lt(a, b)
	}

	twoGe := func(a, b any) bool {
		return n.Ge(a, b)
	}

	twoLe := func(a, b any) bool {
		return n.Le(a, b)
	}

	twoNe := func(a, b any) bool {
		return n.Ne(a, b)
	}

	for _, test := range []struct {
		tstCompareType
		funcUnderTest func(a, b any) bool
	}{
		{tstGt, twoGt},
		{tstLt, twoLt},
		{tstGe, twoGe},
		{tstLe, twoLe},
		{tstEq, twoEq},
		{tstNe, twoNe},
	} {
		doTestCompare(t, test.tstCompareType, test.funcUnderTest)
	}
}

func doTestCompare(t *testing.T, tp tstCompareType, funcUnderTest func(a, b any) bool) {
	for i, test := range []struct {
		left            any
		right           any
		expectIndicator int
	}{
		{5, 8, -1},
		{8, 5, 1},
		{5, 5, 0},
		{int(5), int64(5), 0},
		{int32(5), int(5), 0},
		{int16(4), 4, 0},
		{uint8(4), 4, 0},
		{uint16(4), 4, 0},
		{uint16(4), 4, 0},
		{uint32(4), uint16(4), 0},
		{uint32(4), uint16(3), 1},
		{uint64(4), 4, 0},
		{4, uint64(4), 0},
		{uint64(math.MaxUint32), uint32(math.MaxUint32), 0},
		{uint64(math.MaxUint16), int(math.MaxUint16), 0},
		{int16(4), int(5), -1},
		{uint(15), uint64(15), 0},
		{-2, 1, -1},
		{2, -5, 1},
		{0.0, 1.23, -1},
		{1.1, 1.1, 0},
		{float32(1.0), float64(1.0), 0},
		{1.23, 0.0, 1},
		{"5", "5", 0},
		{"8", "5", 1},
		{"5", "0001", 1},
		{[]int{100, 99}, []int{1, 2, 3, 4}, -1},
		{cast.ToTime("2015-11-20"), cast.ToTime("2015-11-20"), 0},
		{cast.ToTime("2015-11-19"), cast.ToTime("2015-11-20"), -1},
		{cast.ToTime("2015-11-20"), cast.ToTime("2015-11-19"), 1},
		{"a", "a", 0},
		{"a", "b", -1},
		{"b", "a", 1},
		{"infinity", "infinity", 0},
		{"nan", "nan", 0},
		{tstEqerType1("a"), tstEqerType1("a"), 0},
		{tstEqerType1("a"), tstEqerType2("a"), 0},
		{tstEqerType2("a"), tstEqerType1("a"), 0},
		{tstEqerType2("a"), tstEqerType1("b"), -1},
		{hugo.MustParseVersion("0.32.1").Version(), hugo.MustParseVersion("0.32").Version(), 1},
		{hugo.MustParseVersion("0.35").Version(), hugo.MustParseVersion("0.32").Version(), 1},
		{hugo.MustParseVersion("0.36").Version(), hugo.MustParseVersion("0.36").Version(), 0},
		{hugo.MustParseVersion("0.32").Version(), hugo.MustParseVersion("0.36").Version(), -1},
		{hugo.MustParseVersion("0.32").Version(), "0.36", -1},
		{"0.36", hugo.MustParseVersion("0.32").Version(), 1},
		{"0.36", hugo.MustParseVersion("0.36").Version(), 0},
		{"0.37", hugo.MustParseVersion("0.37-DEV").Version(), 1},
		{"0.37-DEV", hugo.MustParseVersion("0.37").Version(), -1},
		{"0.36", hugo.MustParseVersion("0.37-DEV").Version(), -1},
		{"0.37-DEV", hugo.MustParseVersion("0.37-DEV").Version(), 0},
		// https://github.com/gohugoio/hugo/issues/5905
		{nil, nil, 0},
		{testT.NonEmptyInterfaceNil, nil, 0},
		{testT.NonEmptyInterfaceTypedNil, nil, 0},
	} {

		result := funcUnderTest(test.left, test.right)
		success := false

		if test.expectIndicator == 0 {
			if tstIsEq(tp) {
				success = result
			} else {
				success = !result
			}
		}

		if test.expectIndicator < 0 {
			success = result && (tstIsLt(tp) || tp == tstNe)
			success = success || (!result && !tstIsLt(tp))
		}

		if test.expectIndicator > 0 {
			success = result && (tstIsGt(tp) || tp == tstNe)
			success = success || (!result && (!tstIsGt(tp) || tp != tstNe))
		}

		if !success {
			t.Fatalf("[%d][%s] %v compared to %v: %t", i, path.Base(runtime.FuncForPC(reflect.ValueOf(funcUnderTest).Pointer()).Name()), test.left, test.right, result)
		}
	}
}

func TestEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{1, []any{1, 2}, true},
		{1, []any{2, 1}, true},
		{1, []any{2, 3}, false},
		{tstEqerType1("a"), []any{tstEqerType1("a"), tstEqerType1("b")}, true},
		{tstEqerType1("a"), []any{tstEqerType1("b"), tstEqerType1("a")}, true},
		{tstEqerType1("a"), []any{tstEqerType1("b"), tstEqerType1("c")}, false},
	} {

		result := ns.Eq(test.first, test.others...)

		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestNotEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{1, []any{2, 3}, true},
		{1, []any{2, 1}, false},
		{1, []any{1, 2}, false},
	} {
		result := ns.Ne(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestGreaterEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{5, []any{2, 3}, true},
		{5, []any{5, 5}, true},
		{3, []any{4, 2}, false},
		{3, []any{2, 4}, false},
	} {
		result := ns.Ge(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestGreaterThanExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{5, []any{2, 3}, true},
		{5, []any{5, 4}, false},
		{3, []any{4, 2}, false},
	} {
		result := ns.Gt(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLessEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{1, []any{2, 3}, true},
		{1, []any{1, 2}, true},
		{2, []any{1, 2}, false},
		{3, []any{2, 4}, false},
	} {
		result := ns.Le(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLessThanExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	for _, test := range []struct {
		first  any
		others []any
		expect bool
	}{
		{1, []any{2, 3}, true},
		{1, []any{1, 2}, false},
		{2, []any{1, 2}, false},
		{3, []any{2, 4}, false},
	} {
		result := ns.Lt(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestCase(t *testing.T) {
	c := qt.New(t)
	n := New(time.UTC, false)

	c.Assert(n.Eq("az", "az"), qt.Equals, true)
	c.Assert(n.Eq("az", stringType("az")), qt.Equals, true)
}

func TestStringType(t *testing.T) {
	c := qt.New(t)
	n := New(time.UTC, true)

	c.Assert(n.Lt("az", "Za"), qt.Equals, true)
	c.Assert(n.Gt("ab", "Ab"), qt.Equals, true)
}

func TestTimeUnix(t *testing.T) {
	t.Parallel()
	n := New(time.UTC, false)
	var sec int64 = 1234567890
	tv := reflect.ValueOf(time.Unix(sec, 0))
	i := 1

	res := n.toTimeUnix(tv)
	if sec != res {
		t.Errorf("[%d] timeUnix got %v but expected %v", i, res, sec)
	}

	i++
	func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("[%d] timeUnix didn't return an expected error", i)
			}
		}()
		iv := reflect.ValueOf(sec)
		n.toTimeUnix(iv)
	}(t)
}

func TestConditional(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New(time.UTC, false)

	type args struct {
		cond any
		v1   any
		v2   any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"a", args{cond: true, v1: "true", v2: "false"}, "true"},
		{"b", args{cond: false, v1: "true", v2: "false"}, "false"},
		{"c", args{cond: 1, v1: "true", v2: "false"}, "true"},
		{"d", args{cond: 0, v1: "true", v2: "false"}, "false"},
		{"e", args{cond: "foo", v1: "true", v2: "false"}, "true"},
		{"f", args{cond: "", v1: "true", v2: "false"}, "false"},
		{"g", args{cond: []int{6, 7}, v1: "true", v2: "false"}, "true"},
		{"h", args{cond: []int{}, v1: "true", v2: "false"}, "false"},
		{"i", args{cond: nil, v1: "true", v2: "false"}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Assert(ns.Conditional(tt.args.cond, tt.args.v1, tt.args.v2), qt.Equals, tt.want)
		})
	}
}

// Issue 9462
func TestComparisonArgCount(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(time.UTC, false)

	panicMsg := "missing arguments for comparison"

	c.Assert(func() { ns.Eq(1) }, qt.PanicMatches, panicMsg)
	c.Assert(func() { ns.Ge(1) }, qt.PanicMatches, panicMsg)
	c.Assert(func() { ns.Gt(1) }, qt.PanicMatches, panicMsg)
	c.Assert(func() { ns.Le(1) }, qt.PanicMatches, panicMsg)
	c.Assert(func() { ns.Lt(1) }, qt.PanicMatches, panicMsg)
	c.Assert(func() { ns.Ne(1) }, qt.PanicMatches, panicMsg)
}
