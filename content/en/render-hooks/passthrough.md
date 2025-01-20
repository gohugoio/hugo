---
title: Passthrough render hooks
linkTitle: Passthrough
description: Create a passthrough render hook to override the rendering of text snippets captured by the Goldmark passthrough extension.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 80
weight: 80
toc: true
---

{{< new-in 0.132.0 >}}

## Overview

Hugo uses [Goldmark] to render Markdown to HTML. Goldmark supports custom extensions to extend its core functionality. The Goldmark [passthrough extension] captures and preserves raw Markdown within delimited snippets of text, including the delimiters themselves. These are known as _passthrough elements_.

[Goldmark]: https://github.com/yuin/goldmark
[passthrough extension]: /getting-started/configuration-markup/#passthrough

Depending on your choice of delimiters, Hugo will classify a passthrough element as either _block_ or _inline_. Consider this contrived example:

{{< code file=content/sample.md >}}
This is a

\[block\]

passthrough element with opening and closing block delimiters.

This is an \(inline\) passthrough element with opening and closing inline delimiters.
{{< /code >}}

Update your site configuration to enable the passthrough extension and define  opening and closing delimiters for each passthrough element type, either `block` or `inline`. For example:

{{< code-toggle file=hugo >}}
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
block = [['\[', '\]'], ['$$', '$$']]
inline = [['\(', '\)']]
{{< /code-toggle >}}

In the example above there are two sets of `block` delimiters. You may use either one in your Markdown.

The Goldmark passthrough extension is often used in conjunction with the MathJax or KaTeX display engine to render [mathematical expressions] written in the LaTeX markup language.

[mathematical expressions]: /content-management/mathematics/

To enable custom rendering of passthrough elements, create a passthrough render hook.

## Context

Passthrough render hook templates receive the following [context]:

[context]: /getting-started/glossary/#context

###### Attributes

(`map`) The [Markdown attributes], available if you configure your site as follows:

[Markdown attributes]: /content-management/markdown-attributes/

{{< code-toggle file=hugo >}}
[markup.goldmark.parser.attribute]
block = true
{{< /code-toggle >}}

Hugo populates the `Attributes` map for _block_ passthrough elements. Markdown attributes are not applicable to _inline_ elements.

###### Inner
(`string`) The inner content of the passthrough element, excluding the delimiters.

###### Ordinal

(`int`) The zero-based ordinal of the passthrough element on the page.

###### Page

(`page`) A reference to the current page.

###### PageInner

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### Position

(`string`) The position of the passthrough element within the page content.

###### Type

(`string`) The passthrough element type, either `block` or `inline`.

## Example

Instead of client-side JavaScript rendering of mathematical markup using MathJax or KaTeX, create a passthrough render hook which calls the [`transform.ToMath`] function.

[`transform.ToMath`]: /functions/transform/tomath/

{{< code file=layouts/_default/_markup/render-passthrough.html copy=true >}}
{{- $opts := dict "output" "htmlAndMathml" "displayMode" (eq .Type "block") }}
{{- with try (transform.ToMath .Inner $opts) }}
  {{- with .Err }}
    {{ errorf "Unable to render mathematical markup to HTML using the transform.ToMath function. The KaTeX display engine threw the following error: %s: see %s." . $.Position }}
  {{- else }}
    {{- .Value }}
    {{- $.Page.Store.Set "hasMath" true }}
  {{- end }}
{{- end -}}
{{< /code >}}

Then, in your base template, conditionally include the KaTeX CSS within the head element:

{{< code file=layouts/_default/baseof.html copy=true >}}
<head>
  {{ $noop := .WordCount }}
  {{ if .Page.Store.Get "hasMath" }}
    <link href="https://cdn.jsdelivr.net/npm/katex@0.16.21/dist/katex.min.css" rel="stylesheet">
  {{ end }}
</head>
{{< /code >}}

In the above, note the use of a [noop](g) statement to force content rendering before we check the value of `hasMath` with the `Store.Get` method.

Although you can use one template with conditional logic as shown above, you can also create separate templates for each [`Type`](#type) of passthrough element:

```text
layouts/
└── _default/
    └── _markup/
        ├── render-passthrough-block.html
        └── render-passthrough-inline.html
```

{{% include "/render-hooks/_common/pageinner.md" %}}
