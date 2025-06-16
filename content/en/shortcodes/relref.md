---
title: Relref shortcode
linkTitle: Relref
description: Insert a relative permalink to the given page reference using the relref shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `relref` shortcode, copy the [source code] to a file with the same name in the `layouts/_shortcodes` directory.

> [!note]
> When working with Markdown, this shortcode is obsolete. Instead, use a [link render hook] that resolves the link destination using the `GetPage` method on the `Page` object. You can either create your own, or simply enable the [embedded link render hook]. The embedded link render hook is automatically enabled for multilingual single-host projects.

## Usage

The `relref` shortcode accepts either a single positional argument (the path) or one or more named arguments, as listed below.

## Arguments

{{% include "_common/ref-and-relref-options.md" %}}

## Examples

The `relref` shortcode typically provides the destination for a Markdown link.

> [!note]
> Always use [Markdown notation] notation when calling this shortcode.

The following examples show the rendered output for a page on the English version of the site:

```md
[Link A]({{%/* relref "/books/book-1" */%}})

[Link B]({{%/* relref path="/books/book-1" */%}})

[Link C]({{%/* relref path="/books/book-1" lang="de" */%}})

[Link D]({{%/* relref path="/books/book-1" lang="de" outputFormat="json" */%}})
```

Rendered:

```html
<a href="/en/books/book-1/">Link A</a>

<a href="/en/books/book-1/">Link B</a>

<a href="/de/books/book-1/">Link C</a>

<a href="/de/books/book-1/index.json">Link D</a>
```

## Error handling

{{% include "_common/ref-and-relref-error-handling.md" %}}

[content format]: /content-management/formats/
[embedded link render hook]: /render-hooks/links/#default
[link render hook]: /render-hooks/links/
[Markdown notation]: /content-management/shortcodes/#notation
[source code]: {{% eturl relref %}}
