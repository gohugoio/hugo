---
title: Sections
description: Returns a collection of first level section pages.
categories: []
keywords: []
action:
  related:
    - methods/site/AllPages
    - methods/site/Pages
    - methods/site/RegularPages
  returnType: page.Pages
  signatures: [SITE.Sections]
---

Given this content structure:

```text
content/
├── books/
│   ├── book-1.md
│   └── book-2.md
├── films/
│   ├── film-1.md
│   └── film-2.md
└── _index.md
```

This template:

```go-html-template
{{ range .Site.Sections }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

Is rendered to:

```html
<h2><a href="/books/">Books</a></h2>
<h2><a href="/films/">Films</a></h2>
```
