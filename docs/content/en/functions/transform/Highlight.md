---
title: transform.Highlight
description: Renders code with a syntax highlighter.
categories: []
keywords: [highlight]
params:
  functions_and_methods:
    aliases: [highlight]
    returnType: template.HTML
    signatures: ['transform.Highlight CODE LANG [OPTIONS]']
aliases: [/functions/highlight]
---

The `transform.Highlight` function uses the [`alecthomas/chroma`][] package to generate syntax-highlighted HTML from the provided code, [language][], and [options][].

## Arguments

The `transform.Highlight` function takes three arguments.

CODE
: (`string`) The code to highlight.

LANG
: (`string`) The [language][] of the code to highlight. This value is case-insensitive.

OPTIONS
: (`map or string`) A map or comma-separated key-value pairs wrapped in quotation marks. You can set default values for each option in your [project configuration][]. The key names are case-insensitive.

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

{{% include "_common/syntax-highlighting-options.md" %}}

[`alecthomas/chroma`]: https://github.com/alecthomas/chroma
[language]: /content-management/syntax-highlighting#languages
[options]: #options-1
[project configuration]: /configuration/markup#highlight
