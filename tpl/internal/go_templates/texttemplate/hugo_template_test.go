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

package template

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/hreflect"
)

type TestStruct struct {
	S string
	M map[string]string
}

func (t TestStruct) Hello1(arg string) string {
	return arg
}

func (t TestStruct) Hello2(arg1, arg2 string) string {
	return arg1 + " " + arg2
}

type execHelper struct{}

func (e *execHelper) Init(ctx context.Context, tmpl Preparer) {
}

func (e *execHelper) GetFunc(ctx context.Context, tmpl Preparer, name string) (reflect.Value, reflect.Value, bool) {
	if name == "print" {
		return zero, zero, false
	}
	return reflect.ValueOf(func(s string) string {
		return "hello " + s
	}), zero, true
}

func (e *execHelper) GetMapValue(ctx context.Context, tmpl Preparer, m, key reflect.Value) (reflect.Value, bool) {
	key = reflect.ValueOf(strings.ToLower(key.String()))
	return m.MapIndex(key), true
}

func (e *execHelper) GetMethod(ctx context.Context, tmpl Preparer, receiver reflect.Value, name string) (reflect.Value, reflect.Value) {
	if name != "Hello1" {
		return zero, zero
	}
	m := hreflect.GetMethodByName(receiver, "Hello2")
	return m, reflect.ValueOf("v2")
}

func (e *execHelper) OnCalled(ctx context.Context, tmpl Preparer, name string, args []reflect.Value, returnValue reflect.Value) {
}

func TestTemplateExecutor(t *testing.T) {
	c := qt.New(t)

	templ, err := New("").Parse(`
{{ print "foo" }}
{{ printf "hugo" }}
Map: {{ .M.A }}
Method: {{ .Hello1 "v1" }}

`)

	c.Assert(err, qt.IsNil)

	ex := NewExecuter(&execHelper{})

	var b bytes.Buffer
	data := TestStruct{S: "sv", M: map[string]string{"a": "av"}}

	c.Assert(ex.ExecuteWithContext(context.Background(), templ, &b, data), qt.IsNil)
	got := b.String()

	c.Assert(got, qt.Contains, "foo")
	c.Assert(got, qt.Contains, "hello hugo")
	c.Assert(got, qt.Contains, "Map: av")
	c.Assert(got, qt.Contains, "Method: v2 v1")
}

// See issue 15385.
func TestSafeCallFastPaths(t *testing.T) {
	c := qt.New(t)

	vals := func(args ...any) []reflect.Value {
		var argv []reflect.Value
		for _, a := range args {
			argv = append(argv, reflect.ValueOf(a))
		}
		return argv
	}
	ctx := context.WithValue(context.Background(), "k", "v")
	errBoom := errors.New("boom")
	m := map[string]any{"a": 1}
	params := hmaps.Params{"a": 1}

	for _, test := range []struct {
		name   string
		fn     any
		args   []reflect.Value
		expect any
		err    string
		allocs float64
	}{
		{"variadic string", func(args ...any) string { return args[0].(string) }, vals("a", 1), "a", "", 2},
		{"string", func() string { return "s" }, nil, "s", "", 1},
		{"ctx variadic any error", func(c context.Context, args ...any) (any, error) { return c.Value("k"), nil }, vals(ctx, "x"), "v", "", 1},
		{"ctx variadic nil", func(c context.Context, args ...any) (any, error) { return nil, nil }, vals(ctx), nil, "", 0},
		{"ctx variadic error", func(c context.Context, args ...any) (any, error) { return nil, errBoom }, vals(ctx), nil, "boom", 0},
		{"any variadic bool", func(a any, args ...any) bool { return a == args[0] }, vals(1, 1), true, "", 1},
		{"params", func() hmaps.Params { return params }, nil, params, "", 0},
		{"any any bool error", func(a, b any) (bool, error) { return a == b, nil }, vals("a", "b"), false, "", 0},
		{"any string error", func(a any) (string, error) { return "", errBoom }, vals("a"), "", "boom", 0},
		{"ctx string variadic", func(c context.Context, s string, args ...any) (any, error) { return args[0], nil }, vals(ctx, "s", "1"), "1", "", 1},
		{"reflect value bool", func(v reflect.Value) bool { return v.Int() == 1 }, []reflect.Value{reflect.ValueOf(reflect.ValueOf(1))}, true, "", 0},
		{"map", func() map[string]any { return m }, nil, m, "", 0},
		{"variadic map error", func(args ...any) (map[string]any, error) { return m, nil }, vals(1), m, "", 1},
		{"any", func() any { return 42 }, nil, 42, "", 0},
		{"bool", func() bool { return true }, nil, true, "", 0},
		{"string variadic", func(s string, args ...any) string { return s }, vals("s", "1"), "s", "", 2},
		{"any variadic any error", func(a any, args ...any) (any, error) { return args[0], nil }, vals(1, "b"), "b", "", 1},
		{"panic", func() string { panic("bad") }, nil, nil, "bad", -1},
		{"slow path", func(a, b string) string { return a + b }, vals("a", "b"), "ab", "", -1},
	} {
		c.Run(test.name, func(c *qt.C) {
			fn := reflect.ValueOf(test.fn)
			v, err := safeCall(fn, test.args)
			if test.err != "" {
				c.Assert(err, qt.ErrorMatches, test.err)
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(v.IsValid(), qt.IsTrue)
				if test.expect == nil {
					c.Assert(v.Interface(), qt.IsNil)
				} else {
					c.Assert(v.Interface(), qt.DeepEquals, test.expect)
				}
			}
			if test.allocs >= 0 {
				allocs := testing.AllocsPerRun(100, func() { safeCall(fn, test.args) })
				c.Assert(allocs, qt.Equals, test.allocs)
			}
		})
	}
}

func BenchmarkSafeCall(b *testing.B) {
	for _, bench := range []struct {
		name string
		fn   any
		args []reflect.Value
	}{
		{"variadic string", func(args ...any) string { return "" }, []reflect.Value{reflect.ValueOf("a")}},
		{"string", func() string { return "" }, nil},
		{"ctx variadic any error", func(context.Context, ...any) (any, error) { return nil, nil }, []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf("a")}},
		{"slow path", func(a, b string) string { return a + b }, []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b")}},
	} {
		fn := reflect.ValueOf(bench.fn)
		b.Run(bench.name, func(b *testing.B) {
			for b.Loop() {
				safeCall(fn, bench.args)
			}
		})
	}
}
