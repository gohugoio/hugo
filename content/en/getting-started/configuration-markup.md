---
title: Configure markup
description: Configure rendering of markup to HTML.
categories: [fundamentals, getting started]
keywords: [configuration,highlighting]
menu:
  docs:
    parent: getting-started
    weight: 50
weight: 50
slug: configuration-markup
toc: true
---

## Default configuration

See [Goldmark](#goldmark) for settings related to the default markdown handler in Hugo.

Below are all markup related configuration in Hugo with their default settings:

{{< code-toggle config="markup" />}}

**See each section below for details.**

## Goldmark

[Goldmark](https://github.com/yuin/goldmark/) is from Hugo 0.60 the default library used for Markdown. It's fast, it's [CommonMark](https://spec.commonmark.org/0.29/) compliant and it's very flexible.

This is the default configuration:

{{< code-toggle config="markup.goldmark" />}}

For details on the extensions, refer to [this section](https://github.com/yuin/goldmark/#built-in-extensions) of the Goldmark documentation

Some settings explained:

hardWraps
: By default, Goldmark ignores newlines within a paragraph. Set to `true` to render newlines as `<br>` elements.

unsafe
: By default, Goldmark does not render raw HTML and potentially dangerous links. If you have lots of inline HTML and/or JavaScript, you may need to turn this on.

typographer
: This extension substitutes punctuations with typographic entities like [smartypants](https://daringfireball.net/projects/smartypants/).

attribute
: Enable custom attribute support for titles and blocks by adding attribute lists inside single curly brackets (`{.myclass class="class1 class2" }`) and placing it _after the Markdown element it decorates_, on the same line for titles and on a new line directly below for blocks.

Hugo supports adding attributes (e.g. CSS classes) to Markdown blocks, e.g. tables, lists, paragraphs etc.

A blockquote with a CSS class:

```md
> foo
> bar
{.myclass}
```

There are some current limitations: For tables you can currently only apply it to the full table, and for lists the `ul`/`ol`-nodes only, e.g.:

```md
* Fruit
  * Apple
  * Orange
  * Banana
  {.fruits}
* Dairy
  * Milk
  * Cheese
  {.dairies}
{.list}
```

Note that attributes in [code fences](/content-management/syntax-highlighting/#highlighting-in-code-fences) must come after the opening tag, with any other highlighting processing instruction, e.g.:

````txt
```go {.myclass linenos=table,hl_lines=[8,"15-17"],linenostart=199}
// ... code
```
````

autoHeadingIDType ("github")
: The strategy used for creating auto IDs (anchor names). Available types are `github`, `github-ascii` and `blackfriday`. `github` produces GitHub-compatible IDs, `github-ascii` will drop any non-Ascii characters after accent normalization, and `blackfriday` will make the IDs compatible with Blackfriday, the default Markdown engine before Hugo 0.60. Note that if Goldmark is your default Markdown engine, this is also the strategy used in the [anchorize](/functions/anchorize/) template func.

## Highlight

This is the default `highlight` configuration. Note that some of these settings can be set per code block, see [Syntax Highlighting](/content-management/syntax-highlighting/).

{{< code-toggle config="markup.highlight" />}}

For `style`, see these galleries:

* [Short snippets](https://xyproto.github.io/splash/docs/all.html)
* [Long snippets](https://xyproto.github.io/splash/docs/longer/all.html)

For CSS, see [Generate Syntax Highlighter CSS](/content-management/syntax-highlighting/#generate-syntax-highlighter-css).

## Table of contents

{{< code-toggle config="markup.tableOfContents" />}}

These settings only works for the Goldmark renderer:

startLevel
: The heading level, values starting at 1 (`h1`), to start render the table of contents.

endLevel
: The heading level, inclusive, to stop render the table of contents.

ordered
: If `true`, generates an ordered list instead of an unordered list.

## Markdown render hooks

See [Markdown Render Hooks](/templates/render-hooks/).
