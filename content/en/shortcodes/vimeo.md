---
title: Vimeo shortcode
linkTitle: Vimeo
description: Embed a Vimeo video in your content using the vimeo shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `vimeo` shortcode, copy the [source code] to a file with the same name in the `layouts/_shortcodes` directory.

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

## Arguments

id
: (string) The video `id`. Optional if the `id` is the first and only positional argument.

allowFullScreen
: {{< new-in 0.146.0 />}}
: (`bool`) Whether the `iframe` element can activate full screen mode. Default is `true`.

class
: (`string`) The `class` attribute of the wrapping `div` element. Adding one or more CSS classes disables inline styling.

loading
: {{< new-in 0.146.0 />}}
: (`string`) The loading attribute of the `iframe` element, either `eager` or `lazy`. Default is `eager`.

title
: (`string`) The `title` attribute of the `iframe` element.

Here's an example using some of the available arguments:

```text
{{</* vimeo id=55073825 allowFullScreen=false loading=lazy */>}}
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

The source code for the simple version of the shortcode is available [in this file].

[in this file]: <{{% eturl vimeo_simple %}}>
[source code]: <{{% eturl vimeo %}}>
