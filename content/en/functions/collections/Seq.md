---
title: collections.Seq
description: Returns a slice of integers.
categories: []
keywords: []
action:
  aliases: [seq]
  related: []
  returnType: '[]int'
  signatures:
    - collections.Seq LAST
    - collections.Seq FIRST LAST
    - collections.Seq FIRST INCREMENT LAST
aliases: [/functions/seq]
---

```go-html-template
{{ seq 2 }} → [1 2]
{{ seq 0 2 }} → [0 1 2]
{{ seq -2 2 }} → [-2 -1 0 1 2]
{{ seq -2 2 2 }} → [-2 0 2]
```

A contrived example of iterating over a sequence of integers:

```go-html-template
{{ $product := 1 }}
{{ range seq 4 }}
  {{ $product = mul $product . }}
{{ end }}
{{ $product }} → 24
```

{{% note %}}
The slice created by the `seq` function is limited to 2000 elements.
{{% /note %}}
