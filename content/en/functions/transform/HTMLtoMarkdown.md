---
title: transform.HTMLToMarkdown
description: Converts HTML to Markdown.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [transform.HTMLToMarkdown INPUT]
---

{{< new-in "0.151.0" />}}

> [!note]
> This function is experimental and its API may change in the future.

The `transform.HTMLToMarkdown` function converts HTML to Markdown by utilizing the [`html-to-markdown`][] Go package.

## Usage

```go-html-template
{{ .Content | transform.HTMLToMarkdown | safeHTML }}
```

## Plugins

The conversion process is enabled by the following `html-to-markdown` plugins:

Plugin|Description
:--|:--
Base|Implements basic shared functionality
CommonMark|Implements Markdown according to the [Commonmark Spec][]
Table|Implements tables according to the [GitHub Flavored Markdown Spec][]

[`html-to-markdown`]: https://github.com/JohannesKaufmann/html-to-markdown?tab=readme-ov-file#readme
[Commonmark Spec]: https://spec.commonmark.org/current/
[GitHub Flavored Markdown Spec]: https://github.github.com/gfm/
