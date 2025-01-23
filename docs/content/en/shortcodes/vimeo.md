---
title: Vimeo
description: Embed a Vimeo video in your content using the vimeo shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{% note %}}
To override Hugo's embedded `vimeo` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl vimeo %}}
{{% /note %}}

## Example

To display a Vimeo video with this URL:

```text
https://vimeo.com/channels/staffpicks/55073825
```

Include this in your Markdown:

```text
{{</* vimeo 55073825 */>}}
```

Hugo renders this to:

{{< vimeo 55073825 >}}

## Parameters

class
: (`string`) The `class` attribute of the wrapping `div` element. Adding one or more CSS classes disables inline styling.

id
: (`string`) The `id` of the Vimeo video

title
: (`string`) The `title` attribute of the `iframe` element.

If you proivde a `class` or `title` you must use a named parameter for the `id`.

```text
{{</* vimeo id=55073825 class="foo bar" title="My Video" */>}}
```

## Privacy

Adjust the relevant privacy settings in your site configuration.

{{< code-toggle config=privacy.vimeo />}}

disable
: (`bool`) Whether to disable the shortcode. Default is `false`.

enableDNT
: (`bool`) Whether to block the Vimeo player from tracking session data and analytics. Default is `false`.

simple
: (`bool`) Whether to enable simple mode. If `true`, the video thumbnail is fetched from Vimeo and overlaid with a play button. Clicking the thumbnail opens the video in a new Vimeo tab. Default is `false`.

The source code for the simple version of the shortcode is available [here].

[here]: {{% eturl vimeo_simple %}}
