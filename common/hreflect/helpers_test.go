// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"context"
	"math"
	"reflect"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

type zeroStruct struct {
	zero bool
}

func (z zeroStruct) IsZero() bool {
	return z.zero
}

func TestIsTruthful(t *testing.T) {
	c := qt.New(t)

	var nilpointerZero *zeroStruct

	c.Assert(IsTruthful(true), qt.Equals, true)
	c.Assert(IsTruthful(false), qt.Equals, false)
	c.Assert(IsTruthful(time.Now()), qt.Equals, true)
	c.Assert(IsTruthful(time.Time{}), qt.Equals, false)
	c.Assert(IsTruthful(&zeroStruct{zero: false}), qt.Equals, true)
	c.Assert(IsTruthful(&zeroStruct{zero: true}), qt.Equals, false)
	c.Assert(IsTruthful(zeroStruct{zero: false}), qt.Equals, true)
	c.Assert(IsTruthful(zeroStruct{zero: true}), qt.Equals, false)
	c.Assert(IsTruthful(nil), qt.Equals, false)
	c.Assert(IsTruthful(nilpointerZero), qt.Equals, false)
}

func TestGetMethodByName(t *testing.T) {
	c := qt.New(t)
	v := reflect.ValueOf(&testStruct{})
	tp := v.Type()

	c.Assert(GetMethodIndexByName(tp, "Method1"), qt.Equals, 0)
	c.Assert(GetMethodIndexByName(tp, "Method3"), qt.Equals, 2)
	c.Assert(GetMethodIndexByName(tp, "Foo"), qt.Equals, -1)
}

func TestIsContextType(t *testing.T) {
	c := qt.New(t)
	type k string
	ctx := context.Background()
	valueCtx := context.WithValue(ctx, k("key"), 32)
	c.Assert(IsContextType(reflect.TypeOf(ctx)), qt.IsTrue)
	c.Assert(IsContextType(reflect.TypeOf(valueCtx)), qt.IsTrue)
}

func TestToSliceAny(t *testing.T) {
	c := qt.New(t)

	checkOK := func(in any, expected []any) {
		out, ok := ToSliceAny(in)
		c.Assert(ok, qt.Equals, true)
		c.Assert(out, qt.DeepEquals, expected)
	}

	checkOK([]any{1, 2, 3}, []any{1, 2, 3})
	checkOK([]int{1, 2, 3}, []any{1, 2, 3})
}

func BenchmarkIsContextType(b *testing.B) {
	type k string
	b.Run("value", func(b *testing.B) {
		ctx := context.Background()
		ctxs := make([]reflect.Type, b.N)
		for i := 0; i < b.N; i++ {
			ctxs[i] = reflect.TypeOf(context.WithValue(ctx, k("key"), i))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !IsContextType(ctxs[i]) {
				b.Fatal("not context")
			}
		}
	})

	b.Run("background", func(b *testing.B) {
		var ctxt reflect.Type = reflect.TypeOf(context.Background())
		for i := 0; i < b.N; i++ {
			if !IsContextType(ctxt) {
				b.Fatal("not context")
			}
		}
	})
}

func BenchmarkIsTruthFulValue(b *testing.B) {
	var (
		stringHugo  = reflect.ValueOf("Hugo")
		stringEmpty = reflect.ValueOf("")
		zero        = reflect.ValueOf(time.Time{})
		timeNow     = reflect.ValueOf(time.Now())
		boolTrue    = reflect.ValueOf(true)
		boolFalse   = reflect.ValueOf(false)
		nilPointer  = reflect.ValueOf((*zeroStruct)(nil))
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsTruthfulValue(stringHugo)
		IsTruthfulValue(stringEmpty)
		IsTruthfulValue(zero)
		IsTruthfulValue(timeNow)
		IsTruthfulValue(boolTrue)
		IsTruthfulValue(boolFalse)
		IsTruthfulValue(nilPointer)
	}
}

type testStruct struct{}

func (t *testStruct) Method1() string {
	return "Hugo"
}

func (t *testStruct) Method2() string {
	return "Hugo"
}

func (t *testStruct) Method3() string {
	return "Hugo"
}

func (t *testStruct) Method4() string {
	return "Hugo"
}

func (t *testStruct) Method5() string {
	return "Hugo"
}

func BenchmarkGetMethodByName(b *testing.B) {
	v := reflect.ValueOf(&testStruct{})
	methods := []string{"Method1", "Method2", "Method3", "Method4", "Method5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, method := range methods {
			_ = GetMethodByName(v, method)
		}
	}
}

func BenchmarkGetMethodByNamePara(b *testing.B) {
	v := reflect.ValueOf(&testStruct{})
	methods := []string{"Method1", "Method2", "Method3", "Method4", "Method5"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, method := range methods {
				_ = GetMethodByName(v, method)
			}
		}
	})
}

func TestCastIfPossible(t *testing.T) {
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
		// From float to int.
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
			name:  "float64(math.MaxFloat64) to int16",
			value: float64(math.MaxFloat64),
			typ:   int16(0),
			ok:    false, // overflow
		},
		{
			name:     "float64(32767) to int16",
			value:    float64(32767),
			typ:      int16(0),
			ok:       true,
			expected: int16(32767),
		},
	} {

		v, ok := ConvertIfPossible(reflect.ValueOf(test.value), reflect.TypeOf(test.typ))
		c.Assert(ok, qt.Equals, test.ok, qt.Commentf("test case: %s", test.name))
		if test.ok {
			c.Assert(v.Interface(), qt.Equals, test.expected, qt.Commentf("test case: %s", test.name))
		}
	}
}
