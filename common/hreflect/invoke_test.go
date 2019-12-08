// Copyright 2019 The Hugo Authors. All rights reserved.
// Some functions in this file (see comments) is based on the Go source code,
// copyright The Go Authors and  governed by a BSD-style license.
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
	"testing"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/pkg/errors"

	qt "github.com/frankban/quicktest"
)

type testStruct struct {
	Val    string
	Struct testStruct2
	err    error
}

type testStruct2 struct {
	Val2 string
	err  error
}

func (t testStruct) GetStruct() testStruct2 {
	return t.Struct
}

func (t testStruct) GetVal() string {
	return t.Val
}

func (t testStruct) OneArg(arg string) string {
	return arg
}

func (t *testStruct) GetValP() string {
	return t.Val
}

func (t testStruct) GetValError() (string, error) {
	return t.Val, t.err
}

func (t testStruct2) GetVal2() string {
	return t.Val2
}

func (t testStruct2) GetVal2Error() (string, error) {
	return t.Val2, t.err
}

func TestInvoke(t *testing.T) {
	c := qt.New(t)

	hello := testStruct{Val: "hello"}

	for _, test := range []struct {
		name   string
		dot    interface{}
		path   []string
		args   []interface{}
		expect interface{}
	}{
		{"Method", hello, []string{"GetVal"}, nil, "hello"},
		{"Method one arg", hello, []string{"OneArg"}, []interface{}{"hello"}, "hello"},
		{"Method pointer 1", &testStruct{Val: "hello"}, []string{"GetValP"}, nil, "hello"},
		{"Method pointer 2", &testStruct{Val: "hello"}, []string{"GetVal"}, nil, "hello"},
		{"Method error", testStruct{Val: "hello", err: errors.New("This failed")}, []string{"GetValError"}, nil, false},
		{"Method error nil", hello, []string{"GetValError"}, nil, "hello"},
		{"Func", func() testStruct { return hello }, []string{"GetVal"}, nil, "hello"},
		{"Func nested", func() testStruct { return testStruct{Val: "hello", Struct: testStruct2{Val2: "hello2"}} }, []string{"GetStruct", "GetVal2"}, nil, "hello2"},
		{"Field", hello, []string{"Val"}, nil, "hello"},
		{"Field pointer receiver", &testStruct{Val: "hello"}, []string{"Val"}, nil, "hello"},
		{"Method nested", testStruct{Val: "hello", Struct: testStruct2{Val2: "hello2"}}, []string{"GetStruct", "GetVal2"}, nil, "hello2"},
		{"Field nested", testStruct{Val: "hello", Struct: testStruct2{Val2: "hello2"}}, []string{"Struct", "Val2"}, nil, "hello2"},
		{"Method field nested", testStruct{Val: "hello", Struct: testStruct2{Val2: "hello2"}}, []string{"GetStruct", "Val2"}, nil, "hello2"},
		{"Method nested error", testStruct{Val: "hello", Struct: testStruct2{Val2: "hello2", err: errors.New("This failed")}}, []string{"GetStruct", "GetVal2Error"}, nil, false},
		{"Map", map[string]string{"hello": "world"}, []string{"hello"}, nil, "world"},
		{"Map not found", map[string]string{"hello": "world"}, []string{"Hugo"}, nil, nil}, // TODO1 nil type
		{"Map nested", map[string]map[string]string{"hugo": map[string]string{"does": "rock"}}, []string{"hugo", "does"}, nil, "rock"},
		{"Params", maps.Params{"hello": "world"}, []string{"Hello"}, nil, "world"},
	} {
		c.Run(test.name, func(c *qt.C) {
			got, err := Invoke(test.dot, test.path, test.args...)
			if b, ok := test.expect.(bool); ok && !b {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}
			c.Assert(got, qt.DeepEquals, test.expect)
		})
	}
}
