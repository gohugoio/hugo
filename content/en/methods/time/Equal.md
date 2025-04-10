---
title: Equal
description: Reports whether TIME1 is equal to TIME2.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [TIME1.Equal TIME2]
---

```go-html-template
{{ $t1 := time.AsTime "2023-01-01T17:00:00-08:00" }}
{{ $t2 := time.AsTime "2023-01-01T20:00:00-05:00" }}

{{ $t1.Equal $t2 }} → true
```
