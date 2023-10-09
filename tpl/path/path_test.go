// Copyright 2018 The Hugo Authors. All rights reserved.
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

package path

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config/testconfig"
)

func newNs() *Namespace {
	return New(testconfig.GetTestDeps(nil, nil))
}

type tstNoStringer struct{}

func TestBase(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.txt`), `bar.txt`},
		{filepath.FromSlash(`foo/bar/txt `), `txt `},
		{filepath.FromSlash(`foo/bar.t`), `bar.t`},
		{`foo.bar.txt`, `foo.bar.txt`},
		{`.x`, `.x`},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Base(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestBaseName(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.txt`), `bar`},
		{filepath.FromSlash(`foo/bar/txt `), `txt `},
		{filepath.FromSlash(`foo/bar.t`), `bar`},
		{`foo.bar.txt`, `foo.bar`},
		{`.x`, ``},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.BaseName(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestDir(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.txt`), `foo`},
		{filepath.FromSlash(`foo/bar/txt `), `foo/bar`},
		{filepath.FromSlash(`foo/bar.t`), `foo`},
		{`foo.bar.txt`, `.`},
		{`.x`, `.`},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Dir(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestExt(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.json`), `.json`},
		{`foo.bar.txt `, `.txt `},
		{``, ``},
		{`.x`, `.x`},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Ext(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		elements any
		expect   any
	}{
		{
			[]string{"", "baz", filepath.FromSlash(`foo/bar.txt`)},
			`baz/foo/bar.txt`,
		},
		{
			[]any{"", "baz", paths.DirFile{Dir: "big", File: "john"}, filepath.FromSlash(`foo/bar.txt`)},
			`baz/big|john/foo/bar.txt`,
		},
		{nil, ""},
		// errors
		{tstNoStringer{}, false},
		{[]any{"", tstNoStringer{}}, false},
	} {

		result, err := ns.Join(test.elements)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.txt`), paths.DirFile{Dir: `foo/`, File: `bar.txt`}},
		{filepath.FromSlash(`foo/bar/txt `), paths.DirFile{Dir: `foo/bar/`, File: `txt `}},
		{`foo.bar.txt`, paths.DirFile{Dir: ``, File: `foo.bar.txt`}},
		{``, paths.DirFile{Dir: ``, File: ``}},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Split(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestClean(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		path   any
		expect any
	}{
		{filepath.FromSlash(`foo/bar.txt`), `foo/bar.txt`},
		{filepath.FromSlash(`foo/bar/txt`), `foo/bar/txt`},
		{filepath.FromSlash(`foo/bar`), `foo/bar`},
		{filepath.FromSlash(`foo/bar.t`), `foo/bar.t`},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Clean(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}
