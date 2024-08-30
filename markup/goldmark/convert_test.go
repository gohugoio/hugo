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

package goldmark_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark"

	"github.com/gohugoio/hugo/markup/highlight"

	"github.com/gohugoio/hugo/markup/markup_config"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

var cfgStrHighlichgtNoClasses = `
[markup]
[markup.highlight]
noclasses=false
`

func convert(c *qt.C, conf config.AllProvider, content string) converter.ResultRender {
	pconf := converter.ProviderConfig{
		Logger: loggers.NewDefault(),
		Conf:   conf,
	}

	p, err := goldmark.Provider.New(
		pconf,
	)
	c.Assert(err, qt.IsNil)

	mconf := pconf.MarkupConfig()

	h := highlight.New(mconf.Highlight)

	getRenderer := func(t hooks.RendererType, id any) any {
		switch t {
		case hooks.CodeBlockRendererType:
			return h
		case hooks.TableRendererType:
			return tableRenderer(0)
		}

		return nil
	}

	conv, err := p.New(converter.DocumentContext{DocumentID: "thedoc"})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{RenderTOC: true, Src: []byte(content), GetRenderer: getRenderer})
	c.Assert(err, qt.IsNil)

	return b
}

func TestConvert(t *testing.T) {
	c := qt.New(t)

	// Smoke test of the default configuration.
	content := `
## Links

https://github.com/gohugoio/hugo/issues/6528
[Live Demo here!](https://docuapi.netlify.com/)

[I'm an inline-style link with title](https://www.google.com "Google's Homepage")
<https://foo.bar/>
https://bar.baz/
<fake@example.com>
<mailto:fake2@example.com>


## Code Fences

§§§bash
LINE1
§§§

## Code Fences No Lexer

§§§moo
LINE1
§§§

## Custom ID {#custom}

## Auto ID

* Autolink: https://gohugo.io/
* Strikethrough:~~Hi~~ Hello, world!

## Table

| foo | bar |
| --- | --- |
| baz | bim |

## Task Lists (default on)

- [x] Finish my changes[^1]
- [ ] Push my commits to GitHub
- [ ] Open a pull request


## Smartypants (default on)

* Straight double "quotes" and single 'quotes' into “curly” quote HTML entities
* Dashes (“--” and “---”) into en- and em-dash entities
* Three consecutive dots (“...”) into an ellipsis entity
* Apostrophes are also converted: "That was back in the '90s, that's a long time ago"

## Footnotes

That's some text with a footnote.[^1]

## Definition Lists

date
: the datetime assigned to this page.

description
: the description for the content.


## 神真美好

## 神真美好

## 神真美好

[^1]: And that's the footnote.

`

	// Code fences
	content = strings.Replace(content, "§§§", "```", -1)

	cfg := config.FromTOMLConfigString(`
[markup]
[markup.highlight]
noClasses = false
[markup.goldmark.renderer]
unsafe = true

`)

	b := convert(c, testconfig.GetTestConfig(nil, cfg), content)
	got := string(b.Bytes())

	// Links
	c.Assert(got, qt.Contains, `<a href="https://docuapi.netlify.com/">Live Demo here!</a>`)
	c.Assert(got, qt.Contains, `<a href="https://foo.bar/">https://foo.bar/</a>`)
	c.Assert(got, qt.Contains, `<a href="https://bar.baz/">https://bar.baz/</a>`)
	c.Assert(got, qt.Contains, `<a href="mailto:fake@example.com">fake@example.com</a>`)
	c.Assert(got, qt.Contains, `<a href="mailto:fake2@example.com">mailto:fake2@example.com</a></p>`)

	// Header IDs
	c.Assert(got, qt.Contains, `<h2 id="custom">Custom ID</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="auto-id">Auto ID</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好">神真美好</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好-1">神真美好</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好-2">神真美好</h2>`, qt.Commentf(got))

	// Code fences
	c.Assert(got, qt.Contains, "<div class=\"highlight\"><pre tabindex=\"0\" class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\"><span class=\"line\"><span class=\"cl\">LINE1\n</span></span></code></pre></div>")
	c.Assert(got, qt.Contains, "Code Fences No Lexer</h2>\n<pre tabindex=\"0\"><code class=\"language-moo\" data-lang=\"moo\">LINE1\n</code></pre>")

	// Extensions
	c.Assert(got, qt.Contains, `Autolink: <a href="https://gohugo.io/">https://gohugo.io/</a>`)
	c.Assert(got, qt.Contains, `Strikethrough:<del>Hi</del> Hello, world`)
	c.Assert(got, qt.Contains, `Table`)
	c.Assert(got, qt.Contains, `<li><input disabled="" type="checkbox"> Push my commits to GitHub</li>`)

	c.Assert(got, qt.Contains, `Straight double &ldquo;quotes&rdquo; and single &lsquo;quotes&rsquo;`)
	c.Assert(got, qt.Contains, `Dashes (“&ndash;” and “&mdash;”) `)
	c.Assert(got, qt.Contains, `Three consecutive dots (“&hellip;”)`)
	c.Assert(got, qt.Contains, `&ldquo;That was back in the &rsquo;90s, that&rsquo;s a long time ago&rdquo;`)
	c.Assert(got, qt.Contains, `footnote.<sup id="fnref1:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup>`)
	c.Assert(got, qt.Contains, `<div class="footnotes" role="doc-endnotes">`)
	c.Assert(got, qt.Contains, `<dt>date</dt>`)

	toc, ok := b.(converter.TableOfContentsProvider)
	c.Assert(ok, qt.Equals, true)
	tocString := string(toc.TableOfContents().ToHTML(1, 2, false))
	c.Assert(tocString, qt.Contains, "TableOfContents")
}

