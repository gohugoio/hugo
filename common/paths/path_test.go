// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		path   string
		base   string
		expect any
	}{
		{filepath.FromSlash("/a/b"), filepath.FromSlash("/a"), filepath.FromSlash("b")},
		{filepath.FromSlash("/a/b/c/"), filepath.FromSlash("/a"), filepath.FromSlash("b/c/")},
		{filepath.FromSlash("/c"), filepath.FromSlash("/a/b"), filepath.FromSlash("../../c")},
		{filepath.FromSlash("/c"), "", false},
	}
	for i, this := range tests {
		// ultimately a fancy wrapper around filepath.Rel
		result, err := GetRelativePath(this.path, this.base)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] GetRelativePath didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] GetRelativePath failed: %s", i, err)
				continue
			}
			if result != this.expect {
				t.Errorf("[%d] GetRelativePath got %v but expected %v", i, result, this.expect)
			}
		}

	}
}

func TestMakePathRelative(t *testing.T) {
	type test struct {
		inPath, path1, path2, output string
	}

	data := []test{
		{"/abc/bcd/ab.css", "/abc/bcd", "/bbc/bcd", "/ab.css"},
		{"/abc/bcd/ab.css", "/abcd/bcd", "/abc/bcd", "/ab.css"},
	}

	for i, d := range data {
		output, _ := makePathRelative(d.inPath, d.path1, d.path2)
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}
	_, error := makePathRelative("a/b/c.ss", "/a/c", "/d/c", "/e/f")

	if error == nil {
		t.Errorf("Test failed, expected error")
	}
}

