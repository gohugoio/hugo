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

package template

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
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

type execHelper struct {
}

func (e *execHelper) GetFunc(tmpl Preparer, name string) (reflect.Value, bool) {
	if name == "print" {
		return zero, false
	}
	return reflect.ValueOf(func(s string) string {
		return "hello " + s
	}), true
}

func (e *execHelper) GetMapValue(tmpl Preparer, m, key reflect.Value) (reflect.Value, bool) {
	key = reflect.ValueOf(strings.ToLower(key.String()))
	return m.MapIndex(key), true
}

func (e *execHelper) GetMethod(tmpl Preparer, receiver reflect.Value, name string) (method reflect.Value, firstArg reflect.Value) {
	if name != "Hello1" {
		return zero, zero
	}
	m := receiver.MethodByName("Hello2")
	return m, reflect.ValueOf("v2")
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

	c.Assert(ex.Execute(templ, &b, data), qt.IsNil)
	got := b.String()

	c.Assert(got, qt.Contains, "foo")
	c.Assert(got, qt.Contains, "hello hugo")
	c.Assert(got, qt.Contains, "Map: av")
	c.Assert(got, qt.Contains, "Method: v2 v1")

}
