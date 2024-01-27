---
title: Day
description: Returns the day of the month of the given time.Time value.
categories: []
keywords: []
action:
  related:
    - methods/time/Year
    - methods/time/Month
    - methods/time/Hour
    - methods/time/Minute
    - methods/time/Second
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Day]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Day }} → 27
```