func TestConvertAutoIDAsciiOnly(t *testing.T) {
	c := qt.New(t)

	content := `
## God is Good: 神真美好
`

	cfg := config.FromTOMLConfigString(`
[markup]
[markup.goldmark]
[markup.goldmark.parser]
autoHeadingIDType = 'github-ascii'

`)

	b := convert(c, testconfig.GetTestConfig(nil, cfg), content)

	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<h2 id=\"god-is-good-\">")
}

func TestConvertAutoIDBlackfriday(t *testing.T) {
	c := qt.New(t)

	content := `
## Let's try this, shall we?

`

	cfg := config.FromTOMLConfigString(`
[markup]
[markup.goldmark]
[markup.goldmark.parser]
autoHeadingIDType = 'blackfriday'
`)

	b := convert(c, testconfig.GetTestConfig(nil, cfg), content)

	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<h2 id=\"let-s-try-this-shall-we\">")
}

func TestConvertAttributes(t *testing.T) {
	c := qt.New(t)

	withBlockAttributes := func(conf *markup_config.Config) {
		conf.Goldmark.Parser.Attribute.Block = true
		conf.Goldmark.Parser.Attribute.Title = false
	}

	withTitleAndBlockAttributes := func(conf *markup_config.Config) {
		conf.Goldmark.Parser.Attribute.Block = true
		conf.Goldmark.Parser.Attribute.Title = true
	}

	for _, test := range []struct {
		name       string
		withConfig func(conf *markup_config.Config)
		input      string
		expect     any
	}{
		{
			"Title",
			nil,
			"## heading {#id .className attrName=attrValue class=\"class1 class2\"}",
			"<h2 id=\"id\" class=\"className class1 class2\" attrName=\"attrValue\">heading</h2>\n",
		},
		{
			"Blockquote",
			withBlockAttributes,
			"> foo\n> bar\n{#id .className attrName=attrValue class=\"class1 class2\"}\n",
			"<blockquote id=\"id\" class=\"className class1 class2\"><p>foo\nbar</p>\n</blockquote>\n",
		},
		/*{
			// TODO(bep) this needs an upstream fix, see https://github.com/yuin/goldmark/issues/195
			"Code block, CodeFences=false",
			func(conf *markup_config.Config) {
				withBlockAttributes(conf)
				conf.Highlight.CodeFences = false
			},
			"```bash\necho 'foo';\n```\n{.myclass}",
			"TODO",
		},*/
		{
			"Code block, CodeFences=true",
			func(conf *markup_config.Config) {
				withBlockAttributes(conf)
				conf.Highlight.CodeFences = true
			},
			"```bash {.myclass id=\"myid\"}\necho 'foo';\n````\n",
			"<div class=\"highlight myclass\" id=\"myid\"><pre style",
		},
		{
			"Code block, CodeFences=true,linenos=table",
			func(conf *markup_config.Config) {
				withBlockAttributes(conf)
				conf.Highlight.CodeFences = true
			},
			"```bash {linenos=table .myclass id=\"myid\"}\necho 'foo';\n````\n{ .adfadf }",
			[]string{
				"div class=\"highlight myclass\" id=\"myid\"><div s",
				"table style",
			},
		},
		{
			"Code block, CodeFences=true,lineanchors",
			func(conf *markup_config.Config) {
				withBlockAttributes(conf)
				conf.Highlight.CodeFences = true
				conf.Highlight.NoClasses = false
			},
			"```bash {linenos=table, anchorlinenos=true, lineanchors=org-coderef--xyz}\necho 'foo';\n```",
			"<div class=\"highlight\"><div class=\"chroma\">\n<table class=\"lntable\"><tr><td class=\"lntd\">\n<pre tabindex=\"0\" class=\"chroma\"><code><span class=\"lnt\" id=\"org-coderef--xyz-1\"><a href=\"#org-coderef--xyz-1\">1</a>\n</span></code></pre></td>\n<td class=\"lntd\">\n<pre tabindex=\"0\" class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\"><span class=\"line\"><span class=\"cl\"><span class=\"nb\">echo</span> <span class=\"s1\">&#39;foo&#39;</span><span class=\"p\">;</span>\n</span></span></code></pre></td></tr></table>\n</div>\n</div>",
		},
		{
			"Code block, CodeFences=true,lineanchors, default ordinal",
			func(conf *markup_config.Config) {
				withBlockAttributes(conf)
				conf.Highlight.CodeFences = true
				conf.Highlight.NoClasses = false
			},
			"```bash {linenos=inline, anchorlinenos=true}\necho 'foo';\nnecho 'bar';\n```\n\n```bash {linenos=inline, anchorlinenos=true}\necho 'baz';\nnecho 'qux';\n```",
			[]string{
				"<span class=\"ln\" id=\"hl-0-1\"><a class=\"lnlinks\" href=\"#hl-0-1\">1</a></span><span class=\"cl\"><span class=\"nb\">echo</span> <span class=\"s1\">&#39;foo&#39;</span>",
				"<span class=\"ln\" id=\"hl-0-2\"><a class=\"lnlinks\" href=\"#hl-0-2\">2</a></span><span class=\"cl\">necho <span class=\"s1\">&#39;bar&#39;</span>",
				"<span class=\"ln\" id=\"hl-1-2\"><a class=\"lnlinks\" href=\"#hl-1-2\">2</a></span><span class=\"cl\">necho <span class=\"s1\">&#39;qux&#39;</span>",
			},
		},
		{
			"Paragraph",
			withBlockAttributes,
			"\nHi there.\n{.myclass }",
			"<p class=\"myclass\">Hi there.</p>\n",
		},
		{
			"Ordered list",
			withBlockAttributes,
			"\n1. First\n2. Second\n{.myclass }",
			"<ol class=\"myclass\">\n<li>First</li>\n<li>Second</li>\n</ol>\n",
		},
		{
			"Unordered list",
			withBlockAttributes,
			"\n* First\n* Second\n{.myclass }",
			"<ul class=\"myclass\">\n<li>First</li>\n<li>Second</li>\n</ul>\n",
		},
		{
			"Unordered list, indented",
			withBlockAttributes,
			`* Fruit
  * Apple
  * Orange
  * Banana
  {.fruits}
* Dairy
  * Milk
  * Cheese
  {.dairies}
{.list}`,
			[]string{"<ul class=\"list\">\n<li>Fruit\n<ul class=\"fruits\">", "<li>Dairy\n<ul class=\"dairies\">"},
		},
		{
			"Table",
			withBlockAttributes,
			`| A        | B           |
| ------------- |:-------------:| -----:|
| AV      | BV |
{.myclass }`,
			"Table",
		},
		{
			"Title and Blockquote",
			withTitleAndBlockAttributes,
			"## heading {#id .className attrName=attrValue class=\"class1 class2\"}\n> foo\n> bar\n{.myclass}",
			"<h2 id=\"id\" class=\"className class1 class2\" attrName=\"attrValue\">heading</h2>\n<blockquote class=\"myclass\"><p>foo\nbar</p>\n</blockquote>\n",
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			mconf := markup_config.Default
			if test.withConfig != nil {
				test.withConfig(&mconf)
			}
			data, err := toml.Marshal(mconf)
			c.Assert(err, qt.IsNil)
			m := maps.Params{
				"markup": config.FromTOMLConfigString(string(data)).Get(""),
			}
			conf := testconfig.GetTestConfig(nil, config.NewFrom(m))
			b := convert(c, conf, test.input)
			got := string(b.Bytes())

			for _, s := range cast.ToStringSlice(test.expect) {
				c.Assert(got, qt.Contains, s)
			}
		})
	}
}

