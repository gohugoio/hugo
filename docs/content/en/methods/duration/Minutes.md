---
title: Minutes
description: Returns the time.Duration value as a floating point number of minutes.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: float64
    signatures: [DURATION.Minutes]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}
{{ $d.Minutes }} → 212.525
```
