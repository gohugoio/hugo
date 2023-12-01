// Copyright 2021 The Hugo Authors. All rights reserved.
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

package paths

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMakePermalink(t *testing.T) {
	type test struct {
		host, link, output string
	}

	data := []test{
		{"http://abc.com/foo", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com/foo/", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com", "post/bar", "http://abc.com/post/bar"},
		{"http://abc.com", "bar", "http://abc.com/bar"},
		{"http://abc.com/foo/bar", "post/bar", "http://abc.com/foo/bar/post/bar"},
		{"http://abc.com/foo/bar", "post/bar/", "http://abc.com/foo/bar/post/bar/"},
		{"http://abc.com/foo", "post/bar?a=b#c", "http://abc.com/foo/post/bar?a=b#c"},
	}

	for i, d := range data {
		output := MakePermalink(d.host, d.link).String()
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}
}

func TestAddContextRoot(t *testing.T) {
	tests := []struct {
		baseURL  string
		url      string
		expected string
	}{
		{"http://example.com/sub/", "/foo", "/sub/foo"},
		{"http://example.com/sub/", "/foo/index.html", "/sub/foo/index.html"},
		{"http://example.com/sub1/sub2", "/foo", "/sub1/sub2/foo"},
		{"http://example.com", "/foo", "/foo"},
		// cannot guess that the context root is already added int the example below
		{"http://example.com/sub/", "/sub/foo", "/sub/sub/foo"},
		{"http://example.com/тря", "/трям/", "/тря/трям/"},
		{"http://example.com", "/", "/"},
		{"http://example.com/bar", "//", "/bar/"},
	}

	for _, test := range tests {
		output := AddContextRoot(test.baseURL, test.url)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestPretty(t *testing.T) {
	c := qt.New(t)
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name.html"))
	c.Assert("/section/sub/name/index.html", qt.Equals, PrettifyURLPath("/section/sub/name.html"))
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name/"))
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name/index.html"))
	c.Assert("/index.html", qt.Equals, PrettifyURLPath("/index.html"))
	c.Assert("/name/index.xml", qt.Equals, PrettifyURLPath("/name.xml"))
	c.Assert("/", qt.Equals, PrettifyURLPath("/"))
	c.Assert("/", qt.Equals, PrettifyURLPath(""))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name.html"))
	c.Assert("/section/sub/name", qt.Equals, PrettifyURL("/section/sub/name.html"))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name/"))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name/index.html"))
	c.Assert("/", qt.Equals, PrettifyURL("/index.html"))
	c.Assert("/name/index.xml", qt.Equals, PrettifyURL("/name.xml"))
	c.Assert("/", qt.Equals, PrettifyURL("/"))
	c.Assert("/", qt.Equals, PrettifyURL(""))
}

func TestUgly(t *testing.T) {
	c := qt.New(t)
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name.html"))
	c.Assert("/section/sub/name.html", qt.Equals, Uglify("/section/sub/name.html"))
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name/"))
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name/index.html"))
	c.Assert("/index.html", qt.Equals, Uglify("/index.html"))
	c.Assert("/name.xml", qt.Equals, Uglify("/name.xml"))
	c.Assert("/", qt.Equals, Uglify("/"))
	c.Assert("/", qt.Equals, Uglify(""))
}
