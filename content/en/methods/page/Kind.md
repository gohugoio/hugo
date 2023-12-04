---
title: Kind
description: Returns the kind of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Type
  returnType: string
  signatures: [PAGE.Kind]
---

The [page kind] is one of `home`, `page`, `section`, `taxonomy`, or `term`.

```text
content/
├── books/
│   ├── book-1/
│   │   └── index.md    <-- kind = page
│   ├── book-2.md       <-- kind = page
│   └── _index.md       <-- kind = section
├── tags/
│   ├── fiction/
│   │   └── _index.md   <-- kind = term
│   └── _index.md       <-- kind = taxonomy
└── _index.md           <-- kind = home
```

To get the value within a template:

```go-html-template
{{ .Kind }}
```

[page kind]: /getting-started/glossary/#page-kind
