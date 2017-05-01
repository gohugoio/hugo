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

package safe

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tstNoStringer struct{}

func TestCSS(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{`a[href =~ "//example.com"]#foo`, template.CSS(`a[href =~ "//example.com"]#foo`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.CSS(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestHTML(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{`Hello, <b>World</b> &amp;tc!`, template.HTML(`Hello, <b>World</b> &amp;tc!`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.HTML(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestHTMLAttr(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{` dir="ltr"`, template.HTMLAttr(` dir="ltr"`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.HTMLAttr(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestJS(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{`c && alert("Hello, World!");`, template.JS(`c && alert("Hello, World!");`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.JS(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestJSStr(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{`Hello, World & O'Reilly\x21`, template.JSStr(`Hello, World & O'Reilly\x21`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.JSStr(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestURL(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{`greeting=H%69&addressee=(World)`, template.URL(`greeting=H%69&addressee=(World)`)},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.URL(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestSanitizeURL(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		a      interface{}
		expect interface{}
	}{
		{"http://foo/../../bar", "http://foo/bar"},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.SanitizeURL(test.a)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
