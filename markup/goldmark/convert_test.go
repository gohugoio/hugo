// Copyright 2019 The Hugo Authors. All rights reserved.
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

package goldmark

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"

	"github.com/gohugoio/hugo/markup/highlight"

	"github.com/gohugoio/hugo/markup/markup_config"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

func convert(c *qt.C, mconf markup_config.Config, content string) converter.Result {

	p, err := Provider.New(
		converter.ProviderConfig{
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{DocumentID: "thedoc"})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{RenderTOC: true, Src: []byte(content)})
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
	mconf := markup_config.Default
	mconf.Highlight.NoClasses = false
	mconf.Goldmark.Renderer.Unsafe = true

	b := convert(c, mconf, content)
	got := string(b.Bytes())

	// Links
	//	c.Assert(got, qt.Contains, `<a href="https://docuapi.netlify.com/">Live Demo here!</a>`)

	// Header IDs
	c.Assert(got, qt.Contains, `<h2 id="custom">Custom ID</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="auto-id">Auto ID</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好">神真美好</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好-1">神真美好</h2>`, qt.Commentf(got))
	c.Assert(got, qt.Contains, `<h2 id="神真美好-2">神真美好</h2>`, qt.Commentf(got))

	// Code fences
	c.Assert(got, qt.Contains, "<div class=\"highlight\"><pre class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\">LINE1\n</code></pre></div>")
	c.Assert(got, qt.Contains, "Code Fences No Lexer</h2>\n<pre><code class=\"language-moo\" data-lang=\"moo\">LINE1\n</code></pre>")

	// Extensions
	c.Assert(got, qt.Contains, `Autolink: <a href="https://gohugo.io/">https://gohugo.io/</a>`)
	c.Assert(got, qt.Contains, `Strikethrough:<del>Hi</del> Hello, world`)
	c.Assert(got, qt.Contains, `<th>foo</th>`)
	c.Assert(got, qt.Contains, `<li><input disabled="" type="checkbox"> Push my commits to GitHub</li>`)

	c.Assert(got, qt.Contains, `Straight double &ldquo;quotes&rdquo; and single &lsquo;quotes&rsquo;`)
	c.Assert(got, qt.Contains, `Dashes (“&ndash;” and “&mdash;”) `)
	c.Assert(got, qt.Contains, `Three consecutive dots (“&hellip;”)`)
	c.Assert(got, qt.Contains, `&ldquo;That was back in the &rsquo;90s, that&rsquo;s a long time ago&rdquo;`)
	c.Assert(got, qt.Contains, `footnote.<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup>`)
	c.Assert(got, qt.Contains, `<section class="footnotes" role="doc-endnotes">`)
	c.Assert(got, qt.Contains, `<dt>date</dt>`)

	toc, ok := b.(converter.TableOfContentsProvider)
	c.Assert(ok, qt.Equals, true)
	tocHTML := toc.TableOfContents().ToHTML(1, 2, false)
	c.Assert(tocHTML, qt.Contains, "TableOfContents")

}

func TestConvertAutoIDAsciiOnly(t *testing.T) {
	c := qt.New(t)

	content := `
## God is Good: 神真美好
`
	mconf := markup_config.Default
	mconf.Goldmark.Parser.AutoHeadingIDType = goldmark_config.AutoHeadingIDTypeGitHubAscii
	b := convert(c, mconf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<h2 id=\"god-is-good-\">")
}

func TestConvertAutoIDBlackfriday(t *testing.T) {
	c := qt.New(t)

	content := `
## Let's try this, shall we?

`
	mconf := markup_config.Default
	mconf.Goldmark.Parser.AutoHeadingIDType = goldmark_config.AutoHeadingIDTypeBlackfriday
	b := convert(c, mconf, content)
	got := string(b.Bytes())

	c.Assert(got, qt.Contains, "<h2 id=\"let-s-try-this-shall-we\">")
}

func TestCodeFence(t *testing.T) {
	c := qt.New(t)

	lines := `LINE1
LINE2
LINE3
LINE4
LINE5
`

	convertForConfig := func(c *qt.C, conf highlight.Config, code, language string) string {
		mconf := markup_config.Default
		mconf.Highlight = conf

		p, err := Provider.New(
			converter.ProviderConfig{
				MarkupConfig: mconf,
				Logger:       loggers.NewErrorLogger(),
			},
		)

		content := "```" + language + "\n" + code + "\n```"

		c.Assert(err, qt.IsNil)
		conv, err := p.New(converter.DocumentContext{})
		c.Assert(err, qt.IsNil)
		b, err := conv.Convert(converter.RenderContext{Src: []byte(content)})
		c.Assert(err, qt.IsNil)

		return string(b.Bytes())
	}

	c.Run("Basic", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false

		result := convertForConfig(c, cfg, `echo "Hugo Rocks!"`, "bash")
		// TODO(bep) there is a whitespace mismatch (\n) between this and the highlight template func.
		c.Assert(result, qt.Equals, `<div class="highlight"><pre class="chroma"><code class="language-bash" data-lang="bash"><span class="nb">echo</span> <span class="s2">&#34;Hugo Rocks!&#34;</span>
</code></pre></div>`)
		result = convertForConfig(c, cfg, `echo "Hugo Rocks!"`, "unknown")
		c.Assert(result, qt.Equals, "<pre><code class=\"language-unknown\" data-lang=\"unknown\">echo &quot;Hugo Rocks!&quot;\n</code></pre>")

	})

	c.Run("Highlight lines, default config", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false

		result := convertForConfig(c, cfg, lines, `bash {linenos=table,hl_lines=[2 "4-5"],linenostart=3}`)
		c.Assert(result, qt.Contains, "<div class=\"highlight\"><div class=\"chroma\">\n<table class=\"lntable\"><tr><td class=\"lntd\">\n<pre class=\"chroma\"><code><span class")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">4")

		result = convertForConfig(c, cfg, lines, "bash {linenos=inline,hl_lines=[2]}")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n</span>")
		c.Assert(result, qt.Not(qt.Contains), "<table")

		result = convertForConfig(c, cfg, lines, "bash {linenos=true,hl_lines=[2]}")
		c.Assert(result, qt.Contains, "<table")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">2\n</span>")
	})

	c.Run("Highlight lines, linenumbers default on", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true

		result := convertForConfig(c, cfg, lines, "bash")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">2\n</span>")

		result = convertForConfig(c, cfg, lines, "bash {linenos=false,hl_lines=[2]}")
		c.Assert(result, qt.Not(qt.Contains), "class=\"lnt\"")
	})

	c.Run("Highlight lines, linenumbers default on, linenumbers in table default off", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		cfg.LineNumbersInTable = false

		result := convertForConfig(c, cfg, lines, "bash")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n<")
		result = convertForConfig(c, cfg, lines, "bash {linenos=table}")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">1\n</span>")
	})

	c.Run("No language", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		cfg.LineNumbersInTable = false

		result := convertForConfig(c, cfg, lines, "")
		c.Assert(result, qt.Contains, "<pre><code>LINE1\n")
	})

	c.Run("No language, guess syntax", func(c *qt.C) {
		cfg := highlight.DefaultConfig
		cfg.NoClasses = false
		cfg.GuessSyntax = true
		cfg.LineNos = true
		cfg.LineNumbersInTable = false

		result := convertForConfig(c, cfg, lines, "")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n<")
	})
}
