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

package helpers_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
)

func TestMakePath(t *testing.T) {
	tests := []struct {
		input         string
		expected      string
		removeAccents bool
	}{
		{"dot.slash/backslash\\underscore_pound#plus+hyphen-", "dot.slash/backslash\\underscore_pound#plus+hyphen-", true},
		{"abcXYZ0123456789", "abcXYZ0123456789", true},
		{"%20 %2", "%20-2", true},
		{"foo- bar", "foo-bar", true},
		{"  Foo bar  ", "Foo-bar", true},
		{"Foo.Bar/foo_Bar-Foo", "Foo.Bar/foo_Bar-Foo", true},
		{"fOO,bar:foobAR", "fOObarfoobAR", true},
		{"FOo/BaR.html", "FOo/BaR.html", true},
		{"трям/трям", "трям/трям", true},
		{"은행", "은행", true},
		{"Банковский кассир", "Банковскии-кассир", true},
		// Issue #1488
		{"संस्कृत", "संस्कृत", false},
		{"a%C3%B1ame", "a%C3%B1ame", false},         // Issue #1292
		{"this+is+a+test", "this+is+a+test", false}, // Issue #1290
		{"~foo", "~foo", false},                     // Issue #2177
		{"foo--bar", "foo--bar", true},              // Issue #7288
		{"foo@bar", "foo@bar", true},                //	Issue #10548
	}

	for _, test := range tests {
		p := newTestPathSpec("removePathAccents", test.removeAccents)
		output := p.MakePath(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestMakePathSanitized(t *testing.T) {
	p := newTestPathSpec()

	tests := []struct {
		input    string
		expected string
	}{
		{"  FOO bar  ", "foo-bar"},
		{"Foo.Bar/fOO_bAr-Foo", "foo.bar/foo_bar-foo"},
		{"FOO,bar:FooBar", "foobarfoobar"},
		{"foo/BAR.HTML", "foo/bar.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
	}

	for _, test := range tests {
		output := p.MakePathSanitized(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestMakePathSanitizedDisablePathToLower(t *testing.T) {
	p := newTestPathSpec("disablePathToLower", true)

	tests := []struct {
		input    string
		expected string
	}{
		{"  FOO bar  ", "FOO-bar"},
		{"Foo.Bar/fOO_bAr-Foo", "Foo.Bar/fOO_bAr-Foo"},
		{"FOO,bar:FooBar", "FOObarFooBar"},
		{"foo/BAR.HTML", "foo/BAR.HTML"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
	}

	for _, test := range tests {
		output := p.MakePathSanitized(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
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
		output, _ := helpers.MakePathRelative(d.inPath, d.path1, d.path2)
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}
	_, error := helpers.MakePathRelative("a/b/c.ss", "/a/c", "/d/c", "/e/f")

	if error == nil {
		t.Errorf("Test failed, expected error")
	}
}

func TestGetDottedRelativePath(t *testing.T) {
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
		output := helpers.GetDottedRelativePath(d.input)
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
		output := helpers.MakeTitle(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func TestDirExists(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	data := []test{
		{".", true},
		{"./", true},
		{"..", true},
		{"../", true},
		{"./..", true},
		{"./../", true},
		{os.TempDir(), true},
		{os.TempDir() + helpers.FilePathSeparator, true},
		{"/", true},
		{"/some-really-random-directory-name", false},
		{"/some/really/random/directory/name", false},
		{"./some-really-random-local-directory-name", false},
		{"./some/really/random/local/directory/name", false},
	}

	for i, d := range data {
		exists, _ := helpers.DirExists(filepath.FromSlash(d.input), new(afero.OsFs))
		if d.expected != exists {
			t.Errorf("Test %d failed. Expected %t got %t", i, d.expected, exists)
		}
	}
}

func TestIsDir(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}
	data := []test{
		{"./", true},
		{"/", true},
		{"./this-directory-does-not-existi", false},
		{"/this-absolute-directory/does-not-exist", false},
	}

	for i, d := range data {

		exists, _ := helpers.IsDir(d.input, new(afero.OsFs))
		if d.expected != exists {
			t.Errorf("Test %d failed. Expected %t got %t", i, d.expected, exists)
		}
	}
}

func createZeroSizedFileInTempDir(t *testing.T) *os.File {
	t.Helper()

	filePrefix := "_path_test_"
	f, err := os.CreateTemp(t.TempDir(), filePrefix)
	if err != nil {
		t.Error(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}
	return f
}

func createNonZeroSizedFileInTempDir(t *testing.T) *os.File {
	t.Helper()

	f := createZeroSizedFileInTempDir(t)
	byteString := []byte("byteString")
	err := os.WriteFile(f.Name(), byteString, 0o644)
	if err != nil {
		t.Error(err)
	}
	return f
}

func TestExists(t *testing.T) {
	zeroSizedFile := createZeroSizedFileInTempDir(t)
	nonZeroSizedFile := createNonZeroSizedFileInTempDir(t)
	emptyDirectory := t.TempDir()
	nonExistentFile := os.TempDir() + "/this-file-does-not-exist.txt"
	nonExistentDir := os.TempDir() + "/this/directory/does/not/exist/"

	type test struct {
		input          string
		expectedResult bool
		expectedErr    error
	}

	data := []test{
		{zeroSizedFile.Name(), true, nil},
		{nonZeroSizedFile.Name(), true, nil},
		{emptyDirectory, true, nil},
		{nonExistentFile, false, nil},
		{nonExistentDir, false, nil},
	}
	for i, d := range data {
		exists, err := helpers.Exists(d.input, new(afero.OsFs))
		if d.expectedResult != exists {
			t.Errorf("Test %d failed. Expected result %t got %t", i, d.expectedResult, exists)
		}
		if d.expectedErr != err {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expectedErr, err)
		}
	}
}

func TestAbsPathify(t *testing.T) {
	type test struct {
		inPath, workingDir, expected string
	}
	data := []test{
		{os.TempDir(), filepath.FromSlash("/work"), filepath.Clean(os.TempDir())}, // TempDir has trailing slash
		{"dir", filepath.FromSlash("/work"), filepath.FromSlash("/work/dir")},
	}

	windowsData := []test{
		{"c:\\banana\\..\\dir", "c:\\foo", "c:\\dir"},
		{"\\dir", "c:\\foo", "c:\\foo\\dir"},
		{"c:\\", "c:\\foo", "c:\\"},
	}

	unixData := []test{
		{"/banana/../dir/", "/work", "/dir"},
	}

	for i, d := range data {
		// todo see comment in AbsPathify
		ps := newTestPathSpec("workingDir", d.workingDir)

		expected := ps.AbsPathify(d.inPath)
		if d.expected != expected {
			t.Errorf("Test %d failed. Expected %q but got %q", i, d.expected, expected)
		}
	}
	t.Logf("Running platform specific path tests for %s", runtime.GOOS)
	if runtime.GOOS == "windows" {
		for i, d := range windowsData {
			ps := newTestPathSpec("workingDir", d.workingDir)

			expected := ps.AbsPathify(d.inPath)
			if d.expected != expected {
				t.Errorf("Test %d failed. Expected %q but got %q", i, d.expected, expected)
			}
		}
	} else {
		for i, d := range unixData {
			ps := newTestPathSpec("workingDir", d.workingDir)

			expected := ps.AbsPathify(d.inPath)
			if d.expected != expected {
				t.Errorf("Test %d failed. Expected %q but got %q", i, d.expected, expected)
			}
		}
	}
}

func TestExtractAndGroupRootPaths(t *testing.T) {
	in := []string{
		filepath.FromSlash("/a/b/c/d"),
		filepath.FromSlash("/a/b/c/e"),
		filepath.FromSlash("/a/b/e/f"),
		filepath.FromSlash("/a/b"),
		filepath.FromSlash("/a/b/c/b/g"),
		filepath.FromSlash("/c/d/e"),
	}

	inCopy := make([]string, len(in))
	copy(inCopy, in)

	result := helpers.ExtractAndGroupRootPaths(in)

	c := qt.New(t)
	c.Assert(fmt.Sprint(result), qt.Equals, filepath.FromSlash("[/a/b/{c,e} /c/d/e]"))

	// Make sure the original is preserved
	c.Assert(in, qt.DeepEquals, inCopy)
}

func TestExtractRootPaths(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{{
		[]string{
			filepath.FromSlash("a/b"), filepath.FromSlash("a/b/c/"), "b",
			filepath.FromSlash("/c/d"), filepath.FromSlash("d/"), filepath.FromSlash("//e//"),
		},
		[]string{"a", "a", "b", "c", "d", "e"},
	}}

	for _, test := range tests {
		output := helpers.ExtractRootPaths(test.input)
		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestFindCWD(t *testing.T) {
	type test struct {
		expectedDir string
		expectedErr error
	}

	// cwd, _ := os.Getwd()
	data := []test{
		//{cwd, nil},
		// Commenting this out. It doesn't work properly.
		// There's a good reason why we don't use os.Getwd(), it doesn't actually work the way we want it to.
		// I really don't know a better way to test this function. - SPF 2014.11.04
	}
	for i, d := range data {
		dir, err := helpers.FindCWD()
		if d.expectedDir != dir {
			t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedDir, dir)
		}
		if d.expectedErr != err {
			t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedErr, err)
		}
	}
}

func TestSafeWriteToDisk(t *testing.T) {
	emptyFile := createZeroSizedFileInTempDir(t)
	tmpDir := t.TempDir()

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	fileExists := fmt.Errorf("%v already exists", emptyFile.Name())

	type test struct {
		filename    string
		expectedErr error
	}

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	data := []test{
		{emptyFile.Name(), fileExists},
		{tmpDir + "/" + nowStr, nil},
	}

	for i, d := range data {
		e := helpers.SafeWriteToDisk(d.filename, reader, new(afero.OsFs))
		if d.expectedErr != nil {
			if d.expectedErr.Error() != e.Error() {
				t.Errorf("Test %d failed. Expected error %q but got %q", i, d.expectedErr.Error(), e.Error())
			}
		} else {
			if d.expectedErr != e {
				t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedErr, e)
			}
			contents, _ := os.ReadFile(d.filename)
			if randomString != string(contents) {
				t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
			}
		}
		reader.Seek(0, 0)
	}
}

func TestWriteToDisk(t *testing.T) {
	emptyFile := createZeroSizedFileInTempDir(t)
	tmpDir := t.TempDir()

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	type test struct {
		filename    string
		expectedErr error
	}

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	data := []test{
		{emptyFile.Name(), nil},
		{tmpDir + "/" + nowStr, nil},
	}

	for i, d := range data {
		e := helpers.WriteToDisk(d.filename, reader, new(afero.OsFs))
		if d.expectedErr != e {
			t.Errorf("Test %d failed. WriteToDisk Error Expected %q but got %q", i, d.expectedErr, e)
		}
		contents, e := os.ReadFile(d.filename)
		if e != nil {
			t.Errorf("Test %d failed. Could not read file %s. Reason: %s\n", i, d.filename, e)
		}
		if randomString != string(contents) {
			t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
		}
		reader.Seek(0, 0)
	}
}

func TestGetTempDir(t *testing.T) {
	dir := os.TempDir()
	if helpers.FilePathSeparator != dir[len(dir)-1:] {
		dir = dir + helpers.FilePathSeparator
	}
	testDir := "hugoTestFolder" + helpers.FilePathSeparator
	tests := []struct {
		input    string
		expected string
	}{
		{"", dir},
		{testDir + "  Foo bar  ", dir + testDir + "  Foo bar  " + helpers.FilePathSeparator},
		{testDir + "Foo.Bar/foo_Bar-Foo", dir + testDir + "Foo.Bar/foo_Bar-Foo" + helpers.FilePathSeparator},
		{testDir + "fOO,bar:foo%bAR", dir + testDir + "fOObarfoo%bAR" + helpers.FilePathSeparator},
		{testDir + "fOO,bar:foobAR", dir + testDir + "fOObarfoobAR" + helpers.FilePathSeparator},
		{testDir + "FOo/BaR.html", dir + testDir + "FOo/BaR.html" + helpers.FilePathSeparator},
		{testDir + "трям/трям", dir + testDir + "трям/трям" + helpers.FilePathSeparator},
		{testDir + "은행", dir + testDir + "은행" + helpers.FilePathSeparator},
		{testDir + "Банковский кассир", dir + testDir + "Банковский кассир" + helpers.FilePathSeparator},
	}

	for _, test := range tests {
		output := helpers.GetTempDir(test.input, new(afero.MemMapFs))
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}
