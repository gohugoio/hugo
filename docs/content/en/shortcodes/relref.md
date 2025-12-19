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
> When working with Markdown this shortcode is obsolete. Instead, to properly resolve Markdown link destinations, use the [embedded link render hook] or create your own.
>
> In its default configuration, Hugo automatically uses the embedded link render hook for multilingual single-host sites, specifically when the [duplication of shared page resources] feature is disabled. This is the default behavior for such sites. If custom link render hooks are defined by your project, modules, or themes, these will be used instead.
>
> You can also configure Hugo to `always` use the embedded link render hook, use it only as a `fallback`, or `never` use it. See&nbsp;[details](/configuration/markup/#renderhookslinkuseembedded).

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

[duplication of shared page resources]: /configuration/markup/#duplicateresourcefiles
[embedded link render hook]: /render-hooks/links/#embedded
[Markdown notation]: /content-management/shortcodes/#notation
[source code]: <{{% eturl relref %}}>
