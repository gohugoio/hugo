---
title: IsBranch
description: Reports whether the given page is a branch.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [PAGE.IsBranch]
---

{{< new-in 0.163.0 />}}

{{% glossary-term branch %}}

```tree
content/
├── books/
│   ├── book-1/
│   │   └── index.md    <-- kind = page      IsBranch = false
│   ├── book-2.md       <-- kind = page      IsBranch = false
│   └── _index.md       <-- kind = section   IsBranch = true
├── tags
│   ├── fiction   
│   │   └── _index.md   <-- kind = term      IsBranch = true
│   └── _index.md       <-- kind = taxonomy  IsBranch = true
└── _index.md           <-- kind = home      IsBranch = true
```

```go-html-template
{{ .IsBranch }}
```
