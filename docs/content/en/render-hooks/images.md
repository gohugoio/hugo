---
title: Image render hooks
linkTitle: Images
description: Create an image render to hook override the rendering of Markdown images to HTML.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 60
weight: 60
toc: true
---

## Markdown

A Markdown image has three components: the image description, the image destination, and optionally the image title.

```text
![white kitten](/images/kitten.jpg "A kitten!")
  ------------  ------------------  ---------
  description      destination        title
```

These components are passed into the render hook [context] as shown below.

[context]: /getting-started/glossary/#context

## Context

Image render hook templates receive the following context:

###### Attributes

(`map`) The [Markdown attributes], available if you configure your site as follows:

[Markdown attributes]: /content-management/markdown-attributes/

{{< code-toggle file=hugo >}}
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = false
[markup.goldmark.parser.attribute]
block = true
{{< /code-toggle >}}

###### Destination

(`string`) The image destination.

###### IsBlock

(`bool`) Returns true if a standalone image is not wrapped within a paragraph element.

###### Ordinal

(`int`) The zero-based ordinal of the image on the page.

###### Page

(`page`) A reference to the current page.

###### PageInner

{{< new-in 0.125.0 >}}

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### PlainText

(`string`) The image description as plain text.

###### Text

(`template.HTML`) The image description.

###### Title

(`string`) The image title.

## Examples

{{% note %}}
With inline elements such as images and links, remove leading and trailing whitespace using the `{{‑ ‑}}` delimiter notation to prevent whitespace between adjacent inline elements and text.
{{% /note %}}

In its default configuration, Hugo renders Markdown images according to the [CommonMark specification]. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

{{< code file=layouts/_default/_markup/render-image.html copy=true >}}
<img src="{{ .Destination | safeURL }}"
  {{- with .Text }} alt="{{ . }}"{{ end -}}
  {{- with .Title }} title="{{ . }}"{{ end -}}
>
{{- /* chomp trailing newline */ -}}
{{< /code >}}

To render standalone images within `figure` elements:

{{< code file=layouts/_default/_markup/render-image.html copy=true >}}
{{- if .IsBlock -}}
  <figure>
    <img src="{{ .Destination | safeURL }}"
      {{- with .Text }} alt="{{ . }}"{{ end -}}
    >
    {{- with .Title }}<figcaption>{{ . }}</figcaption>{{ end -}}
  </figure>
{{- else -}}
  <img src="{{ .Destination | safeURL }}"
    {{- with .Text }} alt="{{ . }}"{{ end -}}
    {{- with .Title }} title="{{ . }}"{{ end -}}
  >
{{- end -}}
{{< /code >}}

Note that the above requires the following site configuration:

{{< code-toggle file=hugo >}}
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = false
{{< /code-toggle >}}

## Default

{{< new-in 0.123.0 >}}

Hugo includes an [embedded image render hook] to resolve Markdown image destinations. Disabled by default, you can enable it in your site configuration:

[embedded image render hook]: {{% eturl render-image %}}

{{< code-toggle file=hugo >}}
[markup.goldmark.renderHooks.image]
enableDefault = true
{{< /code-toggle >}}

A custom render hook, even when provided by a theme or module, will override the embedded render hook regardless of the configuration setting above.

{{% note %}}
The embedded image render hook is automatically enabled for multilingual single-host sites if [duplication of shared page resources] is disabled. This is the default configuration for multilingual single-host sites.

[duplication of shared page resources]: /getting-started/configuration-markup/#duplicateresourcefiles
{{% /note %}}

The embedded image render hook resolves internal Markdown destinations by looking for a matching [page resource], falling back to a matching [global resource]. Remote destinations are passed through, and the render hook will not throw an error or warning if unable to resolve a destination.

[page resource]: /getting-started/glossary/#page-resource
[global resource]: /getting-started/glossary/#global-resource

You must place global resources in the assets directory. If you have placed your resources in the static directory, and you are unable or unwilling to move them, you must mount the static directory to the assets directory by including both of these entries in your site configuration:

{{< code-toggle file=hugo >}}
[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'static'
target = 'assets'
{{< /code-toggle >}}

Note that the embedded image render hook does not perform image processing. Its sole purpose is to resolve Markdown image destinations.

{{% include "/render-hooks/_common/pageinner.md" %}}
