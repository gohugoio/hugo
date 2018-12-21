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

package transform

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tstNoStringer struct{}

func TestEmojify(t *testing.T) {
	t.Parallel()

	v := viper.New()
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		expect interface{}
	}{
		{":notamoji:", template.HTML(":notamoji:")},
		{"I :heart: Hugo", template.HTML("I ❤️ Hugo")},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %s", i, test.s)

		result, err := ns.Emojify(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestHighlight(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		lang   string
		opts   string
		expect interface{}
	}{
		{"func boo() {}", "go", "", "boo"},
		// Issue #4179
		{`<Foo attr=" &lt; "></Foo>`, "xml", "", `&amp;lt;`},
		{tstNoStringer{}, "go", "", false},
	} {
		errMsg := fmt.Sprintf("[%d]", i)

		result, err := ns.Highlight(test.s, test.lang, test.opts)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Contains(t, result, test.expect.(string), errMsg)
	}
}

func TestHTMLEscape(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		expect interface{}
	}{
		{`"Foo & Bar's Diner" <y@z>`, `&#34;Foo &amp; Bar&#39;s Diner&#34; &lt;y@z&gt;`},
		{"Hugo & Caddy > Wordpress & Apache", "Hugo &amp; Caddy &gt; Wordpress &amp; Apache"},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %s", i, test.s)

		result, err := ns.HTMLEscape(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestHTMLUnescape(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		expect interface{}
	}{
		{`&quot;Foo &amp; Bar&#39;s Diner&quot; &lt;y@z&gt;`, `"Foo & Bar's Diner" <y@z>`},
		{"Hugo &amp; Caddy &gt; Wordpress &amp; Apache", "Hugo & Caddy > Wordpress & Apache"},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %s", i, test.s)

		result, err := ns.HTMLUnescape(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestMarkdownify(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		expect interface{}
	}{
		{"Hello **World!**", template.HTML("Hello <strong>World!</strong>")},
		{[]byte("Hello Bytes **World!**"), template.HTML("Hello Bytes <strong>World!</strong>")},
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %s", i, test.s)

		result, err := ns.Markdownify(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

// Issue #3040
func TestMarkdownifyBlocksOfText(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	text := `
#First 

This is some *bold* text.

## Second

This is some more text.

And then some.
`

	result, err := ns.Markdownify(text)
	assert.NoError(err)
	assert.Equal(template.HTML(
		"<p>#First</p>\n\n<p>This is some <em>bold</em> text.</p>\n\n<h2 id=\"second\">Second</h2>\n\n<p>This is some more text.</p>\n\n<p>And then some.</p>\n"),
		result)

}

func TestPlainify(t *testing.T) {
	t.Parallel()

	v := viper.New()
	ns := New(newDeps(v))

	for i, test := range []struct {
		s      interface{}
		expect interface{}
	}{
		{"<em>Note:</em> blah <b>blah</b>", "Note: blah blah"},
		// errors
		{tstNoStringer{}, false},
	} {
		errMsg := fmt.Sprintf("[%d] %s", i, test.s)

		result, err := ns.Plainify(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func newDeps(cfg config.Provider) *deps.Deps {
	cfg.Set("contentDir", "content")
	cfg.Set("i18nDir", "i18n")

	l := langs.NewLanguage("en", cfg)

	cs, err := helpers.NewContentSpec(l)
	if err != nil {
		panic(err)
	}

	return &deps.Deps{
		Cfg:         cfg,
		Fs:          hugofs.NewMem(l),
		ContentSpec: cs,
	}
}
