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

(`template.HTML`) Applicable when [`Type`](#type) is `alert`, this is the alert title. See the [alerts](#alerts) section below.

###### AlertSign

{{< new-in 0.134.0 >}}

(`string`) Applicable when [`Type`](#type) is `alert`, this is the alert sign. Typically used to indicate whether an alert is graphically foldable, this is one of&nbsp;`+`,&nbsp;`-`,&nbsp;or an empty string. See the [alerts](#alerts) section below.

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

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### Position

(`string`) The position of the blockquote within the page content.

###### Text
(`template.HTML`) The blockquote text, excluding the first line if [`Type`](#type) is `alert`. See the [alerts](#alerts) section below.

###### Type

(`bool`) The blockquote type. Returns `alert` if the blockquote has an alert designator, else `regular`. See the [alerts](#alerts) section below.

## Examples

In its default configuration, Hugo renders Markdown blockquotes according to the [CommonMark specification]. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

{{< code file=layouts/_default/_markup/render-blockquote.html copy=true >}}
<blockquote>
  {{ .Text }}
</blockquote>
{{< /code >}}

To render a blockquote as an HTML `figure` element with an optional citation and caption:

{{< code file=layouts/_default/_markup/render-blockquote.html copy=true >}}
<figure>
  <blockquote {{ with .Attributes.cite }}cite="{{ . }}"{{ end }}>
    {{ .Text }}
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

Also known as _callouts_ or _admonitions_, alerts are blockquotes used to emphasize critical information.

### Basic syntax

With the basic Markdown syntax, the first line of each alert is an alert designator consisting of an exclamation point followed by the alert type, wrapped within brackets. For example:

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

The basic syntax is compatible with [GitHub], [Obsidian], and [Typora].

[GitHub]: https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#alerts
[Obsidian]: https://help.obsidian.md/Editing+and+formatting/Callouts
[Typora]: https://support.typora.io/Markdown-Reference/#callouts--github-style-alerts

### Extended syntax

With the extended Markdown syntax, you may optionally include an alert sign and/or an alert title. The alert sign is one of&nbsp;`+`&nbsp;or&nbsp;`-`, typically used to indicate whether an alert is graphically foldable. For example:

{{< code file=content/example.md lang=text >}}
> [!WARNING]+ Radiation hazard
> Do not approach or handle without protective gear.
{{< /code >}}

The extended syntax is compatible with [Obsidian].

{{% note %}}
The extended syntax is not compatible with GitHub or Typora. If you include an alert sign or an alert title, these applications render the Markdown as a blockquote.
{{% /note %}}

### Example

This blockquote render hook renders a multilingual alert if an alert designator is present, otherwise it renders a blockquote according to the CommonMark specification.

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
      {{ with .AlertTitle }}
        {{ . }}
      {{ else }}
        {{ or (i18n .AlertType) (title .AlertType) }}
      {{ end }}
    </p>
    {{ .Text }}
  </blockquote>
{{ else }}
  <blockquote>
    {{ .Text }}
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

{{% include "/render-hooks/_common/pageinner.md" %}}
