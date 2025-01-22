---
title: Relref
description: Insert a relative permalink to the given page reference using the relref shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
---

{{% note %}}
To override Hugo's embedded `relref` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl relref %}}
{{% /note %}}

{{% note %}}
When working with the Markdown [content format], this shortcode has become largely redundant. Its functionality is now primarily handled by [link render hooks], specifically the embedded one provided by Hugo. This hook effectively addresses all the use cases previously covered by this shortcode.

[content format]: /content-management/formats/
[link render hooks]: /render-hooks/links/
{{% /note %}}

The `relref` shortcode returns the relative permalink of the given page reference.

Example usage:

```text
[Post 1]({{%/* relref "/posts/post-1" */%}})
[Post 1]({{%/* relref "/posts/post-1.md" */%}})
[Post 1]({{%/* relref "/posts/post-1#foo" */%}})
[Post 1]({{%/* relref "/posts/post-1.md#foo" */%}})
```

Rendered:

```html
<a href="/posts/post-1/">Post 1</a>
<a href="/posts/post-1/">Post 1</a>
<a href="/posts/post-1/#foo">Post 1</a>
<a href="/posts/post-1/#foo">Post 1</a>
```

{{% note %}}
Always use the `{{%/* */%}}` notation when calling this shortcode.
{{% /note %}}
