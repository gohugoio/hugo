---
title: transform.Highlight
description: Renders code with a syntax highlighter.
categories: []
keywords: [highlight]
params:
  functions_and_methods:
    aliases: [highlight]
    returnType: template.HTML
    signatures: ['transform.Highlight CODE [LANG] [OPTIONS]']
aliases: [/functions/highlight]
---

The `transform.Highlight` function uses the [`alecthomas/chroma`][] package to generate syntax-highlighted HTML from the provided code, [language][], and [options](#options-1).

## Arguments

`CODE`
: (`string`) The code to highlight.

`LANG`
: (`string`) The [language][] of the code to highlight. This value is case-insensitive. Optional; you can also set the language with the `type` key in OPTIONS. {{< new-in 0.162.0 />}}

`OPTIONS`
: (`map or string`) A map or comma-separated key-value pairs wrapped in quotation marks. See the [options](#options-1) below; you can set default values for each option in your [project configuration][]. The key names are case-insensitive.

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

{{ $input := `print("Hello World!")` }}
{{ $opts := dict "type" "python" "style" "dracula" }}
{{ transform.Highlight $input $opts }}
```

## Options

The `transform.Highlight` function accepts an options map.

{{% include "_common/syntax-highlighting-options.md" %}}

`code`
: {{< new-in 0.162.0 />}}
: (`string`) Overrides the `CODE` argument.

`type`
: {{< new-in 0.162.0 />}}
: (`string`) Overrides the `LANG` argument.

[`alecthomas/chroma`]: https://github.com/alecthomas/chroma
[language]: /content-management/syntax-highlighting#languages
[project configuration]: /configuration/markup#highlight