func TestMakeTitle(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"Make-Title", "Make Title"},
		{"MakeTitle", "MakeTitle"},
		{"make_title", "make_title"},
	}
	for i, d := range data {
		output := MakeTitle(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

// Replace Extension is probably poorly named, but the intent of the
// function is to accept a path and return only the file name with a
// new extension. It's intentionally designed to strip out the path
// and only provide the name. We should probably rename the function to
// be more explicit at some point.
func TestReplaceExtension(t *testing.T) {
	type test struct {
		input, newext, expected string
	}
	data := []test{
		// These work according to the above definition
		{"/some/random/path/file.xml", "html", "file.html"},
		{"/banana.html", "xml", "banana.xml"},
		{"./banana.html", "xml", "banana.xml"},
		{"banana/pie/index.html", "xml", "index.xml"},
		{"../pies/fish/index.html", "xml", "index.xml"},
		// but these all fail
		{"filename-without-an-ext", "ext", "filename-without-an-ext.ext"},
		{"/filename-without-an-ext", "ext", "filename-without-an-ext.ext"},
		{"/directory/mydir/", "ext", ".ext"},
		{"mydir/", "ext", ".ext"},
	}

	for i, d := range data {
		output := ReplaceExtension(filepath.FromSlash(d.input), d.newext)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func TestExtNoDelimiter(t *testing.T) {
	c := qt.New(t)
	c.Assert(ExtNoDelimiter(filepath.FromSlash("/my/data.json")), qt.Equals, "json")
}

func TestFilename(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"index.html", "index"},
		{"./index.html", "index"},
		{"/index.html", "index"},
		{"index", "index"},
		{"/tmp/index.html", "index"},
		{"./filename-no-ext", "filename-no-ext"},
		{"/filename-no-ext", "filename-no-ext"},
		{"filename-no-ext", "filename-no-ext"},
		{"directory/", ""}, // no filename case??
		{"directory/.hidden.ext", ".hidden"},
		{"./directory/../~/banana/gold.fish", "gold"},
		{"../directory/banana.man", "banana"},
		{"~/mydir/filename.ext", "filename"},
		{"./directory//tmp/filename.ext", "filename"},
	}

	for i, d := range data {
		output := Filename(filepath.FromSlash(d.input))
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func TestFileAndExt(t *testing.T) {
	type test struct {
		input, expectedFile, expectedExt string
	}
	data := []test{
		{"index.html", "index", ".html"},
		{"./index.html", "index", ".html"},
		{"/index.html", "index", ".html"},
		{"index", "index", ""},
		{"/tmp/index.html", "index", ".html"},
		{"./filename-no-ext", "filename-no-ext", ""},
		{"/filename-no-ext", "filename-no-ext", ""},
		{"filename-no-ext", "filename-no-ext", ""},
		{"directory/", "", ""}, // no filename case??
		{"directory/.hidden.ext", ".hidden", ".ext"},
		{"./directory/../~/banana/gold.fish", "gold", ".fish"},
		{"../directory/banana.man", "banana", ".man"},
		{"~/mydir/filename.ext", "filename", ".ext"},
		{"./directory//tmp/filename.ext", "filename", ".ext"},
	}

	for i, d := range data {
		file, ext := fileAndExt(filepath.FromSlash(d.input), fpb)
		if d.expectedFile != file {
			t.Errorf("Test %d failed. Expected filename %q got %q.", i, d.expectedFile, file)
		}
		if d.expectedExt != ext {
			t.Errorf("Test %d failed. Expected extension %q got %q.", i, d.expectedExt, ext)
		}
	}
}

func TestSanitize(t *testing.T) {
	c := qt.New(t)
	tests := []struct {
		input    string
		expected string
	}{
		{"  Foo bar  ", "Foo-bar"},
		{"Foo.Bar/foo_Bar-Foo", "Foo.Bar/foo_Bar-Foo"},
		{"fOO,bar:foobAR", "fOObarfoobAR"},
		{"FOo/BaR.html", "FOo/BaR.html"},
		{"FOo/Ba---R.html", "FOo/Ba---R.html"}, /// See #10104
		{"FOo/Ba       R.html", "FOo/Ba-R.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
		{"Банковский кассир", "Банковский-кассир"},
		// Issue #1488
		{"संस्कृत", "संस्कृत"},
		{"a%C3%B1ame", "a%C3%B1ame"},         // Issue #1292
		{"this+is+a+test", "this+is+a+test"}, // Issue #1290
		{"~foo", "~foo"},                     // Issue #2177

	}

	for _, test := range tests {
		c.Assert(Sanitize(test.input), qt.Equals, test.expected)
	}
}

func BenchmarkSanitize(b *testing.B) {
	const (
		allAlowedPath = "foo/bar"
		spacePath     = "foo bar"
	)

	// This should not allocate any memory.
	b.Run("All allowed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			got := Sanitize(allAlowedPath)
			if got != allAlowedPath {
				b.Fatal(got)
			}
		}
	})

	// This will allocate some memory.
	b.Run("Spaces", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			got := Sanitize(spacePath)
			if got != "foo-bar" {
				b.Fatal(got)
			}
		}
	})
}

func TestDir(t *testing.T) {
	c := qt.New(t)
	c.Assert(Dir("/a/b/c/d"), qt.Equals, "/a/b/c")
	c.Assert(Dir("/a"), qt.Equals, "/")
	c.Assert(Dir("/"), qt.Equals, "/")
	c.Assert(Dir(""), qt.Equals, "")
}

func TestFieldsSlash(t *testing.T) {
	c := qt.New(t)

	c.Assert(FieldsSlash("a/b/c"), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(FieldsSlash("/a/b/c"), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(FieldsSlash("/a/b/c/"), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(FieldsSlash("a/b/c/"), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(FieldsSlash("/"), qt.DeepEquals, []string{})
	c.Assert(FieldsSlash(""), qt.DeepEquals, []string{})
}

func TestCommonDirPath(t *testing.T) {
	c := qt.New(t)

	for _, this := range []struct {
		a, b, expected string
	}{
		{"/a/b/c", "/a/b/d", "/a/b"},
		{"/a/b/c", "a/b/d", "/a/b"},
		{"a/b/c", "/a/b/d", "/a/b"},
		{"a/b/c", "a/b/d", "a/b"},
		{"/a/b/c", "/a/b/c", "/a/b/c"},
		{"/a/b/c", "/a/b/c/d", "/a/b/c"},
		{"/a/b/c", "/a/b", "/a/b"},
		{"/a/b/c", "/a", "/a"},
		{"/a/b/c", "/d/e/f", ""},
	} {
		c.Assert(CommonDirPath(this.a, this.b), qt.Equals, this.expected, qt.Commentf("a: %s b: %s", this.a, this.b))
	}
}

func TestIsSameFilePath(t *testing.T) {
	c := qt.New(t)

	for _, this := range []struct {
		a, b     string
		expected bool
	}{
		{"/a/b/c", "/a/b/c", true},
		{"/a/b/c", "/a/b/c/", true},
		{"/a/b/c", "/a/b/d", false},
		{"/a/b/c", "/a/b", false},
		{"/a/b/c", "/a/b/c/d", false},
		{"/a/b/c", "/a/b/cd", false},
		{"/a/b/c", "/a/b/cc", false},
		{"/a/b/c", "/a/b/c/", true},
		{"/a/b/c", "/a/b/c//", true},
		{"/a/b/c", "/a/b/c/.", true},
		{"/a/b/c", "/a/b/c/./", true},
		{"/a/b/c", "/a/b/c/./.", true},
		{"/a/b/c", "/a/b/c/././", true},
		{"/a/b/c", "/a/b/c/././.", true},
		{"/a/b/c", "/a/b/c/./././", true},
		{"/a/b/c", "/a/b/c/./././.", true},
		{"/a/b/c", "/a/b/c/././././", true},
	} {
		c.Assert(IsSameFilePath(filepath.FromSlash(this.a), filepath.FromSlash(this.b)), qt.Equals, this.expected, qt.Commentf("a: %s b: %s", this.a, this.b))
	}
}
