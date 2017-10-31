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
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicNSArithmetic(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := test.fn(test.a, test.b)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestDoArithmetic(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		a      interface{}
		b      interface{}
		op     rune
		expect interface{}
	}{
		{3, 2, '+', int64(5)},
		{3, 2, '-', int64(1)},
		{3, 2, '*', int64(6)},
		{3, 2, '/', int64(1)},
		{3.0, 2, '+', float64(5)},
		{3.0, 2, '-', float64(1)},
		{3.0, 2, '*', float64(6)},
		{3.0, 2, '/', float64(1.5)},
		{3, 2.0, '+', float64(5)},
		{3, 2.0, '-', float64(1)},
		{3, 2.0, '*', float64(6)},
		{3, 2.0, '/', float64(1.5)},
		{3.0, 2.0, '+', float64(5)},
		{3.0, 2.0, '-', float64(1)},
		{3.0, 2.0, '*', float64(6)},
		{3.0, 2.0, '/', float64(1.5)},
		{uint(3), uint(2), '+', uint64(5)},
		{uint(3), uint(2), '-', uint64(1)},
		{uint(3), uint(2), '*', uint64(6)},
		{uint(3), uint(2), '/', uint64(1)},
		{uint(3), 2, '+', uint64(5)},
		{uint(3), 2, '-', uint64(1)},
		{uint(3), 2, '*', uint64(6)},
		{uint(3), 2, '/', uint64(1)},
		{3, uint(2), '+', uint64(5)},
		{3, uint(2), '-', uint64(1)},
		{3, uint(2), '*', uint64(6)},
		{3, uint(2), '/', uint64(1)},
		{uint(3), -2, '+', int64(1)},
		{uint(3), -2, '-', int64(5)},
		{uint(3), -2, '*', int64(-6)},
		{uint(3), -2, '/', int64(-1)},
		{-3, uint(2), '+', int64(-1)},
		{-3, uint(2), '-', int64(-5)},
		{-3, uint(2), '*', int64(-6)},
		{-3, uint(2), '/', int64(-1)},
		{uint(3), 2.0, '+', float64(5)},
		{uint(3), 2.0, '-', float64(1)},
		{uint(3), 2.0, '*', float64(6)},
		{uint(3), 2.0, '/', float64(1.5)},
		{3.0, uint(2), '+', float64(5)},
		{3.0, uint(2), '-', float64(1)},
		{3.0, uint(2), '*', float64(6)},
		{3.0, uint(2), '/', float64(1.5)},
		{0, 0, '+', 0},
		{0, 0, '-', 0},
		{0, 0, '*', 0},
		{"foo", "bar", '+', "foobar"},
		{3, 0, '/', false},
		{3.0, 0, '/', false},
		{3, 0.0, '/', false},
		{uint(3), uint(0), '/', false},
		{3, uint(0), '/', false},
		{-3, uint(0), '/', false},
		{uint(3), 0, '/', false},
		{3.0, uint(0), '/', false},
		{uint(3), 0.0, '/', false},
		{3, "foo", '+', false},
		{3.0, "foo", '+', false},
		{uint(3), "foo", '+', false},
		{"foo", 3, '+', false},
		{"foo", "bar", '-', false},
		{3, 2, '%', false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := DoArithmetic(test.a, test.b, test.op)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestCeil(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Ceil(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestFloor(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Floor(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestLog(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Log(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		// we compare only 4 digits behind point if its a real float
		// otherwise we usually get different float values on the last positions
		if result != math.Inf(-1) {
			result = float64(int(result*10000)) / 10000
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestMod(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Mod(test.a, test.b)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestModBool(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.ModBool(test.a, test.b)

		if test.expect == nil {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestRound(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
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
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Round(test.x)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
