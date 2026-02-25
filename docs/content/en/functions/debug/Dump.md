---
title: debug.Dump
description: Returns an object dump as a string.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [debug.Dump VALUE]
---

```go-html-template
<pre>{{ debug.Dump hugo.Data.books }}</pre>
```

```json
[
  {
    "author": "Victor Hugo",
    "rating": 4,
    "title": "The Hunchback of Notre Dame"
  },
  {
    "author": "Victor Hugo",
    "rating": 5,
    "title": "Les Misérables"
  }
]
```

> [!note]
> Output from this function may change from one release to the next. Use for debugging only.
