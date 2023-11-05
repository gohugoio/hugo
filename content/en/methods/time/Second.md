---
title: Second
description: Returns the second offset within the minute of the given time.Time value, in the range [0, 59].
categories: []
keywords: []
action:
  related:
    - methods/time/Year
    - methods/time/Month
    - methods/time/Day
    - methods/time/Hour
    - methods/time/Minute
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Second]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Second }} → 58
```
