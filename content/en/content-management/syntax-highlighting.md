---
title: Syntax Highlighting
description: Hugo comes with really fast syntax highlighting from Chroma.
date: 2017-02-01
publishdate: 2017-02-01
keywords: [highlighting,chroma,code blocks,syntax]
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


Hugo uses [Chroma](https://github.com/alecthomas/chroma) as its code highlighter; it is built in Go and is really, really fast -- and for the most important parts compatible with Pygments we used before.

## Configure Syntax Highlighter

See [Configure Highlight](/getting-started/configuration-markup#highlight).


## Generate Syntax Highlighter CSS

If you run with `pygmentsUseClasses=true` in your site config, you need a style sheet.

You can generate one with Hugo:

```bash
hugo gen chromastyles --style=monokai > syntax.css
```

Run `hugo gen chromastyles -h` for more options. See https://xyproto.github.io/splash/docs/ for a gallery of available styles.


## Highlight Shortcode

Highlighting is carried out via the [built-in shortcode](/content-management/shortcodes/) `highlight`. `highlight` takes exactly one required parameter for the programming language to be highlighted and requires a closing shortcode. Note that `highlight` is *not* used for client-side javascript highlighting.

Options:

* `linenos`: Valid values are `true`, `false`, `table`, `inline`. `table` will give copy-and-paste friendly code blocks) turns on line numbers.
* Setting `linenos` to `false` will turn off linenumbers if it's configured to be on in site config.{{< new-in "0.60.0" >}}
* `hl_lines` lists a set of line numbers or line number ranges to be highlighted.
* `linenostart=199` starts the line number count from 199.

### Example: Highlight Shortcode

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
// - "Chicago" (see https://www.chicagomanualofstyle.org/home.html)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
  switch strings.ToLower(style) {
  case "go":
    return strings.Title
  case "chicago":
    return transform.NewTitleConverter(transform.ChicagoStyle)
  default:
    return transform.NewTitleConverter(transform.APStyle)
  }
}
{{< / highlight >}}



## Highlight Template Func

See [Highlight](/functions/highlight/).

## Highlighting in Code Fences

Highlighting in code fences is enabled by default.{{< new-in "0.60.0" >}}

````
```go {linenos=table,hl_lines=[8,"15-17"],linenostart=199}
// ... code
```
````


Gives this:

```go {linenos=table,hl_lines=[8,"15-17"],linenostart=199}
// GetTitleFunc returns a func that can be used to transform a string to
// title case.
//
// The supported styles are
//
// - "Go" (strings.Title)
// - "AP" (see https://www.apstylebook.com/)
// - "Chicago" (see https://www.chicagomanualofstyle.org/home.html)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
  switch strings.ToLower(style) {
  case "go":
    return strings.Title
  case "chicago":
    return transform.NewTitleConverter(transform.ChicagoStyle)
  default:
    return transform.NewTitleConverter(transform.APStyle)
  }
}
```

{{< new-in "0.60.0" >}}Note that only Goldmark supports passing attributes such as `hl_lines`, and it's important that it does not contain any spaces. See [goldmark-highlighting](https://github.com/yuin/goldmark-highlighting) for more information.

The options are the same as in the [highlighting shortcode](/content-management/syntax-highlighting/#highlight-shortcode),including `linenos=false`, but note the slightly different Markdown attribute syntax.

## List of Chroma Highlighting Languages

The full list of Chroma lexers and their aliases (which is the identifier used in the `highlight` template func or when doing highlighting in code fences):

{{< chroma-lexers >}}

[Prism]: https://prismjs.com
[prismdownload]: https://prismjs.com/download.html
[Highlight.js]: https://highlightjs.org/
[Rainbow]: https://craig.is/making/rainbows
[Syntax Highlighter]: https://alexgorbatchev.com/SyntaxHighlighter/
[Google Prettify]: https://github.com/google/code-prettify
[Yandex]: https://yandex.ru/
