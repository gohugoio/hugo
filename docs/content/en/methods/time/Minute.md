---
title: Minute
description: Returns the minute offset within the hour of the given time.Time value, in the range [0, 59].
categories: []
keywords: []
action:
  related:
    - methods/time/Year
    - methods/time/Month
    - methods/time/Day
    - methods/time/Hour
    - methods/time/Second
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Minute]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Minute }} → 44
```
