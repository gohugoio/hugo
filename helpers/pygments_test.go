// Copyright 2015 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alecthomas/chroma/formatters/html"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestParsePygmentsArgs(t *testing.T) {
	assert := require.New(t)

	for i, this := range []struct {
		in                 string
		pygmentsStyle      string
		pygmentsUseClasses bool
		expect1            interface{}
	}{
		{"", "foo", true, "encoding=utf8,noclasses=false,style=foo"},
		{"style=boo,noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"Style=boo, noClasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=foo"},
		{"style=boo", "foo", true, "encoding=utf8,noclasses=false,style=boo"},
		{"boo=invalid", "foo", false, false},
		{"style", "foo", false, false},
	} {

		v := viper.New()
		v.Set("pygmentsStyle", this.pygmentsStyle)
		v.Set("pygmentsUseClasses", this.pygmentsUseClasses)
		spec, err := NewContentSpec(v)
		assert.NoError(err)

		result1, err := spec.createPygmentsOptionsString(this.in)
		if b, ok := this.expect1.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parsePygmentArgs didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
				continue
			}
			if result1 != this.expect1 {
				t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result1, this.expect1)
			}

		}
	}
}

func TestParseDefaultPygmentsArgs(t *testing.T) {
	assert := require.New(t)

	expect := "encoding=utf8,noclasses=false,style=foo"

	for i, this := range []struct {
		in                 string
		pygmentsStyle      interface{}
		pygmentsUseClasses interface{}
		pygmentsOptions    string
	}{
		{"", "foo", true, "style=override,noclasses=override"},
		{"", nil, nil, "style=foo,noclasses=false"},
		{"style=foo,noclasses=false", nil, nil, "style=override,noclasses=override"},
		{"style=foo,noclasses=false", "override", false, "style=override,noclasses=override"},
	} {
		v := viper.New()

		v.Set("pygmentsOptions", this.pygmentsOptions)

		if s, ok := this.pygmentsStyle.(string); ok {
			v.Set("pygmentsStyle", s)
		}

		if b, ok := this.pygmentsUseClasses.(bool); ok {
			v.Set("pygmentsUseClasses", b)
		}

		spec, err := NewContentSpec(v)
		assert.NoError(err)

		result, err := spec.createPygmentsOptionsString(this.in)
		if err != nil {
			t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
			continue
		}
		if result != expect {
			t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result, expect)
		}
	}
}

type chromaInfo struct {
	classes            bool
	lineNumbers        bool
	lineNumbersInTable bool
	highlightRangesLen int
	highlightRangesStr string
	baseLineNumber     int
}

func formatterChromaInfo(f *html.Formatter) chromaInfo {
	v := reflect.ValueOf(f).Elem()
	c := chromaInfo{}
	// Hack:

	c.classes = f.Classes
	c.lineNumbers = v.FieldByName("lineNumbers").Bool()
	c.lineNumbersInTable = v.FieldByName("lineNumbersInTable").Bool()
	c.baseLineNumber = int(v.FieldByName("baseLineNumber").Int())
	vv := v.FieldByName("highlightRanges")
	c.highlightRangesLen = vv.Len()
	c.highlightRangesStr = fmt.Sprint(vv)

	return c
}

func TestChromaHTMLHighlight(t *testing.T) {
	assert := require.New(t)

	v := viper.New()
	v.Set("pygmentsUseClasses", true)
	spec, err := NewContentSpec(v)
	assert.NoError(err)

	result, err := spec.Highlight(`echo "Hello"`, "bash", "")
	assert.NoError(err)

	assert.Contains(result, `<div class="highlight"><pre class="chroma"><code class="language-bash" data-lang="bash"><span class="nb">echo</span> <span class="s2">&#34;Hello&#34;</span></code></pre></div>`)

}

