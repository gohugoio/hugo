---
title: Unix
description: Returns the given time.Time value expressed as the number of seconds elapsed since January 1, 1970 UTC.
categories: []
params:
  functions_and_methods:
    returnType: int64
    signatures: [TIME.Unix]
aliases: [/functions/unix]
---

See [Unix epoch](https://en.wikipedia.org/wiki/Unix_time).

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Unix }} → 1674891898
```
