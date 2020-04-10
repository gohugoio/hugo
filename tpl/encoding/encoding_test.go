// Copyright 2020 The Hugo Authors. All rights reserved.
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

package encoding

import (
	"html/template"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"
)

type tstNoStringer struct{}

func TestBase64Decode(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"YWJjMTIzIT8kKiYoKSctPUB+", "abc123!?$*&()'-=@~"},
		// errors
		{t, false},
	} {

		result, err := ns.Base64Decode(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestBase64Encode(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for _, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"YWJjMTIzIT8kKiYoKSctPUB+", "WVdKak1USXpJVDhrS2lZb0tTY3RQVUIr"},
		// errors
		{t, false},
	} {

		result, err := ns.Base64Encode(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestJsonify(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for _, test := range []struct {
		opts   interface{}
		v      interface{}
		expect interface{}
	}{
		{nil, []string{"a", "b"}, template.HTML(`["a","b"]`)},
		{map[string]string{"indent": "<i>"}, []string{"a", "b"}, template.HTML("[\n<i>\"a\",\n<i>\"b\"\n]")},
		{map[string]string{"prefix": "<p>"}, []string{"a", "b"}, template.HTML("[\n<p>\"a\",\n<p>\"b\"\n<p>]")},
		{map[string]string{"prefix": "<p>", "indent": "<i>"}, []string{"a", "b"}, template.HTML("[\n<p><i>\"a\",\n<p><i>\"b\"\n<p>]")},
		{nil, tstNoStringer{}, template.HTML("{}")},
		{nil, nil, template.HTML("null")},
		// errors
		{nil, math.NaN(), false},
		{tstNoStringer{}, []string{"a", "b"}, false},
	} {
		args := []interface{}{}

		if test.opts != nil {
			args = append(args, test.opts)
		}

		args = append(args, test.v)

		result, err := ns.Jsonify(args...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}
