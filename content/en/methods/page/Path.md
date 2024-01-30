---
title: Path
description: Returns the canonical page path of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/File
    - methods/page/RelPermalink
  returnType: string
  signatures: [PAGE.Path]
---

The `Path` method on a `Page` object returns the canonical page path of the given page, regardless of whether the page is backed by a file.

```go-html-template
{{ .Path }} → /posts/post-1
```

This value is neither a file path nor a relative URL. It is a canonical identifier for each page, independent of content format, language, and URL.

{{% note %}}
Beginning with the release of [v0.92.0] in January 2022, Hugo emitted a warning whenever the `Path` method was called. The warning indicated that this method would change in a future release.

The meaning of, and value returned by, the `Path` method on a `Page` object changed with the release of [v0.123.0] in February 2024.

[v0.92.0]: https://github.com/gohugoio/hugo/releases/tag/v0.92.0
[v0.123.0]: https://github.com/gohugoio/hugo/releases/tag/v0.123.0
{{% /note %}}

To determine the canonical page path for pages backed by a file, Hugo starts with the file path, relative to the content directory, and then:

1. Strips the file extension
2. Strips the language identifier
3. Converts the result to lower case
4. Replaces all spaces with hyphens

The value returned by the `Path` method on a `Page` object is independent of content format, language, and URL modifiers such as the `slug` and `url` front matter fields.

File path|Front matter slug|Value returned by .Path
:--|:--|:--
`content/_index.md`||`/`
`content/posts/_index.md`||`/posts`
`content/posts/post-1.md`|`foo`|`/posts/post-1`
`content/posts/post-2.html`|`bar`|`/posts/post-2`

On a multilingual site, note that the return value is the same regardless of language:

File path|Front matter slug|Value returned by .Path
:--|:--|:--
`content/_index.en.md`||`/`
`content/_index.de.md`||`/`
`content/posts/_index.en.md`||`/posts`
`content/posts/_index.de.md`||`/posts`
`content/posts/posts-1.en.md`|`foo`|`/posts/post-1`
`content/posts/posts-1.de.md`|`foo`|`/posts/post-1`
`content/posts/posts-2.en.html`|`bar`|`/posts/post-2`
`content/posts/posts-2.de.html`|`bar`|`/posts/post-2`

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
