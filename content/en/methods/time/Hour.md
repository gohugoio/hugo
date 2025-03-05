---
title: Hour
description: Returns the hour within the day of the given time.Time value, in the range [0, 23].
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [TIME.Hour]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Hour }} → 23
```
