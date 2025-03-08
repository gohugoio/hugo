---
title: Ref shortcode
linkTitle: Ref
description: Insert a permalink to the given page reference using the ref shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `ref` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

> [!note]
> When working with the Markdown [content format], this shortcode has become largely redundant. Its functionality is now primarily handled by [link render hooks], specifically the embedded one provided by Hugo. This hook effectively addresses all the use cases previously covered by this shortcode.

## Usage

The `ref` shortcode accepts either a single positional argument (the path) or one or more named arguments, as listed below.

## Arguments

{{% include "_common/ref-and-relref-options.md" %}}

## Examples

The `ref` shortcode typically provides the destination for a Markdown link.

> [!note]
> Always use [Markdown notation] notation when calling this shortcode.

The following examples show the rendered output for a page on the English version of the site:

```md
[Link A]({{%/* ref "/books/book-1" */%}})

[Link B]({{%/* ref path="/books/book-1" */%}})

[Link C]({{%/* ref path="/books/book-1" lang="de" */%}})

[Link D]({{%/* ref path="/books/book-1" lang="de" outputFormat="json" */%}})
```

Rendered:

```html
<a href="https://example.org/en/books/book-1/">Link A</a>

<a href="https://example.org/en/books/book-1/">Link B</a>

<a href="https://example.org/de/books/book-1/">Link C</a>

<a href="https://example.org/de/books/book-1/index.json">Link D</a>
```

## Error handling

{{% include "_common/ref-and-relref-error-handling.md" %}}

[content format]: /content-management/formats/
[link render hooks]: /render-hooks/images/#default
[Markdown notation]: /content-management/shortcodes/#notation
[source code]: {{% eturl ref %}}
