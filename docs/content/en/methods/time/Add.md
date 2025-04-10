---
title: Add
description: Returns the given time plus the given duration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: time.Time
    signatures: [TIME.Add DURATION]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}

{{ $d1 = time.ParseDuration "3h20m10s" }}
{{ $d2 = time.ParseDuration "-3h20m10s" }}

{{ $t.Add $d1 }} → 2023-01-28 03:05:08 -0800 PST
{{ $t.Add $d2 }} → 2023-01-27 20:24:48 -0800 PST
```
