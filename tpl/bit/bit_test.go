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

package bit

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestAnd(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(0)},
		{-500, 0, int64(0)},
		{0b100, 0b10, int64(0)},
		{0b100, 0b100, int64(0b100)},
		{0b1010, 0b1100, int64(0b1000)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.And(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Expect at least 2 numbers
	_, err := ns.And(0)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.And()
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestClear(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(2000)},
		{-500, 0, int64(-500)},
		{0b100, 0b10, int64(0b100)},
		{0b100, 0b100, int64(0)},
		{0b1010, 0b1100, int64(0b10)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.Clear(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
	
	// Expect at least 2 numbers
	_, err := ns.Clear(0)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Clear()
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestExtract(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		n      any
		l      any
		s      any
		expect any
	}{
		{0, 0, 0, int64(0)},
		{0xFFFF, 0, 4, int64(0)},
		{-1, 0, 4, int64(0)},
		{0, 1, 0, int64(0)},
		{0xFFFF, 1, 4, int64(1)},
		{-1, 1, 4, int64(1)},
		{0b110100, 3, 2, int64(0b101)},
		{"abc", 3, 2, false},
		{52, "abc", 2, false},
		{52, 3, "abc", false},
	} {
		result, err := ns.Extract(test.n, test.l, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestLeadingZeros(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		n      any
		expect any
	}{
		{0, int64(64)},
		{1, int64(63)},
		{2, int64(62)},
		{-1, int64(0)},
		{-0x7FFFFFFFFFFFFFFF, int64(0)},
		{"abc", false},
	} {
		result, err := ns.LeadingZeros(test.n)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestNot(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		n      any
		expect any
	}{
		{0, int64(-1)},
		{255, int64(-256)},
		{-1, int64(0)},
		{-256, int64(255)},
		{"abc", false},
	} {
		result, err := ns.Not(test.n)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestOnesCount(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		n      any
		expect any
	}{
		{0, int64(0)},
		{0b1, int64(1)},
		{0b11, int64(2)},
		{0b100100, int64(2)},
		{-1, int64(64)},
		{"abc", false},
	} {
		result, err := ns.OnesCount(test.n)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestOr(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(2000)},
		{-500, 0, int64(-500)},
		{0b100, 0b10, int64(0b110)},
		{0b100, 0b100, int64(0b100)},
		{0b1010, 0b1100, int64(0b1110)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.Or(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Expect at least 2 numbers
	_, err := ns.Or(0)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Or()
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestShiftLeft(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(2000)},
		{-500, 0, int64(-500)},
		{0, 1, int64(0)},
		{2000, 1, int64(4000)},
		{-1000, 1, int64(-2000)},
		{0xF, 16, int64(0xF0000)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.ShiftLeft(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestShiftRight(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(2000)},
		{-500, 0, int64(-500)},
		{0, 1, int64(0)},
		{2000, 1, int64(1000)},
		{-1000, 1, int64(-500)},
		{0xF, 1, int64(0x7)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.ShiftRight(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrailingZeros(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		n      any
		expect any
	}{
		{0, int64(64)},
		{1, int64(0)},
		{2, int64(1)},
		{0b100000, int64(5)},
		{-2, int64(1)},
		{-4, int64(2)},
		{"abc", false},
	} {
		result, err := ns.TrailingZeros(test.n)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestXnor(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(-1)},
		{2000, 0, int64(-2001)},
		{-500, 0, int64(499)},
		{0b100, 0b10, int64(-0b111)},
		{0b100, 0b100, int64(-0b1)},
		{0b1010, 0b1100, int64(-0b111)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.Xnor(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Expect at least 2 numbers
	_, err := ns.Xnor(0)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Xnor()
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestXor(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		x      any
		y      any
		expect any
	}{
		{0, 0, int64(0)},
		{2000, 0, int64(2000)},
		{-500, 0, int64(-500)},
		{0b100, 0b10, int64(0b110)},
		{0b100, 0b100, int64(0)},
		{0b1010, 0b1100, int64(0b110)},
		{"abc", 7, false},
		{7, "def", false},
		{"abc", "def", false},
	} {
		result, err := ns.Xor(test.x, test.y)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Expect at least 2 numbers
	_, err := ns.Xor(0)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Xor()
	c.Assert(err, qt.Not(qt.IsNil))
}
