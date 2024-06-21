---
title: Path
description: Returns the logical path of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/File
    - methods/page/RelPermalink
  returnType: string
  signatures: [PAGE.Path]
toc: true
---

{{< new-in 0.123.0 >}}

The `Path` method on a `Page` object returns the logical path of the given page, regardless of whether the page is backed by a file.

[logical path]: /getting-started/glossary#logical-path

```go-html-template
{{ .Path }} → /posts/post-1
```

This value is neither a file path nor a relative URL. It is a logical identifier for each page, independent of content format, language, and URL modifiers.

{{% note %}}
Beginning with the release of [v0.92.0] in January 2022, Hugo emitted a warning whenever calling the `Path` method. The warning indicated that this method would change in a future release.

The meaning of, and value returned by, the `Path` method on a `Page` object changed with the release of [v0.123.0] in February 2024.

[v0.92.0]: https://github.com/gohugoio/hugo/releases/tag/v0.92.0
[v0.123.0]: https://github.com/gohugoio/hugo/releases/tag/v0.123.0
{{% /note %}}

To determine the logical path for pages backed by a file, Hugo starts with the file path, relative to the content directory, and then:

1. Strips the file extension
2. Strips the language identifier
3. Converts the result to lower case
4. Replaces spaces with hyphens

The value returned by the `Path` method on a `Page` object is independent of content format, language, and URL modifiers such as the `slug` and `url` front matter fields.

## Examples

### Monolingual site

Note that the logical path is independent of content format and URL modifiers.

File path|Front matter slug|Logical path
:--|:--|:--
`content/_index.md`||`/`
`content/posts/_index.md`||`/posts`
`content/posts/post-1.md`|`foo`|`/posts/post-1`
`content/posts/post-2.html`|`bar`|`/posts/post-2`

### Multilingual site

Note that the logical path is independent of content format, language identifiers, and URL modifiers.

File path|Front matter slug|Logical path
:--|:--|:--
`content/_index.en.md`||`/`
`content/_index.de.md`||`/`
`content/posts/_index.en.md`||`/posts`
`content/posts/_index.de.md`||`/posts`
`content/posts/posts-1.en.md`|`foo`|`/posts/post-1`
`content/posts/posts-1.de.md`|`foo`|`/posts/post-1`
`content/posts/posts-2.en.html`|`bar`|`/posts/post-2`
`content/posts/posts-2.de.html`|`bar`|`/posts/post-2`

### Pages not backed by a file

The `Path` method on a `Page` object returns a value regardless of whether the page is backed by a file.

```text
content/
└── posts/
    └── post-1.md  <-- front matter: tags = ['hugo']
```

When you build the site:

```text
public/
├── posts/
│   ├── post-1/
│   │   └── index.html    .Page.Path = /posts/post-1
│   └── index.html        .Page.Path = /posts
├── tags/
│   ├── hugo/
│   │   └── index.html    .Page.Path = /tags/hugo
│   └── index.html        .Page.Path = /tags
└── index.html            .Page.Path = /
```

## Finding pages

These methods, functions, and shortcodes use the logical path to find the given page:

Methods|Functions|Shortcodes
:--|:--|:--
[`Site.GetPage`]|[`urls.Ref`]|[`ref`]
[`Page.GetPage`]|[`urls.RelRef`]|[`relref`]
[`Page.Ref`]||
[`Page.RelRef`]||
[`Shortcode.Ref`]||
[`Shortcode.RelRef`]||

[`urls.Ref`]: /functions/urls/ref/
[`urls.RelRef`]: /functions/urls/relref/
[`Page.GetPage`]: /methods/page/getpage/
[`Site.GetPage`]: /methods/site/getpage/
[`ref`]: /content-management/shortcodes/#ref
[`relref`]: /content-management/shortcodes/#relref
[`Page.Ref`]: /methods/page/ref/
[`Page.RelRef`]: /methods/page/relref/
[`Shortcode.Ref`]: /methods/shortcode/ref
[`Shortcode.RelRef`]: /methods/shortcode/relref

{{% note %}}
Specify the logical path when using any of these methods, functions, or shortcodes. If you include a file extension or language identifier, Hugo will strip these values before finding the page in the logical tree.
{{% /note %}}


## Logical tree

Just as file paths form a file tree, logical paths form a logical tree.

A file tree:

```text
content/
└── s1/
    ├── p1/
    │   └── index.md 
    └── p2.md
```

The same content represented as a logical tree:

```text
content/
└── s1/
    ├── p1
    └── p2 
```

A key difference between these trees is the relative path from p1 to p2:

- In the file tree, the relative path from p1 to p2 is `../p2.md`
- In the logical tree, the relative path is `p2`

{{% note %}}
Remember to use the logical path when using any of the methods, functions, or shortcodes listed in the previous section. If you include a file extension or language identifier, Hugo will strip these values before finding the page in the logical tree.
{{% /note %}}