func TestConvertIssues(t *testing.T) {
	c := qt.New(t)

	// https://github.com/gohugoio/hugo/issues/7619
	c.Run("Hyphen in HTML attributes", func(c *qt.C) {
		mconf := markup_config.Default
		mconf.Goldmark.Renderer.Unsafe = true
		input := `<custom-element>
    <div>This will be "slotted" into the custom element.</div>
</custom-element>
`

		b := convert(c, unsafeConf(), input)
		got := string(b.Bytes())

		c.Assert(got, qt.Contains, "<custom-element>\n    <div>This will be \"slotted\" into the custom element.</div>\n</custom-element>\n")
	})
}

func TestCodeFence(t *testing.T) {
	c := qt.New(t)

	lines := `LINE1
LINE2
LINE3
LINE4
LINE5
`

	convertForConfig := func(c *qt.C, confStr, code, language string) string {
		cfg := config.FromTOMLConfigString(confStr)
		conf := testconfig.GetTestConfig(nil, cfg)
		pcfg := converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		}
		p, err := goldmark.Provider.New(
			pcfg,
		)

		h := highlight.New(pcfg.MarkupConfig().Highlight)

		getRenderer := func(t hooks.RendererType, id any) any {
			if t == hooks.CodeBlockRendererType {
				return h
			}
			return nil
		}

		content := "```" + language + "\n" + code + "\n```"

		c.Assert(err, qt.IsNil)
		conv, err := p.New(converter.DocumentContext{})
		c.Assert(err, qt.IsNil)
		b, err := conv.Convert(converter.RenderContext{Src: []byte(content), GetRenderer: getRenderer})
		c.Assert(err, qt.IsNil)

		return string(b.Bytes())
	}

	c.Run("Basic", func(c *qt.C) {
		confStr := `
[markup]
[markup.highlight]
noclasses=false
`

		result := convertForConfig(c, confStr, `echo "Hugo Rocks!"`, "bash")
		// TODO(bep) there is a whitespace mismatch (\n) between this and the highlight template func.
		c.Assert(result, qt.Equals, "<div class=\"highlight\"><pre tabindex=\"0\" class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\"><span class=\"line\"><span class=\"cl\"><span class=\"nb\">echo</span> <span class=\"s2\">&#34;Hugo Rocks!&#34;</span>\n</span></span></code></pre></div>")
		result = convertForConfig(c, confStr, `echo "Hugo Rocks!"`, "unknown")
		c.Assert(result, qt.Equals, "<pre tabindex=\"0\"><code class=\"language-unknown\" data-lang=\"unknown\">echo &#34;Hugo Rocks!&#34;\n</code></pre>")
	})

	c.Run("Highlight lines, default config", func(c *qt.C) {
		result := convertForConfig(c, cfgStrHighlichgtNoClasses, lines, `bash {linenos=table,hl_lines=[2 "4-5"],linenostart=3}`)
		c.Assert(result, qt.Contains, "<div class=\"highlight\"><div class=\"chroma\">\n<table class=\"lntable\"><tr><td class=\"lntd\">\n<pre tabindex=\"0\" class=\"chroma\"><code><span class")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">4")

		result = convertForConfig(c, cfgStrHighlichgtNoClasses, lines, "bash {linenos=inline,hl_lines=[2]}")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span><span class=\"cl\">LINE2\n</span></span>")
		c.Assert(result, qt.Not(qt.Contains), "<table")

		result = convertForConfig(c, cfgStrHighlichgtNoClasses, lines, "bash {linenos=true,hl_lines=[2]}")
		c.Assert(result, qt.Contains, "<table")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">2\n</span>")
	})

	c.Run("Highlight lines, linenumbers default on", func(c *qt.C) {
		confStr := `
[markup]
[markup.highlight]
noclasses=false
linenos=true
`

		result := convertForConfig(c, confStr, lines, "bash")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">2\n</span>")

		result = convertForConfig(c, confStr, lines, "bash {linenos=false,hl_lines=[2]}")
		c.Assert(result, qt.Not(qt.Contains), "class=\"lnt\"")
	})

	c.Run("Highlight lines, linenumbers default on, linenumbers in table default off", func(c *qt.C) {
		confStr := `
[markup]
[markup.highlight]
noClasses = false
lineNos = true
lineNumbersInTable = false
`

		result := convertForConfig(c, confStr, lines, "bash")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span><span class=\"cl\">LINE2\n</span>")
		result = convertForConfig(c, confStr, lines, "bash {linenos=table}")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">1\n</span>")
	})

	c.Run("No language", func(c *qt.C) {
		confStr := `
[markup]
[markup.highlight]
noClasses = false
lineNos = true
lineNumbersInTable = false
`
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		cfg.LineNumbersInTable = false

		result := convertForConfig(c, confStr, lines, "")
		c.Assert(result, qt.Contains, "<pre tabindex=\"0\"><code>LINE1\n")
	})

	c.Run("No language, guess syntax", func(c *qt.C) {
		confStr := `
[markup]
[markup.highlight]
noClasses = false
lineNos = true
lineNumbersInTable = false
guessSyntax = true
`

		result := convertForConfig(c, confStr, lines, "")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span><span class=\"cl\">LINE2\n</span></span>")
	})
}

