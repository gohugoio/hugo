---
title: Sub
description: Returns the duration computed by subtracting TIME2 from TIME1.
categories: []
keywords: []
action:
  related:
    - functions/time/AsTime
  returnType: time.Duration
  signatures: [TIME1.Sub TIME2]
---

```go-html-template
{{ $t1 := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t2 := time.AsTime "2023-01-26T22:34:38-08:00" }}

{{ $t1.Sub $t2 }} → 25h10m20s
```
