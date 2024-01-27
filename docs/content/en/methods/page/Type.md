---
title: Type
description: Returns the content type of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Kind
    - methods/page/Layout
    - methods/page/Type
  returnType: string
  signatures: [PAGE.Type]
---

The `Type` method on a `Page` object returns the [content type] of the given page. The content type is defined by the `type` field in front matter,  or inferred from the top-level directory name if the `type` field in front matter is not defined.

With this content structure:

```text
content/
├── auction/
│   ├── _index.md
│   ├── item-1.md
│   └── item-2.md  <-- front matter: type = books
├── books/
│   ├── _index.md
│   ├── book-1.md
│   └── book-2.md
├── films/
│   ├── _index.md
│   ├── film-1.md 
│   └── film-2.md
└── _index.md
```

To list the books, regardless of [section]:

```go-html-template
{{ range where .Site.RegularPages.ByTitle "Type" "books" }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```

Hugo renders this to;

```html
<h2><a href="/books/book-1/">Book 1</a></h2>
<h2><a href="/books/book-2/">Book 2</a></h2>
<h2><a href="/auction/item-2/">Item 2</a></h2>
```

The `type` field in front matter is also useful for targeting a template. See&nbsp;[details].

[content type]: /getting-started/glossary/#content-type
[details]: /templates/lookup-order/#target-a-template
[section]: /getting-started/glossary/#section
