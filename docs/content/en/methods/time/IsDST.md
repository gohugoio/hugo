---
title: IsDST
description: Reports whether the given time.Time value is in Daylight Savings Time.
categories: []
keywords: []
action:
  related:
    - functions/time/AsTime
  returnType: bool
  signatures: [TIME.IsDST]
---

```go-html-template
{{ $t1 := time.AsTime "2023-01-01T00:00:00-08:00" }}
{{ $t2 := time.AsTime "2023-07-01T00:00:00-07:00" }}

{{ $t1.IsDST }} → false
{{ $t2.IsDST }} → true
```
