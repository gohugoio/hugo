---
title: Blockquote render hooks
linkTitle: Blockquotes
description: Create a blockquote render hook to override the rendering of Markdown blockquotes to HTML.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 30
weight: 30
toc: true
---

{{< new-in 0.132.0 >}}

## Context

Blockquote render hook templates receive the following [context]:

[context]: /getting-started/glossary/#context

###### AlertType

(`string`) Applicable when [`Type`](#type) is `alert`, this is the alert type converted to lowercase. See the [alerts](#alerts) section below.

###### AlertTitle

{{< new-in 0.134.0 >}}

(`hstring.HTML`) Applicable when [`Type`](#type) is `alert` when using [Obsidian callouts] syntax, this is the alert title converted to HTML. 

###### AlertSign

{{< new-in 0.134.0 >}}

(`string`) Applicable when [`Type`](#type) is `alert` when using [Obsidian callouts] syntax, this is one of "+", "-" or "" (empty string) to indicate the presence of a foldable sign.

[Obsidian callouts]: https://help.obsidian.md/Editing+and+formatting/Callouts

###### Attributes

(`map`) The [Markdown attributes], available if you configure your site as follows:

[Markdown attributes]: /content-management/markdown-attributes/

{{< code-toggle file=hugo >}}
[markup.goldmark.parser.attribute]
block = true
{{< /code-toggle >}}

###### Ordinal

(`int`) The zero-based ordinal of the blockquote on the page.

###### Page

(`page`) A reference to the current page.

###### PageInner

(`page`) A reference to a page nested via the [`RenderShortcodes`] method.

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### Position

(`string`) The position of the blockquote within the page content.

###### Text
(`string`) The blockquote text, excluding the alert designator if present. See the [alerts](#alerts) section below.

###### Type

(`bool`) The blockquote type. Returns `alert` if the blockquote has an alert designator, else `regular`. See the [alerts](#alerts) section below.

## Examples

In its default configuration, Hugo renders Markdown blockquotes according to the [CommonMark specification]. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

{{< code file=layouts/_default/_markup/render-blockquote.html copy=true >}}
<blockquote>
  {{ .Text | safeHTML }}
</blockquote>
{{< /code >}}

To render a blockquote as an HTML `figure` element with an optional citation and caption:

{{< code file=layouts/_default/_markup/render-blockquote.html copy=true >}}
<figure>
  <blockquote {{ with .Attributes.cite }}cite="{{ . }}"{{ end }}>
    {{ .Text | safeHTML }}
  </blockquote>
  {{ with .Attributes.caption }}
    <figcaption class="blockquote-caption">
      {{ . | safeHTML }}
    </figcaption>
  {{ end }}
</figure>
{{< /code >}}

Then in your markdown:

```text
> Some text
{cite="https://gohugo.io" caption="Some caption"}
```

## Alerts

Also known as _callouts_ or _admonitions_, alerts are blockquotes used to emphasize critical information. For example:

{{< code file=content/example.md lang=text >}}
> [!NOTE]
> Useful information that users should know, even when skimming content.

> [!TIP]
> Helpful advice for doing things better or more easily.

> [!IMPORTANT]
> Key information users need to know to achieve their goal.

> [!WARNING]
> Urgent info that needs immediate user attention to avoid problems.

> [!CAUTION]
> Advises about risks or negative outcomes of certain actions.
{{< /code >}}


{{% note %}}
This syntax is compatible with both the GitHub Alert Markdown extension and Obsidian's callout syntax. 
But note that GitHub will not recognize callouts with one of Obsidian's extensions (e.g. callout title or the foldable sign).
{{% /note %}}

The first line of each alert is an alert designator consisting of an exclamation point followed by the alert type, wrapped within brackets.

The blockquote render hook below renders a multilingual alert if an alert designator is present, otherwise it renders a blockquote according to the CommonMark specification.

{{< code file=layouts/_default/_markup/render-blockquote.html copy=true >}}
{{ $emojis := dict
  "caution" ":exclamation:"
  "important" ":information_source:"
  "note" ":information_source:"
  "tip" ":bulb:"
  "warning" ":information_source:"
}}

{{ if eq .Type "alert" }}
  <blockquote class="alert alert-{{ .AlertType }}">
    <p class="alert-heading">
      {{ transform.Emojify (index $emojis .AlertType) }}
      {{ or (i18n .AlertType) (title .AlertType) }}
    </p>
    {{ .Text | safeHTML }}
  </blockquote>
{{ else }}
  <blockquote>
    {{ .Text | safeHTML }}
  </blockquote>
{{ end }}
{{< /code >}}

To override the label, create these entries in your i18n files:

{{< code-toggle file=i18n/en.toml >}}
caution = 'Caution'
important = 'Important'
note = 'Note'
tip = 'Tip'
warning = 'Warning'
{{< /code-toggle >}}


Although you can use one template with conditional logic as shown above, you can also create separate templates for each [`Type`](#type) of blockquote:

```text
layouts/
└── _default/
    └── _markup/
        ├── render-blockquote-alert.html
        └── render-blockquote-regular.html
```
