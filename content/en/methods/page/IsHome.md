---
title: IsHome
description: Reports whether the given page is the home page.
categories: []
keywords: []
action:
  related:
    - methods/page/IsNode
    - methods/page/IsPage
    - methods/page/IsSection
  returnType: bool
  signatures: [PAGE.IsHome]
---

The `IsHome` method on a `Page` object returns `true` if the [page kind] is `home`.

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
{{ .IsHome }}
```

[page kind]: /getting-started/glossary/#page-kind
