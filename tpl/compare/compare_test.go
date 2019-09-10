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
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gohugoio/hugo/htesting/hqt"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hugo"
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

type tstEqerType1 string
type tstEqerType2 string

func (t tstEqerType2) Eq(other interface{}) bool {
	return cast.ToString(t) == cast.ToString(other)
}

func (t tstEqerType2) String() string {
	return string(t)
}

func (t tstEqerType1) Eq(other interface{}) bool {
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
	ns := New(false)

	for i, test := range []struct {
		dflt   interface{}
		given  interface{}
		expect interface{}
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

	n := New(false)

	twoEq := func(a, b interface{}) bool {
		return n.Eq(a, b)
	}

	twoGt := func(a, b interface{}) bool {
		return n.Gt(a, b)
	}

	twoLt := func(a, b interface{}) bool {
		return n.Lt(a, b)
	}

	twoGe := func(a, b interface{}) bool {
		return n.Ge(a, b)
	}

	twoLe := func(a, b interface{}) bool {
		return n.Le(a, b)
	}

	twoNe := func(a, b interface{}) bool {
		return n.Ne(a, b)
	}

	for _, test := range []struct {
		tstCompareType
		funcUnderTest func(a, b interface{}) bool
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

func doTestCompare(t *testing.T, tp tstCompareType, funcUnderTest func(a, b interface{}) bool) {
	for i, test := range []struct {
		left            interface{}
		right           interface{}
		expectIndicator int
	}{
		{5, 8, -1},
		{8, 5, 1},
		{5, 5, 0},
		{int(5), int64(5), 0},
		{int32(5), int(5), 0},
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

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{1, []interface{}{1, 2}, true},
		{1, []interface{}{2, 1}, true},
		{1, []interface{}{2, 3}, false},
		{tstEqerType1("a"), []interface{}{tstEqerType1("a"), tstEqerType1("b")}, true},
		{tstEqerType1("a"), []interface{}{tstEqerType1("b"), tstEqerType1("a")}, true},
		{tstEqerType1("a"), []interface{}{tstEqerType1("b"), tstEqerType1("c")}, false},
	} {

		result := ns.Eq(test.first, test.others...)

		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestNotEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{1, []interface{}{2, 3}, true},
		{1, []interface{}{2, 1}, false},
		{1, []interface{}{1, 2}, false},
	} {
		result := ns.Ne(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestGreaterEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{5, []interface{}{2, 3}, true},
		{5, []interface{}{5, 5}, true},
		{3, []interface{}{4, 2}, false},
		{3, []interface{}{2, 4}, false},
	} {
		result := ns.Ge(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestGreaterThanExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{5, []interface{}{2, 3}, true},
		{5, []interface{}{5, 4}, false},
		{3, []interface{}{4, 2}, false},
	} {
		result := ns.Gt(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLessEqualExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{1, []interface{}{2, 3}, true},
		{1, []interface{}{1, 2}, true},
		{2, []interface{}{1, 2}, false},
		{3, []interface{}{2, 4}, false},
	} {
		result := ns.Le(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLessThanExtend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New(false)

	for _, test := range []struct {
		first  interface{}
		others []interface{}
		expect bool
	}{
		{1, []interface{}{2, 3}, true},
		{1, []interface{}{1, 2}, false},
		{2, []interface{}{1, 2}, false},
		{3, []interface{}{2, 4}, false},
	} {
		result := ns.Lt(test.first, test.others...)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestCase(t *testing.T) {
	c := qt.New(t)
	n := New(false)

	c.Assert(n.Eq("az", "az"), qt.Equals, true)
	c.Assert(n.Eq("az", stringType("az")), qt.Equals, true)

}

func TestStringType(t *testing.T) {
	c := qt.New(t)
	n := New(true)

	c.Assert(n.Lt("az", "Za"), qt.Equals, true)
	c.Assert(n.Gt("ab", "Ab"), qt.Equals, true)
}

func TestTimeUnix(t *testing.T) {
	t.Parallel()
	var sec int64 = 1234567890
	tv := reflect.ValueOf(time.Unix(sec, 0))
	i := 1

	res := toTimeUnix(tv)
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
		toTimeUnix(iv)
	}(t)
}

func TestConditional(t *testing.T) {
	c := qt.New(t)
	n := New(false)
	a, b := "a", "b"

	c.Assert(n.Conditional(true, a, b), qt.Equals, a)
	c.Assert(n.Conditional(false, a, b), qt.Equals, b)
}
