---
title: transform.HighlightCodeBlock
description: Highlights code received in context within a code block render hook.
categories: []
keywords: [highlight]
params:
  functions_and_methods:
    aliases: []
    returnType: highlight.HighlightResult 
    signatures: ['transform.HighlightCodeBlock CONTEXT [OPTIONS]']
---

The `transform.HighlightCodeBlock` function uses the [`alecthomas/chroma`][] package to generate syntax-highlighted HTML from code received in context within a code block render hook. This function is only useful within a code block render hook.

## Arguments

CONTEXT
: The [context][] passed into a code block render hook.

OPTIONS
: (`map`) A map of key-value pairs. See the [options](#options-1) below. The key names are case-insensitive.

## Return value

`transform.HighlightCodeBlock` returns a `HighlightResult` object with two methods.

`Wrapped`
: (`template.HTML`) Returns highlighted code wrapped in `<div>`, `<pre>`, and `<code>` elements. This is identical to the value returned by the `transform.Highlight` function.

`Inner`
: (`template.HTML`) Returns highlighted code without any wrapping elements, allowing you to create your own wrapper.

## Examples

```go-html-template
{{ $result := transform.HighlightCodeBlock . }}
{{ $result.Wrapped }}
```

To override the default options:

```go-html-template
{{ $opts := merge .Options (dict "lineNos" true) }}
{{ $result := transform.HighlightCodeBlock . $opts }}
{{ $result.Wrapped }}
```

To fall back to plain text when the language is not supported by the highlighter:

```go-html-template
{{ $opts := dict }}
{{ if not (transform.CanHighlight .Type) }}
  {{ $opts = dict "type" "text" }}
{{ end }}
{{ $result := transform.HighlightCodeBlock . $opts }}
{{ $result.Wrapped }}
```

## Options

The `transform.HighlightCodeBlock` function accepts an options map.

{{% include "_common/syntax-highlighting-options.md" %}}

`code`
: {{< new-in 0.162.0 />}}
: (`string`) Overrides the code received from the code block context.

`type`
: {{< new-in 0.162.0 />}}
: (`string`) Overrides the language received from the code block context.

[`alecthomas/chroma`]: https://github.com/alecthomas/chroma
[context]: /render-hooks/code-blocks/#context
