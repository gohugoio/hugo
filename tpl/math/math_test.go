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

package math

import (
	"math"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestBasicNSArithmetic(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	type TestCase struct {
		fn     func(inputs ...any) (any, error)
		values []any
		expect any
	}

	for _, test := range []TestCase{
		{ns.Add, []any{4, 2}, int64(6)},
		{ns.Add, []any{4, 2, 5}, int64(11)},
		{ns.Add, []any{1.0, "foo"}, false},
		{ns.Add, []any{0}, false},
		{ns.Sub, []any{4, 2}, int64(2)},
		{ns.Sub, []any{4, 2, 5}, int64(-3)},
		{ns.Sub, []any{1.0, "foo"}, false},
		{ns.Sub, []any{0}, false},
		{ns.Mul, []any{4, 2}, int64(8)},
		{ns.Mul, []any{4, 2, 5}, int64(40)},
		{ns.Mul, []any{1.0, "foo"}, false},
		{ns.Mul, []any{0}, false},
		{ns.Div, []any{4, 2}, int64(2)},
		{ns.Div, []any{4, 2, 5}, int64(0)},
		{ns.Div, []any{1.0, "foo"}, false},
		{ns.Div, []any{0}, false},
	} {

		result, err := test.fn(test.values...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestAbs(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.0, 0.0},
		{1.5, 1.5},
		{-1.5, 1.5},
		{-2, 2.0},
		{"abc", false},
	} {
		result, err := ns.Abs(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestCeil(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.1, 1.0},
		{0.5, 1.0},
		{1.1, 2.0},
		{1.5, 2.0},
		{-0.1, 0.0},
		{-0.5, 0.0},
		{-1.1, -1.0},
		{-1.5, -1.0},
		{"abc", false},
	} {

		result, err := ns.Ceil(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestFloor(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.1, 0.0},
		{0.5, 0.0},
		{1.1, 1.0},
		{1.5, 1.0},
		{-0.1, -1.0},
		{-0.5, -1.0},
		{-1.1, -2.0},
		{-1.5, -2.0},
		{"abc", false},
	} {

		result, err := ns.Floor(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLog(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		expect any
	}{
		{1, 0.0},
		{3, 1.0986},
		{0, math.Inf(-1)},
		{1.0, 0.0},
		{3.1, 1.1314},
		{"abc", false},
	} {

		result, err := ns.Log(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		if result != math.Inf(-1) {
			result = float64(int(result*10000)) / 10000
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Separate test for Log(-1) -- returns NaN
	result, err := ns.Log(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Satisfies, math.IsNaN)
}

func TestSqrt(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		expect any
	}{
		{81, 9.0},
		{0.25, 0.5},
		{0, 0.0},
		{"abc", false},
	} {

		result, err := ns.Sqrt(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		if result != math.Inf(-1) {
			result = float64(int(result*10000)) / 10000
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Separate test for Sqrt(-1) -- returns NaN
	result, err := ns.Sqrt(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Satisfies, math.IsNaN)
}

func TestMod(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		b      any
		expect any
	}{
		{3, 2, int64(1)},
		{3, 1, int64(0)},
		{3, 0, false},
		{0, 3, int64(0)},
		{3.1, 2, int64(1)},
		{3, 2.1, int64(1)},
		{3.1, 2.1, int64(1)},
		{int8(3), int8(2), int64(1)},
		{int16(3), int16(2), int64(1)},
		{int32(3), int32(2), int64(1)},
		{int64(3), int64(2), int64(1)},
		{"3", "2", int64(1)},
		{"3.1", "2", false},
		{"aaa", "0", false},
		{"3", "aaa", false},
	} {

		result, err := ns.Mod(test.a, test.b)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestModBool(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		b      any
		expect any
	}{
		{3, 3, true},
		{3, 2, false},
		{3, 1, true},
		{3, 0, nil},
		{0, 3, true},
		{3.1, 2, false},
		{3, 2.1, false},
		{3.1, 2.1, false},
		{int8(3), int8(3), true},
		{int8(3), int8(2), false},
		{int16(3), int16(3), true},
		{int16(3), int16(2), false},
		{int32(3), int32(3), true},
		{int32(3), int32(2), false},
		{int64(3), int64(3), true},
		{int64(3), int64(2), false},
		{"3", "3", true},
		{"3", "2", false},
		{"3.1", "2", nil},
		{"aaa", "0", nil},
		{"3", "aaa", nil},
	} {

		result, err := ns.ModBool(test.a, test.b)

		if test.expect == nil {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestRound(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.1, 0.0},
		{0.5, 1.0},
		{1.1, 1.0},
		{1.5, 2.0},
		{-0.1, 0.0},
		{-0.5, -1.0},
		{-1.1, -1.0},
		{-1.5, -2.0},
		{"abc", false},
	} {

		result, err := ns.Round(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestPow(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		b      any
		expect any
	}{
		{0, 0, 1.0},
		{2, 0, 1.0},
		{2, 3, 8.0},
		{-2, 3, -8.0},
		{2, -3, 0.125},
		{-2, -3, -0.125},
		{0.2, 3, 0.008},
		{2, 0.3, 1.2311},
		{0.2, 0.3, 0.617},
		{"aaa", "3", false},
		{"2", "aaa", false},
	} {

		result, err := ns.Pow(test.a, test.b)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestMax(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	type TestCase struct {
		values []any
		expect any
	}

	for _, test := range []TestCase{
		// two values
		{[]any{-1, -1}, -1.0},
		{[]any{-1, 0}, 0.0},
		{[]any{-1, 1}, 1.0},
		{[]any{0, -1}, 0.0},
		{[]any{0, 0}, 0.0},
		{[]any{0, 1}, 1.0},
		{[]any{1, -1}, 1.0},
		{[]any{1, 0}, 1.0},
		{[]any{32}, 32.0},
		{[]any{1, 1}, 1.0},
		{[]any{1.2, 1.23}, 1.23},
		{[]any{-1.2, -1.23}, -1.2},
		{[]any{0, "a"}, false},
		{[]any{"a", 0}, false},
		{[]any{"a", "b"}, false},
		// Issue #11030
		{[]any{7, []any{3, 4}}, 7.0},
		{[]any{8, []any{3, 12}, 3}, 12.0},
		{[]any{[]any{3, 5, 2}}, 5.0},
		{[]any{3, []int{3, 6}, 3}, 6.0},
		// No values.
		{[]any{}, false},

		// multi values
		{[]any{-1, -2, -3}, -1.0},
		{[]any{1, 2, 3}, 3.0},
		{[]any{"a", 2, 3}, false},
	} {
		result, err := ns.Max(test.values...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		msg := qt.Commentf("values: %v", test.values)
		c.Assert(err, qt.IsNil, msg)
		c.Assert(result, qt.Equals, test.expect, msg)
	}
}

func TestMin(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	type TestCase struct {
		values []any
		expect any
	}

	for _, test := range []TestCase{
		// two values
		{[]any{-1, -1}, -1.0},
		{[]any{-1, 0}, -1.0},
		{[]any{-1, 1}, -1.0},
		{[]any{0, -1}, -1.0},
		{[]any{0, 0}, 0.0},
		{[]any{0, 1}, 0.0},
		{[]any{1, -1}, -1.0},
		{[]any{1, 0}, 0.0},
		{[]any{1, 1}, 1.0},
		{[]any{2}, 2.0},
		{[]any{1.2, 1.23}, 1.2},
		{[]any{-1.2, -1.23}, -1.23},
		{[]any{0, "a"}, false},
		{[]any{"a", 0}, false},
		{[]any{"a", "b"}, false},
		// Issue #11030
		{[]any{1, []any{3, 4}}, 1.0},
		{[]any{8, []any{3, 2}, 3}, 2.0},
		{[]any{[]any{3, 2, 2}}, 2.0},
		{[]any{8, []int{3, 2}, 3}, 2.0},

		// No values.
		{[]any{}, false},

		// multi values
		{[]any{-1, -2, -3}, -3.0},
		{[]any{1, 2, 3}, 1.0},
		{[]any{"a", 2, 3}, false},
	} {

		result, err := ns.Min(test.values...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil, qt.Commentf("values: %v", test.values))
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestSum(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	mustSum := func(values ...any) any {
		result, err := ns.Sum(values...)
		c.Assert(err, qt.IsNil)
		return result
	}

	c.Assert(mustSum(1, 2, 3), qt.Equals, 6.0)
	c.Assert(mustSum(1, 2, 3.0), qt.Equals, 6.0)
	c.Assert(mustSum(1, 2, []any{3, 4}), qt.Equals, 10.0)
	c.Assert(mustSum(23), qt.Equals, 23.0)
	c.Assert(mustSum([]any{23}), qt.Equals, 23.0)
	c.Assert(mustSum([]any{}), qt.Equals, 0.0)

	_, err := ns.Sum()
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestProduct(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	mustProduct := func(values ...any) any {
		result, err := ns.Product(values...)
		c.Assert(err, qt.IsNil)
		return result
	}

	c.Assert(mustProduct(2, 2, 3), qt.Equals, 12.0)
	c.Assert(mustProduct(1, 2, 3.0), qt.Equals, 6.0)
	c.Assert(mustProduct(1, 2, []any{3, 4}), qt.Equals, 24.0)
	c.Assert(mustProduct(3.0), qt.Equals, 3.0)
	c.Assert(mustProduct([]string{}), qt.Equals, 0.0)

	_, err := ns.Product()
	c.Assert(err, qt.Not(qt.IsNil))
}

// Test trigonometric functions

func TestPi(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	expect := 3.1415
	result := ns.Pi()

	// we compare only 4 digits behind point if its a real float
	// otherwise we usually get different float values on the last positions
	result = float64(int(result*10000)) / 10000

	c.Assert(result, qt.Equals, expect)
}

func TestSin(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		expect any
	}{
		{0, 0.0},
		{1, 0.8414},
		{math.Pi / 2, 1.0},
		{math.Pi, 0.0},
		{-1.0, -0.8414},
		{"abc", false},
	} {

		result, err := ns.Sin(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestCos(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		expect any
	}{
		{0, 1.0},
		{1, 0.5403},
		{math.Pi / 2, 0.0},
		{math.Pi, -1.0},
		{-1.0, 0.5403},
		{"abc", false},
	} {

		result, err := ns.Cos(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTan(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		a      any
		expect any
	}{
		{0, 0.0},
		{1, 1.5574},
		// {math.Pi / 2, math.Inf(1)},
		{math.Pi, 0.0},
		{-1.0, -1.5574},
		{"abc", false},
	} {

		result, err := ns.Tan(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		if result != math.Inf(1) {
			result = float64(int(result*10000)) / 10000
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Separate test for Tan(oo) -- returns NaN
	result, err := ns.Tan(math.Inf(1))
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Satisfies, math.IsNaN)
}

// Test inverse trigonometric functions

func TestAsin(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.0, 0.0},
		{1.0, 1.5707},
		{-1.0, -1.5707},
		{0.5, 0.5235},
		{"abc", false},
	} {
		result, err := ns.Asin(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}
		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Separate test for Asin(2) -- returns NaN
	result, err := ns.Asin(2)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Satisfies, math.IsNaN)
}

func TestAcos(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{1.0, 0.0},
		{0.0, 1.5707},
		{-1.0, 3.1415},
		{0.5, 1.0471},
		{"abc", false},
	} {
		result, err := ns.Acos(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Separate test for Acos(2) -- returns NaN
	result, err := ns.Acos(2)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Satisfies, math.IsNaN)
}

func TestAtan(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.0, 0.0},
		{1, 0.7853},
		{-1.0, -0.7853},
		{math.Inf(1), 1.5707},
		{"abc", false},
	} {
		result, err := ns.Atan(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestAtan2(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{1.0, 1.0, 0.7853},
		{-1.0, 1.0, -0.7853},
		{1.0, -1.0, 2.3561},
		{-1.0, -1.0, -2.3561},
		{1, 0, 1.5707},
		{-1, 0, -1.5707},
		{0, 1, 0.0},
		{0, -1, 3.1415},
		{0.0, 0.0, 0.0},
		{"abc", "def", false},
	} {
		result, err := ns.Atan2(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

// Test angle helper functions

func TestToDegrees(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0.0, 0.0},
		{1, 57.2957},
		{math.Pi / 2, 90.0},
		{math.Pi, 180.0},
		{"abc", false},
	} {
		result, err := ns.ToDegrees(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestToRadians(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		expect any
	}{
		{0, 0.0},
		{57.29577951308232, 1.0},
		{90, 1.5707},
		{180.0, 3.1415},
		{"abc", false},
	} {
		result, err := ns.ToRadians(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		result = float64(int(result*10000)) / 10000

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}
