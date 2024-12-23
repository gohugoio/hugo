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


func TestStatWithMounts(t *testing.T) {
	t.Parallel()
	b := newFileTestBuilderWithMounts(t).Build()
	ns := os.New(b.H.Deps)

	for _, test := range []struct {
		filename string
		expectError  bool
		expectSize   any
		expectIsDir  any
		expectHint string
	}{
		{filepath.FromSlash("f/f1.txt"), false, int64(10), false, "must be present in workdir"},
		{filepath.FromSlash("about.pl.md"), false, int64(2), false, "must be present in contentDir"},
		{filepath.FromSlash("en/about.md"), false, int64(2), false, "must be present in contentDir"},
		{filepath.FromSlash("assets/testing.json"), false, int64(13), false, "mapping must be present"},
		{filepath.FromSlash("assets/virtual/file1.json"), false, int64(13), false, "mapping to virtual directory must be present"},
		{filepath.FromSlash("assets/file2.json"), false, int64(13), false, "mapping to directory must be present"},
		{filepath.FromSlash("assets/files/raw1.txt"), false, int64(12), false, "mapping to subdirectory must be present"},
		{filepath.FromSlash("assets/files"), false, nil, true, "directory mapping to subdirectory must be present"},
		{filepath.FromSlash("assets/virtual"), false, int64(0), true, "directory mapping to virtual directory must be present"},
		{"b", true, nil, nil, "cannot exist"},
		{"", true, nil, nil, "cannot exist"},
		{".", false, nil, true, "root directory exist"},
		{filepath.FromSlash("../content/about.pl.md"), false, int64(2), false, "must be present in contentDir (relative to content dir)"},
	} {
		for _, prefix := range []string{"", filepath.FromSlash("/")} {
			if prefix != "" && test.filename == "" {
				continue
			}
			filename := prefix + test.filename
			result, err := ns.Stat(filename)
			if test.expectError == true {
				b.Assert(err, qt.Not(qt.IsNil), qt.Commentf("file '%s' %s", filename, test.expectHint))
				continue
			}
	
			b.Assert(err, qt.IsNil,  qt.Commentf("file '%s' %s", filename, test.expectHint))
			if test.expectSize != nil {
				// size for is dir is platform dependent
				b.Assert(result.Size(), qt.Equals, test.expectSize, qt.Commentf("file '%s' invalid size", filename))
			}
			b.Assert(result.IsDir(), qt.Equals, test.expectIsDir, qt.Commentf("file '%s' must be directory", filename))	
		}
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

func newFileTestBuilderWithMounts(t *testing.T) *hugolib.IntegrationTestBuilder {
	files := `
-- f/f1.txt --
f1-content
-- home/f2.txt --
f2-content
-- hugo.toml --
baseURL = "https://example.com/"
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[module]
[[module.imports]]
path = "module1"
[[module.imports.mounts]]
source = "assets"
target = "assets/virtual"
[[module.imports.mounts]]
source = "assets/file1.json"
target = "assets/testing.json"
[[module.imports]]
path = "module2"
[languages]
[languages.en]
contentDir = "content/en"
languageName = "English"
[languages.fr]
contentDir = "content/fr"
languageName = "Français"
[languages.pl]
languageName = "Polski"

-- content/en/about.md --
EN
-- content/fr/about.md --
FR
-- content/about.pl.md --
PL
-- themes/module1/assets/file1.json --
file1-content
-- themes/module2/assets/files/raw1.txt --
raw1-content
-- themes/module2/assets/file2.json --
file2-content
`
	
	return hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			WorkingDir:  "/mywork",
			//NeedsOsFS: true,
		},
	)
}
