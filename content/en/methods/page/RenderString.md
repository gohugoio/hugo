---
title: RenderString
description: Renders markup to HTML.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: template.HTML
    signatures: ['PAGE.RenderString [OPTIONS] MARKUP']
aliases: [/functions/renderstring]
---

The `RenderString` method on a `Page` object renders markup to HTML.

```go-html-template
{{ $s := "An *emphasized* word" }}
{{ $s | .RenderString }} → An <em>emphasized</em> word
```

## Options

The `RenderString` method on a `Page` object accepts an options map.

`display`
: (`string`) Specify either `inline` or `block`. If `inline`, removes surrounding `p` tags from short snippets. Default is `inline`.

`markup`
: (`string`) Specify a [markup identifier][] for the provided markup. Default is the `markup` front matter value, falling back to the value derived from the page's file extension.

## Examples

Render Markdown content to HTML in block display mode:

```go-html-template
{{ $opts := dict "display" "block" }}
{{ $s | .RenderString $opts }} → <p>An <em>emphasized</em> word</p>
```

Render [Pandoc] content to HTML in block display mode:

```go-html-template
{{ $s := "H~2~O" }}

{{ $opts := dict "markup" "pandoc" "display" "block" }}
{{ $s | .RenderString $opts }} → H<sub>2</sub>O
```

[Pandoc]: /content-management/formats/#pandoc
[markup identifier]: /content-management/formats/#classification
