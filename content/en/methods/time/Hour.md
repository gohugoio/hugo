---
title: Hour
description: Returns the hour within the day of the given time.Time value, in the range [0, 23].
categories: []
keywords: []
action:
  related:
    - methods/time/Year
    - methods/time/Month
    - methods/time/Day
    - methods/time/Minute
    - methods/time/Second
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Hour]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Hour }} → 23
```
