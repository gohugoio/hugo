---
title: Before
description: Reports whether TIME1 is before TIME2.
categories: []
keywords: []
action:
  related:
    - methods/time/After
    - methods/time/Equal
    - functions/time/AsTime
  returnType: bool
  signatures: [TIME1.Before TIME2]
---

```go-html-template
{{ $t1 := time.AsTime "2023-01-01T17:00:00-08:00" }}
{{ $t2 := time.AsTime "2030-01-01T17:00:00-08:00" }}

{{ $t1.Before $t2 }} → true
