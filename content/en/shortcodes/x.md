---
title: X
description: Embed an X post in your content using the x shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{< new-in 0.141.0 >}}

{{% note %}}
To override Hugo's embedded `x` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl x %}}
{{% /note %}}

## Example

To display an X post with this URL:

```txt
https://x.com/SanDiegoZoo/status/1453110110599868418
```

Include this in your Markdown:

```text
{{</* x user="SanDiegoZoo" id="1453110110599868418" */>}}
```

Rendered:

{{< x user="SanDiegoZoo" id="1453110110599868418" >}}

## Privacy

Adjust the relevant privacy settings in your site configuration.

{{< code-toggle config=privacy.x />}}

disable
: (`bool`) Whether to disable the shortcode. Default is `false`.

enableDNT
: (`bool`) Whether to prevent X from using post and embedded page data for personalized suggestions and ads. Default is `false`.

simple
: (`bool`) Whether to enable simple mode. If `true`, Hugo builds a static version of the of the post without JavaScript. Default is `false`.

The source code for the simple version of the shortcode is available [here].

[here]: {{% eturl x_simple %}}

If you enable simple mode you may want to disable the hardcoded inline styles by setting `disableInlineCSS` to `true` in your site configuration. The default value for this setting is `false`.

{{< code-toggle config=services.x />}}
