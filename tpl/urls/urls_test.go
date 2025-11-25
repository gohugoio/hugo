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

package urls

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/htesting/hqt"

	qt "github.com/frankban/quicktest"
)

func newNs() *Namespace {
	return New(testconfig.GetTestDeps(nil, nil))
}

type tstNoStringer struct{}

func TestParse(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		rawurl any
		expect any
	}{
		{
			"http://www.google.com",
			&url.URL{
				Scheme: "http",
				Host:   "www.google.com",
			},
		},
		{
			"http://j@ne:password@google.com",
			&url.URL{
				Scheme: "http",
				User:   url.UserPassword("j@ne", "password"),
				Host:   "google.com",
			},
		},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Parse(test.rawurl)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result,
			qt.CmpEquals(hqt.DeepAllowUnexported(&url.URL{}, url.Userinfo{})), test.expect)
	}
}

func TestJoinPath(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		elements any
		expect   any
	}{
		{"", `/`},
		{"a", `a`},
		{"/a/b", `/a/b`},
		{"./../a/b", `a/b`},
		{[]any{""}, `/`},
		{[]any{"a"}, `a`},
		{[]any{"/a", "b"}, `/a/b`},
		{[]any{".", "..", "/a", "b"}, `a/b`},
		{[]any{"https://example.org", "a"}, `https://example.org/a`},
		{[]any{nil}, `/`},
		// errors
		{tstNoStringer{}, false},
		{[]any{tstNoStringer{}}, false},
	} {

		result, err := ns.JoinPath(test.elements)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestPathEscape(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	tests := []struct {
		name     string
		input    any
		want     string
		wantErr  bool
		errCheck string
	}{
		{"string", "A/b/c?d=é&f=g+h", "A%2Fb%2Fc%3Fd=%C3%A9&f=g+h", false, ""},
		{"empty string", "", "", false, ""},
		{"integer", 6, "6", false, ""},
		{"float", 7.42, "7.42", false, ""},
		{"nil", nil, "", false, ""},
		{"slice", []int{}, "", true, "unable to cast"},
		{"map", map[string]string{}, "", true, "unable to cast"},
		{"struct", tstNoStringer{}, "", true, "unable to cast"},
	}

	for _, tt := range tests {
		c.Run(tt.name, func(c *qt.C) {
			got, err := ns.PathEscape(tt.input)
			if tt.wantErr {
				c.Assert(err, qt.IsNotNil, qt.Commentf("PathEscape(%v) should have failed", tt.input))
				if tt.errCheck != "" {
					c.Assert(err, qt.ErrorMatches, ".*"+regexp.QuoteMeta(tt.errCheck)+".*")
				}
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(got, qt.Equals, tt.want)
			}
		})
	}
}

func TestPathUnescape(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	tests := []struct {
		name     string
		input    any
		want     string
		wantErr  bool
		errCheck string
	}{
		{"string", "A%2Fb%2Fc%3Fd=%C3%A9&f=g+h", "A/b/c?d=é&f=g+h", false, ""},
		{"empty string", "", "", false, ""},
		{"integer", 6, "6", false, ""},
		{"float", 7.42, "7.42", false, ""},
		{"nil", nil, "", false, ""},
		{"slice", []int{}, "", true, "unable to cast"},
		{"map", map[string]string{}, "", true, "unable to cast"},
		{"struct", tstNoStringer{}, "", true, "unable to cast"},
		{"malformed hex", "bad%g0escape", "", true, "invalid URL escape"},
		{"incomplete hex", "trailing%", "", true, "invalid URL escape"},
		{"single hex digit", "trail%1", "", true, "invalid URL escape"},
	}

	for _, tt := range tests {
		c.Run(tt.name, func(c *qt.C) {
			got, err := ns.PathUnescape(tt.input)
			if tt.wantErr {
				c.Assert(err, qt.Not(qt.IsNil), qt.Commentf("PathUnescape(%v) should have failed", tt.input))
				if tt.errCheck != "" {
					c.Assert(err, qt.ErrorMatches, ".*"+regexp.QuoteMeta(tt.errCheck)+".*")
				}
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(got, qt.Equals, tt.want)
			}
		})
	}
}
