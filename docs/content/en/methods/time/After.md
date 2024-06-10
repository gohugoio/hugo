---
title: After
description: Reports whether TIME1 is after TIME2.
categories: []
keywords: []
action:
  related:
    - methods/time/Before
    - methods/time/Equal
    - functions/time/AsTime
  returnType: bool
  signatures: [TIME1.After TIME2]
---

```go-html-template
{{ $t1 := time.AsTime "2023-01-01T17:00:00-08:00" }}
{{ $t2 := time.AsTime "2010-01-01T17:00:00-08:00" }}

{{ $t1.After $t2 }} → true
```
