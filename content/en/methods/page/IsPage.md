---
title: IsPage
description: Reports whether the given page is a regular page.
categories: []
keywords: []
action:
  related:
    - methods/page/IsHome
    - methods/page/IsNode
    - methods/page/IsSection
  returnType: bool
  signatures: [PAGE.IsPage]
---

The `IsPage` method on a `Page` object returns `true` if the [page kind] is `page`.

```text
content/
├── books/
│   ├── book-1/
│   │   └── index.md  <-- kind = page
│   ├── book-2.md     <-- kind = page
│   └── _index.md     <-- kind = section
└── _index.md         <-- kind = home
```

```go-html-template
{{ .IsPage }}
```

[page kind]: /getting-started/glossary/#page-kind
