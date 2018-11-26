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
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	then := time.Now()
	now := time.Now()
	ns := New()

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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Default(test.dflt, test.given)

		require.NoError(t, err, errMsg)
		assert.Equal(t, result, test.expect, errMsg)
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	n := New()

	for _, test := range []struct {
		tstCompareType
		funcUnderTest func(a, b interface{}) bool
	}{
		{tstGt, n.Gt},
		{tstLt, n.Lt},
		{tstGe, n.Ge},
		{tstLe, n.Le},
		{tstEq, n.Eq},
		{tstNe, n.Ne},
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
			t.Errorf("[%d][%s] %v compared to %v: %t", i, path.Base(runtime.FuncForPC(reflect.ValueOf(funcUnderTest).Pointer()).Name()), test.left, test.right, result)
		}
	}
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
	assert := require.New(t)
	n := New()
	a, b := "a", "b"

	assert.Equal(a, n.Conditional(true, a, b))
	assert.Equal(b, n.Conditional(false, a, b))
}
