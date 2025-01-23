---
title: Gist
description: Embed a GitHub Gist in your content using the gist shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
expiryDate: 2027-01-31 # Deprecated in v0.143.0. Keep this around for a couple of years because it contains instructions for creating a replacement.
---

{{% deprecated-in 0.143.0 %}}
The `gist` shortcode was deprecated in version 0.143.0 and will be removed in a future release. To continue embedding GitHub Gists in your content, you'll need to create a custom shortcode:

1. Create a new file: Create a file named `gist.html` within the `layouts/shortcodes` directory.
2. Copy the source code: Paste the [original source code] of the gist shortcode into the newly created `gist.html` file.

This will allow you to maintain the functionality of embedding GitHub Gists in your content after the deprecation of the original shortcode.

[original source code]: {{% eturl gist %}}
{{% /deprecated-in %}}

To display a GitHub gist with this URL:

```text
https://gist.github.com/user/50a7482715eac222e230d1e64dd9a89b
```

Include this in your Markdown:

```text
{{</* gist user 23932424365401ffa5e9d9810102a477 */>}}
```

To display a specific file within the gist:

```text
{{</* gist user 23932424365401ffa5e9d9810102a477 list.html */>}}
```
