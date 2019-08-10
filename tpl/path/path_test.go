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
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/viper"
)

var ns = New(&deps.Deps{Cfg: viper.New()})

type tstNoStringer struct{}

func TestBase(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		path   interface{}
		expect interface{}
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

func TestDir(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		path   interface{}
		expect interface{}
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

	for _, test := range []struct {
		path   interface{}
		expect interface{}
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

	for _, test := range []struct {
		elements interface{}
		expect   interface{}
	}{
		{
			[]string{"", "baz", filepath.FromSlash(`foo/bar.txt`)},
			`baz/foo/bar.txt`,
		},
		{
			[]interface{}{"", "baz", DirFile{"big", "john"}, filepath.FromSlash(`foo/bar.txt`)},
			`baz/big|john/foo/bar.txt`,
		},
		{nil, ""},
		// errors
		{tstNoStringer{}, false},
		{[]interface{}{"", tstNoStringer{}}, false},
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

	for _, test := range []struct {
		path   interface{}
		expect interface{}
	}{
		{filepath.FromSlash(`foo/bar.txt`), DirFile{`foo/`, `bar.txt`}},
		{filepath.FromSlash(`foo/bar/txt `), DirFile{`foo/bar/`, `txt `}},
		{`foo.bar.txt`, DirFile{``, `foo.bar.txt`}},
		{``, DirFile{``, ``}},
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
