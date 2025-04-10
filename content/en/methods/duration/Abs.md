---
title: Abs
description: Returns the absolute value of the given time.Duration value.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: time.Duration
    signatures: [DURATION.Abs]
---

```go-html-template
{{ $d = time.ParseDuration "-3h" }}
{{ $d.Abs }} → 3h0m0s
```
