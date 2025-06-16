---
title: transform.ToMath
description: Renders mathematical equations and expressions written in the LaTeX markup language.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: template.HTML
    signatures: ['transform.ToMath INPUT [OPTIONS]']
aliases: [/functions/tomath]
---

{{< new-in 0.132.0 />}}

Hugo uses an embedded instance of the [KaTeX] display engine to render mathematical markup to HTML. You do not need to install the KaTeX display engine.

```go-html-template
{{ transform.ToMath "c = \\pm\\sqrt{a^2 + b^2}" }}
```

> [!note]
> By default, Hugo renders mathematical markup to [MathML], and does not require any CSS to display the result.
>
> To optimize rendering quality and accessibility, use the `htmlAndMathml` output option as described below. This approach requires an external stylesheet.

```go-html-template
{{ $opts := dict "output" "htmlAndMathml" }}
{{ transform.ToMath "c = \\pm\\sqrt{a^2 + b^2}" $opts }}
```

## Options

Pass a map of options as the second argument to the `transform.ToMath` function. The options below are a subset of the KaTeX [rendering options].

displayMode
: (`bool`) Whether to render in display mode instead of inline mode. Default is `false`.

errorColor
: (`string`) The color of the error messages expressed as an RGB [hexadecimal color]. Default is `#cc0000`.

fleqn
: (`bool`) Whether to render flush left with a 2em left margin. Default is `false`.

macros
: (`map`) A map of macros to be used in the math expression. Default is `{}`.

    ```go-html-template
    {{ $macros := dict
      "\\addBar" "\\bar{#1}"
      "\\bold" "\\mathbf{#1}"
    }}
    {{ $opts := dict "macros" $macros }}
    {{ transform.ToMath "\\addBar{y} + \\bold{H}" $opts }}
    ```

minRuleThickness
: (`float`) The minimum thickness of the fraction lines in `em`. Default is `0.04`.

output
: (`string`) Determines the markup language of the output, one of `html`, `mathml`, or `htmlAndMathml`. Default is `mathml`.

    With `html` and `htmlAndMathml` you must include the KaTeX style sheet within the `head` element of your base template.

    ```html
    <link href="https://cdn.jsdelivr.net/npm/katex@0.16.22/dist/katex.min.css" rel="stylesheet">

strict
: {{< new-in 0.147.6 />}}
: (`string`) Controls how KaTeX handles LaTeX features that offer convenience but aren't officially supported, one of `error`, `ignore`, or `warn`. Default is `error`.

  - `error`: Throws an error when convenient, unsupported LaTeX features are encountered.
  - `ignore`: Allows convenient, unsupported LaTeX features without any feedback.
  - `warn`: {{< new-in 0.147.7 />}} Emits a warning when convenient, unsupported LaTeX features are encountered.

: The `newLineInDisplayMode` error code, which flags the use of `\\`
or `\newline` in display mode outside an array or tabular environment, is
intentionally designed not to throw an error, despite this behavior
being questionable.

throwOnError
: (`bool`) Whether to throw a `ParseError` when KaTeX encounters an unsupported command or invalid LaTeX. Default is `true`.

## Error handling

There are three ways to handle errors:

1. Let KaTeX throw an error and fail the build. This is the default behavior.
1. Set the `throwOnError` option to `false` to make KaTeX render the expression as an error instead of throwing an error. See [options](#options).
1. Handle the error in your template.

The example below demonstrates error handing within a template.

## Example

Instead of client-side JavaScript rendering of mathematical markup using MathJax or KaTeX, create a passthrough render hook which calls the `transform.ToMath` function.

### Step 1

Enable and configure the Goldmark [passthrough extension] in your site configuration. The passthrough extension preserves raw Markdown within delimited snippets of text, including the delimiters themselves.

{{< code-toggle file=hugo copy=true >}}
[markup.goldmark.extensions.passthrough]
enable = true

[markup.goldmark.extensions.passthrough.delimiters]
block = [['\[', '\]'], ['$$', '$$']]
inline = [['\(', '\)']]
{{< /code-toggle >}}

> [!note]
> The configuration above precludes the use of the `$...$` delimiter pair for inline equations. Although you can add this delimiter pair to the configuration, you must double-escape the `$` symbol when used outside of math contexts to avoid unintended formatting.

### Step 2

Create a [passthrough render hook] to capture and render the LaTeX markup.

```go-html-template {file="layouts/_markup/render-passthrough.html" copy=true}
{{- $opts := dict "output" "htmlAndMathml" "displayMode" (eq .Type "block") }}
{{- with try (transform.ToMath .Inner $opts) }}
  {{- with .Err }}
    {{- errorf "Unable to render mathematical markup to HTML using the transform.ToMath function. The KaTeX display engine threw the following error: %s: see %s." . $.Position }}
  {{- else }}
    {{- .Value }}
    {{- $.Page.Store.Set "hasMath" true }}
  {{- end }}
{{- end -}}
```

### Step 3

In your base template, conditionally include the KaTeX CSS within the head element.

```go-html-template {file="layouts/baseof.html" copy=true}
<head>
  {{ $noop := .WordCount }}
  {{ if .Page.Store.Get "hasMath" }}
    <link href="https://cdn.jsdelivr.net/npm/katex@0.16.22/dist/katex.min.css" rel="stylesheet">
  {{ end }}
</head>
```

In the above, note the use of a [noop](g) statement to force content rendering before we check the value of `hasMath` with the `Store.Get` method.

### Step 4

Add some mathematical markup to your content, then test.

```text {file="content/example.md"}
This is an inline \(a^*=x-b^*\) equation.

These are block equations:

\[a^*=x-b^*\]

$$a^*=x-b^*$$
```

## Chemistry

{{< new-in 0.144.0 />}}

You can also use the `transform.ToMath` function to render chemical equations, leveraging the `\ce` and `\pu` functions from the [mhchem] package.

```text
$$C_p[\ce{H2O(l)}] = \pu{75.3 J // mol K}$$
```

$$C_p[\ce{H2O(l)}] = \pu{75.3 J // mol K}$$

[hexadecimal color]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color
[KaTeX]: https://katex.org/
[MathML]: https://developer.mozilla.org/en-US/docs/Web/MathML
[mhchem]: https://mhchem.github.io/MathJax-mhchem/
[passthrough extension]: /configuration/markup/#passthrough
[passthrough render hook]: /render-hooks/passthrough/
[rendering options]: https://katex.org/docs/options.html
