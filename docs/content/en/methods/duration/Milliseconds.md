---
title: Milliseconds
description: Returns the time.Duration value as an integer millisecond count.
categories: []
keywords: []
action:
  related:
    - functions/time/Duration
    - functions/time/ParseDuration
  returnType: int64
  signatures: [DURATION.Milliseconds]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}
{{ $d.Milliseconds }} → 12751500
```
