---
title: Instagram shortcode
linkTitle: Instagram
description: Embed an Instagram post in your content using the instagram shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `instagram` shortcode, copy the [source code] to a file with the same name in the `layouts/_shortcodes` directory.

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

[source code]: <{{% eturl instagram %}}>
