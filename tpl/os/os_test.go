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

package os_test

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/os"

	qt "github.com/frankban/quicktest"
)

func TestReadFile(t *testing.T) {
	t.Parallel()

	b := newFileTestBuilder(t).Build()

	// helpers.PrintFs(b.H.PathSpec.BaseFs.Work, "", _os.Stdout)

	ns := os.New(b.H.Deps)

	for _, test := range []struct {
		filename string
		expect   any
	}{
		{filepath.FromSlash("/f/f1.txt"), "f1-content"},
		{filepath.FromSlash("f/f1.txt"), "f1-content"},
		{filepath.FromSlash("../f2.txt"), ""},
		{"", false},
		{"b", ""},
	} {

		result, err := ns.ReadFile(test.filename)

		if bb, ok := test.expect.(bool); ok && !bb {
			b.Assert(err, qt.Not(qt.IsNil), qt.Commentf("filename: %q", test.filename))
			continue
		}

		b.Assert(err, qt.IsNil)
		b.Assert(result, qt.Equals, test.expect)
	}
}

func TestFileExists(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	b := newFileTestBuilder(t).Build()
	ns := os.New(b.H.Deps)

	for _, test := range []struct {
		filename string
		expect   any
	}{
		{filepath.FromSlash("/f/f1.txt"), true},
		{filepath.FromSlash("f/f1.txt"), true},
		{filepath.FromSlash("../f2.txt"), false},
		{"b", false},
		{"", nil},
	} {
		result, err := ns.FileExists(test.filename)

		if test.expect == nil {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestStat(t *testing.T) {
	t.Parallel()
	b := newFileTestBuilder(t).Build()
	ns := os.New(b.H.Deps)

	for _, test := range []struct {
		filename string
		expect   any
	}{
		{filepath.FromSlash("/f/f1.txt"), int64(10)},
		{filepath.FromSlash("f/f1.txt"), int64(10)},
		{"b", nil},
		{"", nil},
	} {
		result, err := ns.Stat(test.filename)

		if test.expect == nil {
			b.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		b.Assert(err, qt.IsNil)
		b.Assert(result.Size(), qt.Equals, test.expect)
	}
}

func newFileTestBuilder(t *testing.T) *hugolib.IntegrationTestBuilder {
	files := `
-- f/f1.txt --
f1-content
-- home/f2.txt --
f2-content
	`

	return hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			WorkingDir:  "/mywork",
		},
	)
}
