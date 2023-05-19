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
