---
title: Truncate
description: Returns the result of rounding TIME down to a multiple of DURATION since January 1, 0001, 00:00:00 UTC.
categories: []
keywords: []
action:
  related:
    - functions/time/AsTime
    - functions/time/ParseDuration
    - methods/time/Round
  returnType: time.Time
  signatures: [TIME.Truncate DURATION]
---

The `Truncate` method operates on TIME as an absolute duration since the [zero time]; it does not operate on the presentation form of the time. If DURATION is a multiple of one hour, `Truncate` may return a time with a non-zero minute, depending on the time zone.

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $d := time.ParseDuration "1h"}}

{{ ($t.Truncate $d).Format "2006-01-02T15:04:05-00:00" }} → 2023-01-27T23:00:00-00:00
```

[zero time]: /getting-started/glossary/#zero-time
