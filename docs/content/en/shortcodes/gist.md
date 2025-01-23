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
---

{{% note %}}
To override Hugo's embedded `gist` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl gist %}}
{{% /note %}}

To display a GitHub gist with this URL:

```text
https://gist.github.com/user/50a7482715eac222e230d1e64dd9a89b
```

Include this in your Markdown:

```text
{{</* gist user 23932424365401ffa5e9d9810102a477 */>}}
```

This will display all files in the gist alphabetically by file name.

{{< gist jmooring 23932424365401ffa5e9d9810102a477 >}}

To display a specific file within the gist:

```text
{{</* gist user 23932424365401ffa5e9d9810102a477 list.html */>}}
```

{{< gist jmooring 23932424365401ffa5e9d9810102a477 list.html >}}
