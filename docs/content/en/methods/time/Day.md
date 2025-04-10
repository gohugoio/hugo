---
title: Day
description: Returns the day of the month of the given time.Time value.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [TIME.Day]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Day }} → 27
```
