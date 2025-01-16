---
title: Ref
description: Insert a permalink to the given page reference using the ref shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
---

{{% note %}}
To override Hugo's embedded `ref` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl ref %}}
{{% /note %}}

{{% note %}}
When working with the Markdown [content format], this shortcode has become largely redundant. Its functionality is now primarily handled by [link render hooks], specifically the embedded one provided by Hugo. This hook effectively addresses all the use cases previously covered by this shortcode.

[content format]: /content-management/formats/
[link render hooks]: /render-hooks/images/#default
{{% /note %}}

The `ref` shortcode returns the permalink of the given page reference.

Example usage:

```text
[Post 1]({{%/* ref "/posts/post-1" */%}})
[Post 1]({{%/* ref "/posts/post-1.md" */%}})
[Post 1]({{%/* ref "/posts/post-1#foo" */%}})
[Post 1]({{%/* ref "/posts/post-1.md#foo" */%}})
```

Rendered:

```html
<a href="https://example.org/posts/post-1/">Post 1</a>
<a href="https://example.org/posts/post-1/">Post 1</a>
<a href="https://example.org/posts/post-1/#foo">Post 1</a>
<a href="https://example.org/posts/post-1/#foo">Post 1</a>
```

{{% note %}}
Always use the `{{%/* */%}}` notation when calling this shortcode.
{{% /note %}}
