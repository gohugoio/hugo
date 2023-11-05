---
title: collections.Slice
description: Creates a slice (array) of all passed arguments.
categories: []
keywords: []
action:
  aliases: [slice]
  returnType: any
  signatures: [collections.Slice ITEM...]
related:
  - collections.Append
  - collections.Apply
  - collections.Delimit
  - collections.In
  - collections.Reverse
  - collections.Seq
  - collections.Slice
aliases: [/functions/slice]
---

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ $s }} → [a b c]
```
