---
title: Nanoseconds
description: Returns the time.Duration value as an integer nanosecond count.
categories: []
keywords: []
action:
  related:
    - functions/time/Duration
    - functions/time/ParseDuration
  returnType: int64
  signatures: [DURATION.Nanoseconds]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}
{{ $d.Nanoseconds }} → 12751500000000
```
