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

// Note: most of these test cases are from the Go stdlib filepath package.

package filepath

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ns = New()

type tstNoStringer struct{}

type PathTest struct {
	path   interface{}
	expect interface{}
}

func TestBase(t *testing.T) {
	t.Parallel()

	tests := []PathTest{
		{"", "."},
		{".", "."},
		{"/.", "."},
		{"/", "/"},
		{"////", "/"},
		{"x/", "x"},
		{"abc", "abc"},
		{"abc/def", "def"},
		{"abc/def/", "def"},
		{"a/b/.x", ".x"},
		{"a/b/c.", "c."},
		{"a/b/c.x", "c.x"},
		// errors
		{tstNoStringer{}, false},
	}

	if runtime.GOOS == "windows" {
		// make unix tests work on windows
		for i := range tests {
			if s, ok := tests[i].expect.(string); ok {
				tests[i].expect = filepath.Clean(s)
			}
		}

		// add windows tests
		for _, test := range []PathTest{
			{`c:\`, `\`},
			{`c:.`, `.`},
			{`c:\a\b`, `b`},
			{`c:a\b`, `b`},
			{`c:a\b\c`, `c`},
			{`\\host\share\`, `\`},
			{`\\host\share\a`, `a`},
			{`\\host\share\a\b`, `b`},
		} {
			tests = append(tests, test)
		}
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Base(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestClean(t *testing.T) {
	t.Parallel()

	tests := []PathTest{
		// Already clean
		{"abc", "abc"},
		{"abc/def", "abc/def"},
		{"a/b/c", "a/b/c"},
		{".", "."},
		{"..", ".."},
		{"../..", "../.."},
		{"../../abc", "../../abc"},
		{"/abc", "/abc"},
		{"/", "/"},

		// Empty is current dir
		{"", "."},

		// Remove trailing slash
		{"abc/", "abc"},
		{"abc/def/", "abc/def"},
		{"a/b/c/", "a/b/c"},
		{"./", "."},
		{"../", ".."},
		{"../../", "../.."},
		{"/abc/", "/abc"},

		// Remove doubled slash
		{"abc//def//ghi", "abc/def/ghi"},
		{"//abc", "/abc"},
		{"///abc", "/abc"},
		{"//abc//", "/abc"},
		{"abc//", "abc"},

		// Remove . elements
		{"abc/./def", "abc/def"},
		{"/./abc/def", "/abc/def"},
		{"abc/.", "abc"},

		// Remove .. elements
		{"abc/def/ghi/../jkl", "abc/def/jkl"},
		{"abc/def/../ghi/../jkl", "abc/jkl"},
		{"abc/def/..", "abc"},
		{"abc/def/../..", "."},
		{"/abc/def/../..", "/"},
		{"abc/def/../../..", ".."},
		{"/abc/def/../../..", "/"},
		{"abc/def/../../../ghi/jkl/../../../mno", "../../mno"},
		{"/../abc", "/abc"},

		// Combinations
		{"abc/./../def", "def"},
		{"abc//./../def", "def"},
		{"abc/../../././../def", "../../def"},

		// errors
		{tstNoStringer{}, false},
	}

	if runtime.GOOS == "windows" {
		// make unix tests work on windows
		for i := range tests {
			if s, ok := tests[i].expect.(string); ok {
				tests[i].expect = filepath.FromSlash(s)
			}
		}

		// add windows tests
		for _, test := range []PathTest{
			{`c:`, `c:.`},
			{`c:\`, `c:\`},
			{`c:\abc`, `c:\abc`},
			{`c:abc\..\..\.\.\..\def`, `c:..\..\def`},
			{`c:\abc\def\..\..`, `c:\`},
			{`c:\..\abc`, `c:\abc`},
			{`c:..\abc`, `c:..\abc`},
			{`\`, `\`},
			{`/`, `\`},
			{`\\i\..\c$`, `\c$`},
			{`\\i\..\i\c$`, `\i\c$`},
			{`\\i\..\I\c$`, `\I\c$`},
			{`\\host\share\foo\..\bar`, `\\host\share\bar`},
			{`//host/share/foo/../baz`, `\\host\share\baz`},
			{`\\a\b\..\c`, `\\a\b\c`},
			{`\\a\b`, `\\a\b`},
		} {
			tests = append(tests, test)
		}
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Clean(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestDir(t *testing.T) {
	t.Parallel()

	tests := []PathTest{
		{"", "."},
		{".", "."},
		{"/.", "/"},
		{"/", "/"},
		{"////", "/"},
		{"/foo", "/"},
		{"x/", "x"},
		{"abc", "."},
		{"abc/def", "abc"},
		{"a/b/.x", "a/b"},
		{"a/b/c.", "a/b"},
		{"a/b/c.x", "a/b"},
		// errors
		{tstNoStringer{}, false},
	}

	if runtime.GOOS == "windows" {
		// make unix tests work on windows
		for i := range tests {
			if s, ok := tests[i].expect.(string); ok {
				tests[i].expect = filepath.Clean(s)
			}
		}

		// add windows tests
		for _, test := range []PathTest{
			{`c:\`, `c:\`},
			{`c:.`, `c:.`},
			{`c:\a\b`, `c:\a`},
			{`c:a\b`, `c:a`},
			{`c:a\b\c`, `c:a\b`},
			{`\\host\share`, `\\host\share`},
			{`\\host\share\`, `\\host\share\`},
			{`\\host\share\a`, `\\host\share\`},
			{`\\host\share\a\b`, `\\host\share\a`},
		} {
			tests = append(tests, test)
		}
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Dir(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestExt(t *testing.T) {
	t.Parallel()

	tests := []PathTest{
		{"path.go", ".go"},
		{"path.pb.go", ".go"},
		{"a.dir/b", ""},
		{"a.dir/b.go", ".go"},
		{"a.dir/", ""},
		// errors
		{tstNoStringer{}, false},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Ext(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestFromSlash(t *testing.T) {
	t.Parallel()

	sep := ns.Separator()

	tests := []PathTest{
		{"", ""},
		{"/", sep},
		{"/a/b", fmt.Sprintf("%sa%sb", sep, sep)},
		{"a//b", fmt.Sprintf("a%s%sb", sep, sep)},
		// errors
		{tstNoStringer{}, false},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.FromSlash(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestSeparator(t *testing.T) {
	assert.Equal(t, string(filepath.Separator), ns.Separator())
}

func TestSplit(t *testing.T) {
	t.Parallel()

	tests := []PathTest{
		{"", []string{"", ""}},
		{"/", []string{"/", ""}},
		{"a", []string{"", "a"}},
		{"a/", []string{"a/", ""}},
		{"a/b", []string{"a/", "b"}},
		{"a/b/", []string{"a/b/", ""}},
		// errors
		{tstNoStringer{}, false},
	}

	if runtime.GOOS == "windows" {
		// add windows tests
		for _, test := range []PathTest{
			{`c:`, []string{`c:`, ``}},
			{`c:/`, []string{`c:/`, ``}},
			{`c:/foo`, []string{`c:/`, `foo`}},
			{`c:/foo/`, []string{`c:/foo/`, ``}},
			{`c:/foo/bar`, []string{`c:/foo/`, `bar`}},
			{`c:/foo/bar/`, []string{`c:/foo/bar/`, ``}},
			{`//host/share`, []string{`//host/share`, ``}},
			{`//host/share/`, []string{`//host/share/`, ``}},
			{`//host/share/foo`, []string{`//host/share/`, `foo`}},
			{`//host/share/foo/`, []string{`//host/share/foo/`, ``}},
			{`//host/share/foo/bar`, []string{`//host/share/foo/`, `bar`}},
			{`//host/share/foo/bar/`, []string{`//host/share/foo/bar/`, ``}},
		} {
			tests = append(tests, test)
		}
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		dir, file, err := ns.Split(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, []string{dir, file}, errMsg)
	}
}
func TestToSlash(t *testing.T) {
	t.Parallel()

	sep := ns.Separator()

	tests := []PathTest{
		{"", ""},
		{sep, "/"},
		{fmt.Sprintf("%sa%sb", sep, sep), "/a/b"},
		{fmt.Sprintf("a%s%sb", sep, sep), "a//b"},
		// errors
		{tstNoStringer{}, false},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.ToSlash(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
