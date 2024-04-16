---
title: Heading render hooks
linkTitle: Headings
description: Create a heading render hook to override the rendering of Markdown headings to HTML.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 40
weight: 40
toc: true
---

## Context

Heading render hook templates receive the following [context]:

[context]: /getting-started/glossary/#context

###### Anchor

(`string`) The `id` attribute of the heading element.

###### Attributes

(`map`) The Markdown attributes, available if you configure your site as follows:

{{< code-toggle file=hugo >}}
[markup.goldmark.parser.attribute]
title = true
{{< /code-toggle >}}

###### Level

(`int`) The heading level, 1 through 6.

###### Page

(`page`) A reference to the current page.

###### PageInner

{{< new-in 0.125.0 >}}

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### PlainText

(`string`) The heading text as plain text.

###### Text

(`string`) The heading text.

## Examples

In its default configuration, Hugo renders Markdown headings according to the [CommonMark specification] with the addition of automatic `id` attributes. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

{{< code file=layouts/_default/_markup/render-heading.html copy=true >}}
<h{{ .Level }} id="{{ .Anchor }}">
  {{- .Text | safeHTML -}}
</h{{ .Level }}>
{{< /code >}}

To add an anchor link to the right of each heading:

{{< code file=layouts/_default/_markup/render-heading.html copy=true >}}
<h{{ .Level }} id="{{ .Anchor }}">
  {{ .Text | safeHTML }}
  <a href="#{{ .Anchor }}">#</a>
</h{{ .Level }}>
{{< /code >}}

{{% include "/render-hooks/_common/pageinner.md" %}}
