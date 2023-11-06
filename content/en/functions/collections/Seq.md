---
title: collections.Seq
description: Returns a slice of integers.
categories: []
keywords: []
action:
  aliases: [seq]
  returnType: '[]int'
  signatures:
    - collections.Seq LAST
    - collections.Seq FIRST LAST
    - collections.Seq FIRST INCREMENT LAST
related:
  - collections.Apply
  - collections.Delimit
  - collections.In
  - collections.Reverse
  - collections.Seq
  - collections.Slice
aliases: [/functions/seq]
---

```go-html-template
{{ seq 2 }} → [1 2]
{{ seq 0 2 }} → [0 1 2]
{{ seq -2 2 }} → [-2 -1 0 1 2]
{{ seq -2 2 2 }} → [-2 0 2]
```

Iterate over a sequence of integers:

```go-html-template
{{ $product := 1 }}
{{ range seq 4 }}
  {{ $product = mul $product . }}
{{ end }}
{{ $product }} → 24
```

The example above is contrived. To calculate the product of 2 or more numbers, use the [`math.Product`] function:

```go-html-template
{{ math.Product (seq 4) }} → 24
```

[`math.Product`]: /functions/math/product
