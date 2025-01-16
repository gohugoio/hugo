---
title: Instagram
description: Embed an Instagram post in your content using the instagram shortcode.
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
To override Hugo's embedded `instagram` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl instagram %}}
{{% /note %}}

## Example

To display an Instagram post with this URL:

```text
https://www.instagram.com/p/CxOWiQNP2MO/
```

Include this in your Markdown:

```text
{{</* instagram CxOWiQNP2MO */>}}
```

Huge renders this to:

{{< instagram CxOWiQNP2MO >}}

## Privacy

Adjust the relevant privacy settings in your site configuration.

{{< code-toggle config=privacy.instagram />}}

disable
: (`bool`) Whether to disable the shortcode. Default is `false`.

simple
: (`bool`) Whether to enable simple mode for image card generation. If `true`, Hugo creates a static card without JavaScript. This mode only supports image cards, and the image is fetched directly from Instagram's servers. Default is `false`.

The source code for the simple version of the shortcode is available [here].

If you enable simple mode you may want to disable the hardcoded inline styles by setting `disableInlineCSS` to `true` in your site configuration. The default value for this setting is `false`.

[here]: {{% eturl instagram_simple %}}

{{< code-toggle config=services.instagram />}}
