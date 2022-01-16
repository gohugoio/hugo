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
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		path   string
		base   string
		expect interface{}
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

func TestGetDottedRelativePath(t *testing.T) {
	// on Windows this will receive both kinds, both country and western ...
	for _, f := range []func(string) string{filepath.FromSlash, func(s string) string { return s }} {
		doTestGetDottedRelativePath(f, t)
	}
}

func doTestGetDottedRelativePath(urlFixer func(string) string, t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"", "./"},
		{urlFixer("/"), "./"},
		{urlFixer("post"), "../"},
		{urlFixer("/post"), "../"},
		{urlFixer("post/"), "../"},
		{urlFixer("tags/foo.html"), "../"},
		{urlFixer("/tags/foo.html"), "../"},
		{urlFixer("/post/"), "../"},
		{urlFixer("////post/////"), "../"},
		{urlFixer("/foo/bar/index.html"), "../../"},
		{urlFixer("/foo/bar/foo/"), "../../../"},
		{urlFixer("/foo/bar/foo"), "../../../"},
		{urlFixer("foo/bar/foo/"), "../../../"},
		{urlFixer("foo/bar/foo/bar"), "../../../../"},
		{"404.html", "./"},
		{"404.xml", "./"},
		{"/404.html", "./"},
	}
	for i, d := range data {
		output := GetDottedRelativePath(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
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
