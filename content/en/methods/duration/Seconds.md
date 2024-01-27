---
title: Seconds
description: Returns the time.Duration value as a floating point number of seconds.
categories: []
keywords: []
action:
  related:
    - functions/time/Duration
    - functions/time/ParseDuration
  returnType: float64
  signatures: [DURATION.Seconds]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}
{{ $d.Seconds }} → 12751.5
```
