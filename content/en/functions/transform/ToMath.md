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
: A map of zero or more [options].

## Options

These are a sub set of the [KaTeX options].

output
: String. Default is `mathml`.\
`html` Outputs HTML only.\
`mathml`: Outputs MathML only.\
`htmlAndMathml`: Outputs HTML for visual rendering and MathML for accessibility.

displayMode
: Boolean. Default is `false`.\
If `true` the math will be rendered in display mode. If false the math will be rendered in `inline` mode.

leqno
: Boolean. Default is `false`.\
If `true` the math will be rendered with the equation numbers on the left.

fleqn
: Boolean. Default is `false`.\
If `true`, render flush left with a 2em left margin.

minRuleThickness
: Float. Default is `0.04`.\
The minimum thickness of the fraction lines in `em`.

macros
: Map. Default is `{}`.\
A map of macros to be used in the math expression.

throwOnError
: Boolean. Default is `true`.\
If `true`, KaTeX will throw a `ParseError` when it encounters an unsupported command or invalid LaTex. See [error handling].

errorColor
: String. Default is `#cc0000`.\
The color of the error messages.\
A color string given in the format "#XXX" or "#XXXXXX"


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



[options]: #options
[error handling]: #error-handling
[KaTeX options]: https://katex.org/docs/options.html
