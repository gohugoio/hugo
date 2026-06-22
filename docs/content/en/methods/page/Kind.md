---
title: Kind
description: Returns the kind of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [PAGE.Kind]
---

The [page kind](g) is one of `home`, `page`, `section`, `taxonomy`, or `term`.

```tree
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
