---
title: collections.In
description: Reports whether a value exists within the given slice or string.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [in]
    returnType: bool
    signatures: [collections.In SLICE|STRING VALUE]
aliases: [/functions/in]
---

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ in $s "b" }} → true
```

```go-html-template
{{ $s := "abc" }}
{{ in $s "b" }} → true
```
