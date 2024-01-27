---
title: Year
description: Returns the year of the given time.Time value.
categories: []
keywords: []
action:
  related:
    - methods/time/Month
    - methods/time/Day
    - methods/time/Hour
    - methods/time/Minute
    - methods/time/Second
    - functions/time/AsTime
  returnType: int
  signatures: [TIME.Year]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Year }} → 2023
```
