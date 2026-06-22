---
title: UnixMicro
description: Returns the given time.Time value expressed as the number of microseconds elapsed since January 1, 1970 UTC.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int64
    signatures: [TIME.UnixMicro]
---

See [Unix epoch][].

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.UnixMicro }} → 1674891898000000
```

[Unix epoch]: https://en.wikipedia.org/wiki/Unix_time
