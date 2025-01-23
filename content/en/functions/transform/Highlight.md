---
title: transform.Highlight
description: Renders code with a syntax highlighter.
categories: []
keywords: []
action:
  aliases: [highlight]
  related:
    - functions/transform/CanHighlight
    - functions/transform/HighlightCodeBlock
  returnType: template.HTML
  signatures: ['transform.Highlight CODE LANG [OPTIONS]']
aliases: [/functions/highlight]
toc: true
---

The `highlight` function uses the [Chroma] syntax highlighter, supporting over 200 languages with more than 40 [available styles].

[chroma]: https://github.com/alecthomas/chroma
[available styles]: https://xyproto.github.io/splash/docs/

## Arguments

The `transform.Highlight` shortcode takes three arguments.

CODE
: (`string`) The code to highlight.

LANG
: (`string`) The language of the code to highlight. Choose from one of the [supported languages]. This value is case-insensitive.

OPTIONS
: (`map or string`) A map or space-separate key-value pairs wrapped in quotation marks. Set default values for each option in your [site configuration]. The key names are case-insensitive.

[site configuration]: /getting-started/configuration-markup#highlight
[supported languages]: /content-management/syntax-highlighting#list-of-chroma-highlighting-languages

## Examples

```go-html-template
{{ $input := `fmt.Println("Hello World!")` }}
{{ transform.Highlight $input "go" }}

{{ $input := `console.log('Hello World!');` }}
{{ $lang := "js" }}
{{ transform.Highlight $input $lang "lineNos=table, style=api" }}

{{ $input := `echo "Hello World!"` }}
{{ $lang := "bash" }}
{{ $opts := dict "lineNos" "table" "style" "dracula" }}
{{ transform.Highlight $input $lang $opts }}
```

## Options

{{% include "functions/_common/highlighting-options" %}}
