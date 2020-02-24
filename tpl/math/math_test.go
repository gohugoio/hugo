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

	for _, test := range []struct {
		fn     func(a, b interface{}) (interface{}, error)
		a      interface{}
		b      interface{}
		expect interface{}
	}{
		{ns.Add, 4, 2, int64(6)},
		{ns.Add, 1.0, "foo", false},
		{ns.Sub, 4, 2, int64(2)},
		{ns.Sub, 1.0, "foo", false},
		{ns.Mul, 4, 2, int64(8)},
		{ns.Mul, 1.0, "foo", false},
		{ns.Div, 4, 2, int64(2)},
		{ns.Div, 1.0, "foo", false},
	} {

		result, err := test.fn(test.a, test.b)

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
		x      interface{}
		expect interface{}
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
		x      interface{}
		expect interface{}
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
		a      interface{}
		expect interface{}
	}{
		{1, float64(0)},
		{3, float64(1.0986)},
		{0, float64(math.Inf(-1))},
		{1.0, float64(0)},
		{3.1, float64(1.1314)},
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
		a      interface{}
		expect interface{}
	}{
		{81, float64(9)},
		{0.25, float64(0.5)},
		{0, float64(0)},
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
		a      interface{}
		b      interface{}
		expect interface{}
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
		a      interface{}
		b      interface{}
		expect interface{}
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
		x      interface{}
		expect interface{}
	}{
		{0.1, 0.0},
		{0.5, 1.0},
		{1.1, 1.0},
		{1.5, 2.0},
		{-0.1, -0.0},
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
