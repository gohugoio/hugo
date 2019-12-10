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

type execHelper struct {
}

func (e *execHelper) GetFunc(name string) (reflect.Value, bool) {
	if name == "print" {
		return zero, false
	}
	return reflect.ValueOf(func(s string) string {
		return "hello " + s
	}), true
}

func (e *execHelper) GetMapValue(m, key reflect.Value) (reflect.Value, bool) {
	key = reflect.ValueOf(strings.ToLower(key.String()))
	return m.MapIndex(key), true
}

func TestTemplateExecutor(t *testing.T) {
	c := qt.New(t)

	templ, err := New("").Parse(`
{{ print "foo" }}
{{ printf "hugo" }}
Map: {{ .M.A }}

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

}
