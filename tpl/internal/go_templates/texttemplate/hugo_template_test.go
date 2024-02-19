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
	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
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
