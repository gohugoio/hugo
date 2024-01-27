---
title: Nanosecond
description: Returns the nanosecond offset within the second of the given time.Time value, in the range [0, 999999999].
categories: []
keywords: []
action:
  related:
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Nanosecond]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Nanosecond }} → 0
```
