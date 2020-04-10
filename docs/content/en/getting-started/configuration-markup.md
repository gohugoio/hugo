---
title: Configure Markup
description: How to handle Markdown and other markup related configuration.
date: 2019-11-15
categories: [getting started,fundamentals]
keywords: [configuration,highlighting]
weight: 65
sections_weight: 65
slug: configuration-markup
toc: true
---

## Configure Markup

{{< new-in "0.60.0" >}}

See [Goldmark](#goldmark) for settings related to the default Markdown handler in Hugo.

Below are all markup related configuration in Hugo with their default settings:

{{< code-toggle config="markup" />}}

**See each section below for details.**

### Goldmark

[Goldmark](https://github.com/yuin/goldmark/) is from Hugo 0.60 the default library used for Markdown. It's fast, it's [CommonMark](https://spec.commonmark.org/0.29/) compliant and it's very flexible. Note that the feature set of Goldmark vs Blackfriday isn't the same; you gain a lot but also lose some, but we will work to bridge any gap in the upcoming Hugo versions.

This is the default configuration:

{{< code-toggle config="markup.goldmark" />}}

Some settings explained:

unsafe
: By default, Goldmark does not render raw HTMLs and potentially dangerous links. If you have lots of inline HTML and/or JavaScript, you may need to turn this on.

typographer
: This extension substitutes punctuations with typographic entities like [smartypants](https://daringfireball.net/projects/smartypants/).

autoHeadingIDType ("github") {{< new-in "0.62.2" >}}
: The strategy used for creating auto IDs (anchor names). Available types are `github`, `github-ascii` and `blackfriday`. `github` produces GitHub-compatible IDs, `github-ascii` will drop any non-Ascii characters after accent normalization, and `blackfriday` will make the IDs work as with [Blackfriday](#blackfriday), the default Markdown engine before Hugo 0.60. Note that if Goldmark is your default Markdown engine, this is also the strategy used in the [anchorize](/functions/anchorize/) template func.

### Blackfriday


[Blackfriday](https://github.com/russross/blackfriday) was Hugo's default Markdown rendering engine, now replaced with Goldmark. But you can still use it: Just set `defaultMarkdownHandler` to `blackfriday` in your top level `markup` config.

This is the default config:

{{< code-toggle config="markup.blackFriday" />}}

### Highlight

This is the default `highlight` configuration. Note that some of these settings can be set per code block, see [Syntax Highlighting](/content-management/syntax-highlighting/).

{{< code-toggle config="markup.highlight" />}}

For `style`, see these galleries:

* [Short snippets](https://xyproto.github.io/splash/docs/all.html)
* [Long snippets](https://xyproto.github.io/splash/docs/longer/all.html)

For CSS, see [Generate Syntax Highlighter CSS](/content-management/syntax-highlighting/#generate-syntax-highlighter-css).

### Table Of Contents

{{< code-toggle config="markup.tableOfContents" />}}

These settings only works for the Goldmark renderer:

startLevel
: The heading level, values starting at 1 (`h1`), to start render the table of contents.

endLevel
: The heading level, inclusive, to stop render the table of contents.

ordered
: Whether or not to generate an ordered list instead of an unordered list.


## Markdown Render Hooks

{{< new-in "0.62.0" >}}

Note that this is only supported with the [Goldmark](#goldmark) renderer.

These Render Hooks allow custom templates to render links and images from markdown.

You can do this by creating templates with base names `render-link` and/or `render-image` inside `layouts/_default/_markup`.

You can define [Output-Format-](/templates/output-formats) and [language-](/content-management/multilingual/)specific templates if needed.[^hooktemplate] Your `layouts` folder may look like this:

```bash
layouts
└── _default
    └── _markup
        ├── render-image.html
        ├── render-image.rss.xml
        └── render-link.html
```

Some use cases for the above:

* Resolve link references using `.GetPage`. This would make links portable as you could translate `./my-post.md` (and similar constructs that would work on GitHub) into `/blog/2019/01/01/my-post/` etc.
* Add `target=_blank` to external links.
* Resolve and [process](/content-management/image-processing/) images.

### Render Hook Templates

Both `render-link` and `render-image` templates will receive this context:

Page
: The [Page](/variables/page/) being rendered.

Destination
: The URL.

Title
: The title attribute.

Text
: The rendered (HTML) link text.

PlainText
: The plain variant of the above.

#### Link with title Markdown example :

```md
[Text](https://www.gohugo.io "Title")
```

Here is a code example for how the render-link.html template could look:

{{< code file="layouts/_default/_markup/render-link.html" >}}
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}{{ if strings.HasPrefix .Destination "http" }} target="_blank"{{ end }}>{{ .Text | safeHTML }}</a>
{{< /code >}}

#### Image Markdown example:

```md
![Text](https://d33wubrfki0l68.cloudfront.net/c38c7334cc3f23585738e40334284fddcaf03d5e/2e17c/images/hugo-logo-wide.svg "Title")
```

Here is a code example for how the render-image.html template could look:

{{< code file="layouts/_default/_markup/render-image.html" >}}
<p class="md__image">
  <img src="{{ .Destination | safeURL }}" alt="{{ .Text }}" {{ with .Title}} title="{{ . }}"{{ end }} />
</p>
{{< /code >}}

[^hooktemplate]: It's currently only possible to have one set of render hook templates, e.g. not per `Type` or `Section`. We may consider that in a future version.

