---
title: Comment
description: Include hidden comments in your content with the comment shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    identifier: shortcodes-comment
    parent: shortcodes
    weight:
weight:
expiryDate: 2025-01-22 # with v0.142.0 and later use HTML comments instead
---

{{% note %}}
To override Hugo's embedded `comment` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl comment %}}
{{% /note %}}

{{< new-in "0.137.1" >}}

Use the `comment` shortcode to include comments in your content. Hugo will ignore the text within these comments when rendering your site.

Use it inline:

```text
{{%/* comment */%}} rewrite the paragraph below {{%/* /comment */%}}
```

Or as a block comment:

```text
{{%/* comment */%}}
rewrite the paragraph below
{{%/* /comment */%}}
```

Although you can call this shortcode using the `{{</* */>}}` notation, computationally it is more efficient to call it using the `{{%/* */%}}` notation as shown above.
