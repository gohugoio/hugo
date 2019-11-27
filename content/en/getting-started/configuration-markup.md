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