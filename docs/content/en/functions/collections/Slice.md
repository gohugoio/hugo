---
title: collections.Slice
description: Returns a slice composed of the given values.
categories: []
keywords: []
action:
  aliases: [slice]
  related:
    - functions/collections/Dictionary
  returnType: any
  signatures: ['collections.Slice [VALUE...]']
aliases: [/functions/slice]
---

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ $s }} → [a b c]
```

To create an empty slice:

```go-html-template
{{ $s := slice }}
```
