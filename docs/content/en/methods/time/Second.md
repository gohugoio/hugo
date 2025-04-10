---
title: Second
description: Returns the second offset within the minute of the given time.Time value, in the range [0, 59].
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [TIME.Second]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Second }} → 58
```
