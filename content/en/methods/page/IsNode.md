---
title: IsNode
description: Reports whether the given page is a node.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [PAGE.IsNode]
---

The `IsNode` method on a `Page` object checks if the [page kind](g) is one of the following: `home`, `section`, `taxonomy`, or `term`. If it is, the method returns `true`, indicating the page is a [node](g). Otherwise, if the page kind is page, it returns `false`.

```text
content/
├── books/
│   ├── book-1/
│   │   └── index.md    <-- kind = page      IsNode = false
│   ├── book-2.md       <-- kind = page      IsNode = false
│   └── _index.md       <-- kind = section   IsNode = true
├── tags
│   ├── fiction   
│   │   └── _index.md   <-- kind = term      IsNode = true
│   └── _index.md       <-- kind = taxonomy  IsNode = true
└── _index.md           <-- kind = home      IsNode = true
```

```go-html-template
{{ .IsNode }}
```
