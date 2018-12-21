---
title: Syntax Highlighting
description: Hugo comes with really fast syntax highlighting from Chroma.
date: 2017-02-01
publishdate: 2017-02-01
keywords: [highlighting,pygments,chroma,code blocks,syntax]
categories: [content management]
menu:
  docs:
    parent: "content-management"
    weight: 300
weight: 20
sections_weight: 20
draft: false
aliases: [/extras/highlighting/,/extras/highlight/,/tools/syntax-highlighting/]
toc: true
---

From Hugo 0.28, the default syntax hightlighter in Hugo is [Chroma](https://github.com/alecthomas/chroma); it is built in Go and is really, really fast -- and for the most important parts compatible with Pygments.

If you want to continue to use Pygments (see below), set `pygmentsUseClassic=true` in your site config.

The example below shows a simple code snippet from the Hugo source highlighted with the `highlight` shortcode. Note that the gohugo.io site is generated with `pygmentsUseClasses=true` (see [Generate Syntax Highlighter CSS](#generate-syntax-highlighter-css)).

* `linenos=inline` or `linenos=table` (`table` will give copy-and-paste friendly code blocks) turns on line numbers.
* `hl_lines` lists a set of line numbers or line number ranges to be highlighted. Note that the hyphen range syntax is only supported for Chroma.
* `linenostart=199` starts the line number count from 199.

With that, this:

```
{{</* highlight go "linenos=table,hl_lines=8 15-17,linenostart=199" */>}}
// ... code
{{</* / highlight */>}}
```

Gives this:

{{< highlight go "linenos=table,hl_lines=8 15-17,linenostart=199" >}}
// GetTitleFunc returns a func that can be used to transform a string to
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
{{< / highlight >}}


## Configure Syntax Highlighter
To make the transition from Pygments to Chroma seamless, they share a common set of configuration options:

pygmentsOptions
:  A comma separated list of options. See below for a full list.

pygmentsCodefences
: Set to true to enable syntax highlighting in code fences with a language tag in markdown (see below for an example).

pygmentsStyle
: The style of code highlighting. Note that this option is not
  relevant when `pygmentsUseClasses` is set.

  Syntax highlighting galleries:
  **Chroma** ([short snippets](https://xyproto.github.io/splash/docs/all.html),
  [long snippets](https://xyproto.github.io/splash/docs/longer/all.html)),
  [Pygments](https://help.farbox.com/pygments.html)

pygmentsUseClasses
: Set to `true` to use CSS classes to format your highlighted code. See [Generate Syntax Highlighter CSS](#generate-syntax-highlighter-css).

pygmentsCodefencesGuessSyntax
: Set to `true` to try to do syntax highlighting on code fenced blocks in markdown without a language tag.

pygmentsUseClassic
: Set to true to use Pygments instead of the much faster Chroma.

### Options

`pygmentsOptions` can be set either in site config or overridden per code block in the Highlight shortcode or template func.

noclasses
: Use inline style.

linenos
: For Chroma, any value in this setting will print line numbers. Pygments has some more fine grained control.

linenostart
: Start the line numbers from this value (default is 1).


hl_lines
: Highlight a space separated list of line numbers. For Chroma, you can provide a list of ranges, i.e. "3-8 10-20".


The full set of supported options for Pygments is: `encoding`, `outencoding`, `nowrap`, `full`, `title`, `style`, `noclasses`, `classprefix`, `cssclass`, `cssstyles`, `prestyles`, `linenos`, `hl_lines`, `linenostart`, `linenostep`, `linenospecial`, `nobackground`, `lineseparator`, `lineanchors`, `linespans`, `anchorlinenos`, `startinline`. See the [Pygments HTML Formatter Documentation](http://pygments.org/docs/formatters/#HtmlFormatter) for details.


## Generate Syntax Highlighter CSS

If you run with `pygmentsUseClasses=true` in your site config, you need a style sheet.

You can generate one with Hugo:

```bash
hugo gen chromastyles --style=monokai > syntax.css
```

Run `hugo gen chromastyles -h` for more options. See https://help.farbox.com/pygments.html for a gallery of available styles.


## Highlight Shortcode

Highlighting is carried out via the [built-in shortcode](/content-management/shortcodes/) `highlight`. `highlight` takes exactly one required parameter for the programming language to be highlighted and requires a closing shortcode. Note that `highlight` is *not* used for client-side javascript highlighting.

### Example `highlight` Shortcode

{{< code file="example-highlight-shortcode-input.md" >}}
{{</* highlight html */>}}
<section id="main">
  <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Pages }}
      {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{</* /highlight */>}}
{{< /code >}}



## Highlight Template Func

See [Highlight](/functions/highlight/).

## Highlight in Code Fences

It is also possible to add syntax highlighting with GitHub flavored code fences. To enable this, set the `pygmentsCodeFences` to `true` in Hugo's [configuration file](/getting-started/configuration/);

````
```go-html-template
<section id="main">
  <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Pages }}
      {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
```
````

## List of Chroma Highlighting Languages

The full list of Chroma lexers and their aliases (which is the identifier used in the `highlight` template func or when doing highlighting in code fences):

{{< chroma-lexers >}}

## Highlight with Pygments Classic

If you for some reason don't want to use the built-in Chroma highlighter, you can set `pygmentsUseClassic=true` in your config and add Pygments to your path.

{{% note "Disclaimers on Pygments" %}}
* Pygments is relatively slow and _causes a performance hit when building your site_, but Hugo has been designed to cache the results to disk.
* The caching can be turned off by setting the `--ignoreCache` flag to `true`.
* The languages available for highlighting depend on your Pygments installation.
{{% /note %}}

If you have never worked with Pygments before, here is a brief primer:

+ Install Python from [python.org](https://www.python.org/downloads/). Version 2.7.x is already sufficient.
+ Run `pip install Pygments` in order to install Pygments. Once installed, Pygments gives you a command `pygmentize`. Make sure it sits in your PATH; otherwise, Hugo will not be able to find and use it.

On Debian and Ubuntu systems, you may also install Pygments by running `sudo apt-get install python3-pygments`.



[Prism]: http://prismjs.com
[prismdownload]: http://prismjs.com/download.html
[Highlight.js]: http://highlightjs.org/
[Rainbow]: http://craig.is/making/rainbows
[Syntax Highlighter]: http://alexgorbatchev.com/SyntaxHighlighter/
[Google Prettify]: https://github.com/google/code-prettify
[Yandex]: http://yandex.ru/
