---
title: Truncate
description: Returns the result of rounding DURATION1 toward zero to a multiple of DURATION2.
categories: []
keywords: []
action:
  related:
  related:
    - functions/time/Duration
    - functions/time/ParseDuration
  returnType: time.Duration
  signatures: [DURATION1.Truncate DURATION2]
---

```go-html-template
{{ $d = time.ParseDuration "3.5h2.5m1.5s" }}

{{ $d.Truncate (time.ParseDuration "2h") }} → 2h0m0s
{{ $d.Truncate (time.ParseDuration "3m") }} → 3h30m0s
{{ $d.Truncate (time.ParseDuration "4s") }} → 3h32m28s
```
