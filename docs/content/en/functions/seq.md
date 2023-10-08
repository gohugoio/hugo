---
title: seq
description: Returns a slice of integers.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
signature: ["seq LAST", "seq FIRST LAST", "seq FIRST INCREMENT LAST"]
relatedfuncs: []
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
