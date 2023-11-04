---
title: collections.Slice
linkTitle: slice
description: Creates a slice (array) of all passed arguments.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [slice]
  returnType: any
  signatures: [collections.Slice ITEM...]
relatedFunctions:
  - collections.Append
  - collections.Apply
  - collections.Delimit
  - collections.In
  - collections.Reverse
  - collections.Seq
  - collections.Slice
aliases: [/functions/slice]
---

One use case is the concatenation of elements in combination with the [`delimit` function]:

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ $s }} → [a b c]
```
