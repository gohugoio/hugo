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

package os

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestReadFile(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	workingDir := "/home/hugo"

	v := viper.New()
	v.Set("workingDir", workingDir)

	// f := newTestFuncsterWithViper(v)
	ns := New(&deps.Deps{Fs: hugofs.NewMem(v)})

	afero.WriteFile(ns.deps.Fs.Source, filepath.Join(workingDir, "/f/f1.txt"), []byte("f1-content"), 0755)
	afero.WriteFile(ns.deps.Fs.Source, filepath.Join("/home", "f2.txt"), []byte("f2-content"), 0755)

	for _, test := range []struct {
		filename string
		expect   interface{}
	}{
		{filepath.FromSlash("/f/f1.txt"), "f1-content"},
		{filepath.FromSlash("f/f1.txt"), "f1-content"},
		{filepath.FromSlash("../f2.txt"), false},
		{"", false},
		{"b", false},
	} {

		result, err := ns.ReadFile(test.filename)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestFileExists(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	workingDir := "/home/hugo"

	v := viper.New()
	v.Set("workingDir", workingDir)

	ns := New(&deps.Deps{Fs: hugofs.NewMem(v)})

	afero.WriteFile(ns.deps.Fs.Source, filepath.Join(workingDir, "/f/f1.txt"), []byte("f1-content"), 0755)
	afero.WriteFile(ns.deps.Fs.Source, filepath.Join("/home", "f2.txt"), []byte("f2-content"), 0755)

	for _, test := range []struct {
		filename string
		expect   interface{}
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
	c := qt.New(t)
	workingDir := "/home/hugo"

	v := viper.New()
	v.Set("workingDir", workingDir)

	ns := New(&deps.Deps{Fs: hugofs.NewMem(v)})

	afero.WriteFile(ns.deps.Fs.Source, filepath.Join(workingDir, "/f/f1.txt"), []byte("f1-content"), 0755)

	for _, test := range []struct {
		filename string
		expect   interface{}
	}{
		{filepath.FromSlash("/f/f1.txt"), int64(10)},
		{filepath.FromSlash("f/f1.txt"), int64(10)},
		{"b", nil},
		{"", nil},
	} {
		result, err := ns.Stat(test.filename)

		if test.expect == nil {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result.Size(), qt.Equals, test.expect)
	}
}