func TestChromaHTMLFormatterFromOptions(t *testing.T) {
	assert := require.New(t)

	for i, this := range []struct {
		in                 string
		pygmentsStyle      interface{}
		pygmentsUseClasses interface{}
		pygmentsOptions    string
		assert             func(c chromaInfo)
	}{
		{"", "monokai", true, "style=manni,noclasses=true", func(c chromaInfo) {
			assert.True(c.classes)
			assert.False(c.lineNumbers)
			assert.Equal(0, c.highlightRangesLen)

		}},
		{"", nil, nil, "style=monokai,noclasses=false", func(c chromaInfo) {
			assert.True(c.classes)
		}},
		{"linenos=sure,hl_lines=1 2 3", nil, nil, "style=monokai,noclasses=false", func(c chromaInfo) {
			assert.True(c.classes)
			assert.True(c.lineNumbers)
			assert.Equal(3, c.highlightRangesLen)
			assert.Equal("[[1 1] [2 2] [3 3]]", c.highlightRangesStr)
			assert.Equal(1, c.baseLineNumber)
		}},
		{"linenos=inline,hl_lines=1,linenostart=4", nil, nil, "style=monokai,noclasses=false", func(c chromaInfo) {
			assert.True(c.classes)
			assert.True(c.lineNumbers)
			assert.False(c.lineNumbersInTable)
			assert.Equal(1, c.highlightRangesLen)
			// This compansates for https://github.com/alecthomas/chroma/issues/30
			assert.Equal("[[4 4]]", c.highlightRangesStr)
			assert.Equal(4, c.baseLineNumber)
		}},
		{"linenos=table", nil, nil, "style=monokai", func(c chromaInfo) {
			assert.True(c.lineNumbers)
			assert.True(c.lineNumbersInTable)
		}},
		{"style=monokai,noclasses=false", nil, nil, "style=manni,noclasses=true", func(c chromaInfo) {
			assert.True(c.classes)
		}},
		{"style=monokai,noclasses=true", "friendly", false, "style=manni,noclasses=false", func(c chromaInfo) {
			assert.False(c.classes)
		}},
	} {
		v := viper.New()

		v.Set("pygmentsOptions", this.pygmentsOptions)

		if s, ok := this.pygmentsStyle.(string); ok {
			v.Set("pygmentsStyle", s)
		}

		if b, ok := this.pygmentsUseClasses.(bool); ok {
			v.Set("pygmentsUseClasses", b)
		}

		spec, err := NewContentSpec(v)
		assert.NoError(err)

		opts, err := spec.parsePygmentsOpts(this.in)
		if err != nil {
			t.Fatalf("[%d] parsePygmentsOpts failed: %s", i, err)
		}

		chromaFormatter, err := spec.chromaFormatterFromOptions(opts)
		if err != nil {
			t.Fatalf("[%d] chromaFormatterFromOptions failed: %s", i, err)
		}

		this.assert(formatterChromaInfo(chromaFormatter.(*html.Formatter)))
	}
}

func TestHlLinesToRanges(t *testing.T) {
	var zero [][2]int

	for _, this := range []struct {
		in        string
		startLine int
		expected  interface{}
	}{
		{"", 1, zero},
		{"1 4", 1, [][2]int{{1, 1}, {4, 4}}},
		{"1 4", 2, [][2]int{{2, 2}, {5, 5}}},
		{"1-4 5-8", 1, [][2]int{{1, 4}, {5, 8}}},
		{" 1   4 ", 1, [][2]int{{1, 1}, {4, 4}}},
		{"1-4    5-8 ", 1, [][2]int{{1, 4}, {5, 8}}},
		{"1-4 5", 1, [][2]int{{1, 4}, {5, 5}}},
		{"4 5-9", 1, [][2]int{{4, 4}, {5, 9}}},
		{" 1 -4 5 - 8  ", 1, true},
		{"a b", 1, true},
	} {
		got, err := hlLinesToRanges(this.startLine, this.in)

		if expectErr, ok := this.expected.(bool); ok && expectErr {
			if err == nil {
				t.Fatal("No error")
			}
		} else if err != nil {
			t.Fatalf("Got error: %s", err)
		} else if !reflect.DeepEqual(this.expected, got) {
			t.Fatalf("Expected\n%v but got\n%v", this.expected, got)
		}
	}
}

func BenchmarkChromaHighlight(b *testing.B) {
	assert := require.New(b)
	v := viper.New()

	v.Set("pygmentsstyle", "trac")
	v.Set("pygmentsuseclasses", false)
	v.Set("pygmentsuseclassic", false)

	code := `// GetTitleFunc returns a func that can be used to transform a string to
// title case.
//
// The supported styles are
//
// - "Go" (strings.Title)
// - "AP" (see https://www.apstylebook.com/)
// - "Chicago" (see http://www.chicagomanualofstyle.org/home.html)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
  switch strings.ToLower(style) {
  case "go":
    return strings.Title
  case "chicago":
    tc := transform.NewTitleConverter(transform.ChicagoStyle)
    return tc.Title
  default:
    tc := transform.NewTitleConverter(transform.APStyle)
    return tc.Title
  }
}
`

	spec, err := NewContentSpec(v)
	assert.NoError(err)

	for i := 0; i < b.N; i++ {
		_, err := spec.Highlight(code, "go", "linenos=inline,hl_lines=8 15-17")
		if err != nil {
			b.Fatal(err)
		}
	}
}
