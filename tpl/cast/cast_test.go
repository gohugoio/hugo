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

package cast

import (
	"html/template"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestToInt(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for i, test := range []struct {
		v      any
		expect any
	}{
		{"1", 1},
		{template.HTML("2"), 2},
		{template.CSS("3"), 3},
		{template.HTMLAttr("4"), 4},
		{template.JS("5"), 5},
		{template.JSStr("6"), 6},
		{"a", false},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.v)

		result, err := ns.ToInt(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestToString(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		v      any
		expect any
	}{
		{1, "1"},
		{template.HTML("2"), "2"},
		{"a", "a"},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.v)

		result, err := ns.ToString(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestToFloat(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		v      any
		expect any
	}{
		{"1", 1.0},
		{template.HTML("2"), 2.0},
		{template.CSS("3"), 3.0},
		{template.HTMLAttr("4"), 4.0},
		{template.JS("-5.67"), -5.67},
		{template.JSStr("6"), 6.0},
		{"1.23", 1.23},
		{"-1.23", -1.23},
		{"0", 0.0},
		{float64(2.12), 2.12},
		{int64(123), 123.0},
		{2, 2.0},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.v)

		result, err := ns.ToFloat(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestToBool(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		v      any
		expect any
		error  any
	}{
		{"true", true, nil},
		{"false", false, nil},
		{"TRUE", true, nil},
		{"FALSE", false, nil},
		{"t", true, nil},
		{"f", false, nil},
		{"T", true, nil},
		{"F", false, nil},
		{"1", true, nil},
		{"0", false, nil},
		{1, true, nil},
		{0, false, nil},
		{true, true, nil},
		{false, false, nil},
		{nil, false, nil},

		{"cheese", false, nil},
		{"", false, nil},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.v)

		result, err := ns.ToBool(test.v)

		if b, ok := test.error.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestToTruth(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		v      any
		expect any
	}{
		{"true", true},
		{"false", true},
		{"TRUE", true},
		{"FALSE", true},
		{"t", true},
		{"f", true},
		{"T", true},
		{"F", true},
		{"1", true},
		{"0", true},
		{1, true},
		{0, false},
		{"cheese", true},
		{"", false},
		{1.67, true},
		{template.HTML("2"), true},
		{template.CSS("3"), true},
		{template.HTMLAttr("4"), true},
		{template.JS("-5.67"), true},
		{template.JSStr("6"), true},
		{t, true},
		{nil, false},
		{"null", true},
		{"undefined", true},
		{"NaN", true},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.v)

		result, err := ns.ToTruth(test.v)

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}
