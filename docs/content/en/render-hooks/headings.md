---
title: Heading render hooks
linkTitle: Headings
description: Create a heading render hook to override the rendering of Markdown headings to HTML.
categories: []
keywords: []
---

## Context

Heading render hook templates receive the following [context](g):

Anchor
: (`string`) The `id` attribute of the heading element.

Attributes
: (`map`) The [Markdown attributes], available if you configure your site as follows:

  {{< code-toggle file=hugo >}}
  [markup.goldmark.parser.attribute]
  title = true
  {{< /code-toggle >}}

Level
: (`int`) The heading level, 1 through 6.

Page
: (`page`) A reference to the current page.

PageInner
: {{< new-in 0.125.0 />}}
: (`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

PlainText
: (`string`) The heading text as plain text.

Text
: (`template.HTML`) The heading text.

[Markdown attributes]: /content-management/markdown-attributes/
[`RenderShortcodes`]: /methods/page/rendershortcodes

## Examples

In its default configuration, Hugo renders Markdown headings according to the [CommonMark specification] with the addition of automatic `id` attributes. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

```go-html-template {file="layouts/_markup/render-heading.html" copy=true}
<h{{ .Level }} id="{{ .Anchor }}" {{- with .Attributes.class }} class="{{ . }}" {{- end }}>
  {{- .Text -}}
</h{{ .Level }}>
```

To add an anchor link to the right of each heading:

```go-html-template {file="layouts/_markup/render-heading.html" copy=true}
<h{{ .Level }} id="{{ .Anchor }}" {{- with .Attributes.class }} class="{{ . }}" {{- end }}>
  {{ .Text }}
  <a href="#{{ .Anchor }}">#</a>
</h{{ .Level }}>
```

{{% include "/_common/render-hooks/pageinner.md" %}}
