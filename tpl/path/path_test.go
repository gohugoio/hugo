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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ns = New(&deps.Deps{Cfg: viper.New()})

type tstNoStringer struct{}

func TestBase(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		path   interface{}
		expect interface{}
	}{
		{filepath.FromSlash(`foo/bar.txt`), `bar.txt`},
		{filepath.FromSlash(`foo/bar/txt `), `txt `},
		{filepath.FromSlash(`foo\bar.txt`), filepath.FromSlash(`foo\bar.txt`)},
		{filepath.FromSlash(`foo\bar\txt `), filepath.FromSlash(`foo\bar\txt `)},
		{filepath.FromSlash(`foo/bar.t`), `bar.t`},
		{`foo.bar.txt`, `foo.bar.txt`},
		{`.x`, `.x`},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {
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

func TestDir(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		path   interface{}
		expect interface{}
	}{
		{filepath.FromSlash(`foo/bar.txt`), `foo`},
		{filepath.FromSlash(`foo/bar/txt `), `foo/bar`},
		{filepath.FromSlash(`foo\bar.txt`), `.`},
		{filepath.FromSlash(`foo\bar\txt `), `.`},
		{filepath.FromSlash(`foo/bar.t`), `foo`},
		{`foo.bar.txt`, `.`},
		{`.x`, `.`},
		{``, `.`},
		// errors
		{tstNoStringer{}, false},
	} {
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

	for i, test := range []struct {
		path   interface{}
		expect interface{}
	}{
		{filepath.FromSlash(`foo/bar.json`), `.json`},
		{filepath.FromSlash(`foo\bar.txt`), `.txt`},
		{filepath.FromSlash(`foo\bar txt`), ``},
		{`foo.bar.txt `, `.txt `},
		{``, ``},
		{`.x`, `.x`},
		// errors
		{tstNoStringer{}, false},
	} {
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

func TestJoin(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		elements interface{}
		expect   interface{}
	}{
		{
			[]string{"", "baz", filepath.FromSlash(`foo/bar.txt`)},
			filepath.FromSlash(`baz/foo/bar.txt`),
		},
		{
			[]interface{}{"", "baz", DirFile{"big", "john"}, filepath.FromSlash(`foo/bar.txt`)},
			filepath.FromSlash(`baz/big|john/foo/bar.txt`),
		},
		{nil, ""},
		// errors
		{tstNoStringer{}, false},
		{[]interface{}{"", tstNoStringer{}}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Join(test.elements)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		path   interface{}
		expect interface{}
	}{
		{filepath.FromSlash(`foo/bar.txt`), DirFile{filepath.FromSlash(`foo/`), `bar.txt`}},
		{filepath.FromSlash(`foo/bar/txt `), DirFile{filepath.FromSlash(`foo/bar/`), `txt `}},
		{`foo.bar.txt`, DirFile{``, `foo.bar.txt`}},
		{``, DirFile{``, ``}},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Split(test.path)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
