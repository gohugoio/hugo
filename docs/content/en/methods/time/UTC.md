---
title: UTC
description: Returns the given time.Time value with the location set to UTC.
categories: []
keywords: []
action:
  related:
    - methods/time/Local
    - functions/time/AsTime
  returnType: time.Time
  signatures: [TIME.UTC]
---

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.UTC }} → 2023-01-28 07:44:58 +0000 UTC