func TestTypographerConfig(t *testing.T) {
	c := qt.New(t)

	content := `
A "quote" and 'another quote' and a "quote with a 'nested' quote" and a 'quote with a "nested" quote' and an ellipsis...
`

	confStr := `
[markup]
[markup.goldmark]
[markup.goldmark.extensions]
[markup.goldmark.extensions.typographer]
leftDoubleQuote = "&laquo;"
rightDoubleQuote = "&raquo;"
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>A &laquo;quote&raquo; and &lsquo;another quote&rsquo; and a &laquo;quote with a &rsquo;nested&rsquo; quote&raquo; and a &lsquo;quote with a &laquo;nested&raquo; quote&rsquo; and an ellipsis&hellip;</p>\n")
}

// Issue #11045
func TestTypographerImageAltText(t *testing.T) {
	c := qt.New(t)

	content := `
!["They didn't even say 'hello'!" I exclaimed.](https://example.com/image.jpg)
`

	confStr := `
[markup]
[markup.goldmark]

`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "&ldquo;They didn&rsquo;t even say &lsquo;hello&rsquo;!&rdquo; I exclaimed.")
}

func unsafeConf() config.AllProvider {
	cfg := config.FromTOMLConfigString(`
[markup]
[markup.goldmark.renderer]
unsafe = true
`)
	return testconfig.GetTestConfig(nil, cfg)
}

func TestConvertCJK(t *testing.T) {
	c := qt.New(t)

	content := `
私は太郎です。
プログラミングが好きです。\ 運動が苦手です。
`

	confStr := `
[markup]
[markup.goldmark]
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>私は太郎です。\nプログラミングが好きです。\\ 運動が苦手です。</p>\n")
}

func TestConvertCJKWithExtensionWithEastAsianLineBreaksOption(t *testing.T) {
	c := qt.New(t)

	content := `
私は太郎です。
プログラミングが好きで、
運動が苦手です。
`

	confStr := `
[markup]
[markup.goldmark]
[markup.goldmark.extensions.CJK]
enable=true
eastAsianLineBreaks=true
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>私は太郎です。プログラミングが好きで、運動が苦手です。</p>\n")
}

func TestConvertCJKWithExtensionWithEastAsianLineBreaksOptionWithSimple(t *testing.T) {
	c := qt.New(t)

	content := `
私は太郎です。
Programming が好きで、
運動が苦手です。
`

	confStr := `
[markup]
[markup.goldmark]
[markup.goldmark.extensions.CJK]
enable=true
eastAsianLineBreaks=true
eastAsianLineBreaksStyle="simple"
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>私は太郎です。\nProgramming が好きで、運動が苦手です。</p>\n")
}

func TestConvertCJKWithExtensionWithEastAsianLineBreaksOptionWithStyle(t *testing.T) {
	c := qt.New(t)

	content := `
私は太郎です。
Programming が好きで、
運動が苦手です。
`

	confStr := `
[markup]
[markup.goldmark]
[markup.goldmark.extensions.CJK]
enable=true
eastAsianLineBreaks=true
eastAsianLineBreaksStyle="css3draft"
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>私は太郎です。Programming が好きで、運動が苦手です。</p>\n")
}

func TestConvertCJKWithExtensionWithEscapedSpaceOption(t *testing.T) {
	c := qt.New(t)

	content := `
私は太郎です。
プログラミングが好きです。\ 運動が苦手です。
`

	confStr := `
[markup]
[markup.goldmark]
[markup.goldmark.extensions.CJK]
enable=true
escapedSpace=true
`

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)

	b := convert(c, conf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<p>私は太郎です。\nプログラミングが好きです。運動が苦手です。</p>\n")
}

type tableRenderer int

func (hr tableRenderer) RenderTable(cctx context.Context, w hugio.FlexiWriter, ctx hooks.TableContext) error {
	// This is set up with a render hook in the hugolib package, make it simple here.
	fmt.Fprintln(w, "Table")
	return nil
}
