---
title: Configure permalinks
linkTitle: Permalinks
description: Configure permalinks.
categories: []
keywords: []
---

This is the default configuration:

{{< code-toggle config=permalinks />}}

Define a URL pattern for each top-level section. Each URL pattern can target a given language and/or page kind.

> [!note]
> The [`url`] front matter field overrides any matching permalink pattern.

## Monolingual example

With this content structure:

```text
content/
├── posts/
│   ├── bash-in-slow-motion.md
│   └── tls-in-a-nutshell.md
├── tutorials/
│   ├── git-for-beginners.md
│   └── javascript-bundling-with-hugo.md
└── _index.md
```

Render tutorials under "training", and render the posts under "articles" with a date-base hierarchy:

{{< code-toggle file=hugo >}}
[permalinks.page]
posts = '/articles/:year/:month/:slug/'
tutorials = '/training/:slug/'
[permalinks.section]
posts = '/articles/'
tutorials = '/training/'
{{< /code-toggle >}}

The structure of the published site will be:

```text
public/
├── articles/
│   ├── 2023/
│   │   ├── 04/
│   │   │   └── bash-in-slow-motion/
│   │   │       └── index.html
│   │   └── 06/
│   │       └── tls-in-a-nutshell/
│   │           └── index.html
│   └── index.html
├── training/
│   ├── git-for-beginners/
│   │   └── index.html
│   ├── javascript-bundling-with-hugo/
│   │   └── index.html
│   └── index.html
└── index.html
```

To create a date-based hierarchy for regular pages in the content root:

{{< code-toggle file=hugo >}}
[permalinks.page]
'/' = '/:year/:month/:slug/'
{{< /code-toggle >}}

Use the same approach with taxonomy terms. For example, to omit the taxonomy segment of the URL:

{{< code-toggle file=hugo >}}
[permalinks.term]
'tags' = '/:slug/'
{{< /code-toggle >}}

## Multilingual example

Use the `permalinks` configuration as a component of your localization strategy.

With this content structure:

```text
content/
├── en/
│   ├── books/
│   │   ├── les-miserables.md
│   │   └── the-hunchback-of-notre-dame.md
│   └── _index.md
└── es/
    ├── books/
    │   ├── les-miserables.md
    │   └── the-hunchback-of-notre-dame.md
    └── _index.md
```

And this project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
contentDir = 'content/en'
direction = 'ltr'
label = 'English'
locale = 'en-US'
weight = 1

[languages.en.permalinks.page]
books = '/books/:slug/'

[languages.en.permalinks.section]
books = '/books/'

[languages.es]
contentDir = 'content/es'
direction = 'ltr'
label = 'Español'
locale = 'es-ES'
weight = 2

[languages.es.permalinks.page]
books = '/libros/:slug/'

[languages.es.permalinks.section]
books = '/libros/'
{{< /code-toggle >}}

The structure of the published site will be:

```text
public/
├── en/
│   ├── books/
│   │   ├── les-miserables/
│   │   │   └── index.html
│   │   ├── the-hunchback-of-notre-dame/
│   │   │   └── index.html
│   │   └── index.html
│   └── index.html
├── es/
│   ├── libros/
│   │   ├── les-miserables/
│   │   │   └── index.html
│   │   ├── the-hunchback-of-notre-dame/
│   │   │   └── index.html
│   │   └── index.html
│   └── index.html
└── index.html
```

## Tokens

Use these tokens when defining a URL pattern.

{{% include "/_common/permalink-tokens.md" %}}

[`url`]: /content-management/front-matter/#url
