// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hreflect

import (
	"math"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCompareNumbers(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		a, b   any
		expect int
	}{
		{1, 1, 0},
		{1, 2, -1},
		{int8(-1), int64(-2), 1},
		{uint8(1), uint64(1), 0},
		{uint64(math.MaxUint64), uint8(1), 1},
		{-1, uint64(math.MaxUint64), -1},
		{uint64(math.MaxUint64), -1, 1},
		{int64(math.MaxInt64), uint64(math.MaxInt64), 0},
		{int64(math.MaxInt64), uint64(1 << 63), -1},
		{1, 1.0, 0},
		{1.0, 1, 0},
		{1, 1.5, -1},
		{2, 1.5, 1},
		{-1, -1.5, 1},
		{-2, -1.5, -1},
		{1.5, 1, 1},
		{-1.5, -1, -1},
		{int64(1<<53 + 1), float64(1 << 53), 1},
		{float64(1 << 53), int64(1<<53 + 1), -1},
		{int64(1 << 53), float64(1 << 53), 0},
		{int64(math.MaxInt64), float64(1 << 63), -1},
		{int64(math.MinInt64), float64(-(1 << 63)), 0},
		{int64(math.MinInt64), -1e300, 1},
		{uint64(math.MaxUint64), float64(1 << 64), -1},
		{uint64(0), -0.5, 1},
		{uint64(0), math.Copysign(0, -1), 0},
		{uint64(1<<53 + 1), float64(1 << 53), 1},
		{float32(1.5), 1.5, 0},
		{float32(1.1), 1.1, 1},
		{math.Inf(1), int64(math.MaxInt64), 1},
		{math.Inf(-1), uint64(0), -1},
	} {
		got, ok := CompareNumbers(reflect.ValueOf(test.a), reflect.ValueOf(test.b))
		c.Assert(ok, qt.IsTrue, qt.Commentf("%v vs %v", test.a, test.b))
		c.Assert(got, qt.Equals, test.expect, qt.Commentf("%v vs %v", test.a, test.b))
	}

	for _, test := range [][2]any{
		{math.NaN(), 1},
		{1, math.NaN()},
		{math.NaN(), math.NaN()},
		{"1", 1},
		{1, "1"},
		{true, 1},
		{nil, 1},
	} {
		_, ok := CompareNumbers(reflect.ValueOf(test[0]), reflect.ValueOf(test[1]))
		c.Assert(ok, qt.IsFalse, qt.Commentf("%v vs %v", test[0], test[1]))
	}
}
