---
title: Microseconds
description: Returns the time.Duration value as an integer microsecond count.
categories: []
keywords: []
params:
  functions_and_methods:
    related:
      - functions/time/Duration
      - functions/time/ParseDuration
    returnType: int64
    signatures: [DURATION.Microseconds]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}
{{ $d.Microseconds }} → 12751500000
```
