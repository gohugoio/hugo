// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/htesting/hqt"
)

func TestToFuncs(t *testing.T) {
	c := qt.New(t)

	c.Assert(ToInt64(reflect.ValueOf(int(42))), qt.Equals, int64(42))
	c.Assert(ToFloat64(reflect.ValueOf(float32(3.14))), hqt.IsSameFloat64, float64(3.14))
	c.Assert(ToString(reflect.ValueOf("hello")), qt.Equals, "hello")
}

func TestConvertIfPossible(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		value    any
		typ      any
		expected any
		ok       bool
	}{
		// From uint to int.
		{
			name:  "uint64(math.MaxUint64) to int16",
			value: uint64(math.MaxUint64),
			typ:   int16(0),
			ok:    false, // overflow
		},

		{
			name:  "uint64(math.MaxUint64) to int64",
			value: uint64(math.MaxUint64),
			typ:   int64(0),
			ok:    false, // overflow
		},
		{
			name:     "uint64(math.MaxInt16) to int16",
			value:    uint64(math.MaxInt16),
			typ:      int64(0),
			ok:       true,
			expected: int64(math.MaxInt16),
		},
		// From int to int.
		{
			name:  "int64(math.MaxInt64) to int16",
			value: int64(math.MaxInt64),
			typ:   int16(0),
			ok:    false, // overflow
		},
		{
			name:     "int64(math.MaxInt16) to int",
			value:    int64(math.MaxInt16),
			typ:      int(0),
			ok:       true,
			expected: int(math.MaxInt16),
		},

		{
			name:     "int64(math.MaxInt16) to int",
			value:    int64(math.MaxInt16),
			typ:      int(0),
			ok:       true,
			expected: int(math.MaxInt16),
		},
		// From float64 to int.
		{
			name:  "float64(1.5) to int",
			value: float64(1.5),
			typ:   int(0),
			ok:    false, // loss of precision
		},
		{
			name:     "float64(1.0) to int",
			value:    float64(1.0),
			typ:      int(0),
			ok:       true,
			expected: int(1),
		},
		{
			name:  "float64(math.MaxInt16+1) to int16",
			value: float64(math.MaxInt16 + 1),
			typ:   int16(0),
			ok:    false, // overflow
		},
		{
			name:  "float64(math.MaxFloat64) to int64",
			value: float64(math.MaxFloat64),
			typ:   int64(0),
			ok:    false, // overflow
		},
		{
			name:     "float64(32767) to int16",
			value:    float64(32767),
			typ:      int16(0),
			ok:       true,
			expected: int16(32767),
		},
		// From float32 to int.
		{
			name:  "float32(1.5) to int",
			value: float32(1.5),
			typ:   int(0),
			ok:    false, // loss of precision
		},
		{
			name:     "float32(1.0) to int",
			value:    float32(1.0),
			typ:      int(0),
			ok:       true,
			expected: int(1),
		},
		{
			name:  "float32(math.MaxFloat32) to int16",
			value: float32(math.MaxFloat32),
			typ:   int16(0),
			ok:    false, // overflow
		},
		{
			name:  "float32(math.MaxFloat32) to int64",
			value: float32(math.MaxFloat32),
			typ:   int64(0),
			ok:    false, // overflow
		},
		{
			name:     "float32(math.MaxInt16) to int16",
			value:    float32(math.MaxInt16),
			typ:      int16(0),
			ok:       true,
			expected: int16(32767),
		},
		{
			name:  "float32(math.MaxInt16+1) to int16",
			value: float32(math.MaxInt16 + 1),
			typ:   int16(0),
			ok:    false, // overflow
		},
		// Int to float.
		{
			name:     "int16(32767) to float32",
			value:    int16(32767),
			typ:      float32(0),
			ok:       true,
			expected: float32(32767),
		},
		{
			name:     "int64(32767) to float32",
			value:    int64(32767),
			typ:      float32(0),
			ok:       true,
			expected: float32(32767),
		},
		{
			name:     "int64(math.MaxInt64) to float32",
			value:    int64(math.MaxInt64),
			typ:      float32(0),
			ok:       true,
			expected: float32(math.MaxInt64),
		},
		{
			name:     "int64(math.MaxInt64) to float64",
			value:    int64(math.MaxInt64),
			typ:      float64(0),
			ok:       true,
			expected: float64(math.MaxInt64),
		},
		// Int to uint.
		{
			name:     "int16(32767) to uint16",
			value:    int16(32767),
			typ:      uint16(0),
			ok:       true,
			expected: uint16(32767),
		},
		{
			name:  "int16(32767) to uint8",
			value: int16(32767),
			typ:   uint8(0),
			ok:    false,
		},
		{
			name:  "float64(3.14) to uint64",
			value: float64(3.14),
			typ:   uint64(0),
			ok:    false,
		},
		{
			name:     "float64(3.0) to uint64",
			value:    float64(3.0),
			typ:      uint64(0),
			ok:       true,
			expected: uint64(3),
		},
		// From uint to float.
		{
			name:     "uint64(math.MaxInt16) to float64",
			value:    uint64(math.MaxInt16),
			typ:      float64(0),
			ok:       true,
			expected: float64(math.MaxInt16),
		},
		// Float to float.
		{
			name:     "float64(3.14) to float32",
			value:    float64(3.14),
			typ:      float32(0),
			ok:       true,
			expected: float32(3.14),
		},
		{
			name:     "float32(3.14) to float64",
			value:    float32(3.14),
			typ:      float64(0),
			ok:       true,
			expected: float64(3.14),
		},
		{
			name:     "float64(3.14) to float64",
			value:    float64(3.14),
			typ:      float64(0),
			ok:       true,
			expected: float64(3.14),
		},
	} {

		v, ok := ConvertIfPossible(reflect.ValueOf(test.value), reflect.TypeOf(test.typ))
		c.Assert(ok, qt.Equals, test.ok, qt.Commentf("test case: %s", test.name))
		if test.ok {
			c.Assert(v.Interface(), hqt.IsSameNumber, test.expected, qt.Commentf("test case: %s", test.name))
		}
	}
}

func TestConvertIfPossibleMisc(t *testing.T) {
	c := qt.New(t)
	type s string

	var (
		i          = int32(42)
		i64        = int64(i)
		iv     any = i
		ip         = &i
		inil   any = (*int32)(nil)
		shello     = s("hello")
	)

	convertOK := func(v any, typ any) any {
		rv, ok := ConvertIfPossible(reflect.ValueOf(v), reflect.TypeOf(typ))
		c.Assert(ok, qt.IsTrue)
		return rv.Interface()
	}

	c.Assert(convertOK(shello, ""), qt.Equals, "hello")
	c.Assert(convertOK(ip, int64(0)), qt.Equals, i64)
	c.Assert(convertOK(iv, int64(0)), qt.Equals, i64)
	c.Assert(convertOK(inil, int64(0)), qt.Equals, int64(0))
}

func BenchmarkToInt64(b *testing.B) {
	v := reflect.ValueOf(int(42))
	for i := 0; i < b.N; i++ {
		ToInt64(v)
	}
}
