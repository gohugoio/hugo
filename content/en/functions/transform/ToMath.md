---
title: transform.ToMath
description: Renders a math expression using KaTeX.
categories: []
keywords: []
action:
  aliases: []
  related:
    - content-management/mathematics
  returnType: types.Result[template.HTML]
  signatures: ['transform.ToMath EXPRESSION [OPTIONS]']
aliases: [/functions/tomath]
toc: true
---

{{< new-in "0.132.0" >}}

{{% note %}}
This feature was introduced in Hugo 0.132.0 and is marked as experimental.

This does not mean that it's going to be removed, but this is our first use of WASI/Wasm in Hugo, and we need to see how it [works in the wild](https://github.com/gohugoio/hugo/issues/12736) before we can set it in stone.
{{% /note %}}

## Arguments

EXPRESSION
: The math expression to render using KaTeX.

OPTIONS
: A map of zero or more options.

## Options

These are a subset of the [KaTeX options].

output
: (`string`). Determines the markup language of the output. One of `html`, `mathml`, or `htmlAndMathml`. Default is `mathml`.

    <!-- Indent to prevent spliting the description list. -->

    With `html` and `htmlAndMathml` you must include KaTeX CSS within the `head` element of your base template. For example:

    ```html
    <head>
      ...
      <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css" integrity="sha384-nB0miv6/jRmo5UMMR1wu3Gz6NLsoTkbqJghGIsx//Rlm+ZU03BU6SQNC66uf4l5+" crossorigin="anonymous">
      ...
    </head>
    ```

displayMode
: (`bool`) If `true` render in display mode, else render in inline mode. Default is `false`.

leqno
: (`bool`) If `true` render with the equation numbers on the left. Default is `false`.

fleqn
: (`bool`) If `true` render flush left with a 2em left margin. Default is `false`.

minRuleThickness
: (`float`) The minimum thickness of the fraction lines in `em`. Default is `0.04`.

macros
: (`map`) A map of macros to be used in the math expression. Default is `{}`.

throwOnError
: (`bool`) If `true` throw a `ParseError` when KaTeX encounters an unsupported command or invalid LaTex. See [error handling]. Default is `true`.

errorColor
: (`string`) The color of the error messages expressed as an RGB [hexadecimal color]. Default is `#cc0000`.

## Examples

### Basic

```go-html-template
{{ transform.ToMath "c = \\pm\\sqrt{a^2 + b^2}" }}
```

### Macros

```go-html-template
{{ $macros := dict 
    "\\addBar" "\\bar{#1}"
    "\\bold" "\\mathbf{#1}"
}}
{{ $opts := dict "macros" $macros }}
{{ transform.ToMath "\\addBar{y} + \\bold{H}" $opts }}
```

## Error handling

There are 3 ways to handle errors from KaTeX:

1. Let KaTeX throw an error and make the build fail. This is the default behavior.
1. Handle the error in your template. See the render hook example below.
1. Set the `throwOnError` option to `false` to make KaTeX render the expression as an error instead of throwing an error. See [options].

{{< code file=layouts/_default/_markup/render-passthrough-inline.html copy=true >}}
{{ with transform.ToMath .Inner }}
  {{ with .Err }}
    {{ errorf "Failed to render KaTeX: %q. See %s" . $.Position }}
  {{ else }}
    {{ . }}
  {{ end }}
{{ end }}
{{- /* chomp trailing newline */ -}}
{{< /code >}}

[error handling]: #error-handling
[KaTeX options]: https://katex.org/docs/options.html
[hexadecimal color]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color
