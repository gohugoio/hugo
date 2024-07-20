---
title: time.ParseDuration
description: Returns a time.Duration value by parsing the given duration string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/time/AsTime
    - functions/time/Duration
    - functions/time/Format
    - functions/time/Now
  returnType: time.Duration
  signatures: [time.ParseDuration DURATION]
aliases: [/functions/time.parseduration]
---

The `time.ParseDuration` function returns a time.Duration value that you can use with any of the `Duration` [methods].


A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as `300ms`, `-1.5h` or `2h45m`. Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

This template:

```go-html-template
{{ $duration := time.ParseDuration "24h" }}
{{ printf "There are %.0f seconds in one day." $duration.Seconds }}
```

Is rendered to:

```text
There are 86400 seconds in one day.
```

[`time.Duration`]: https://pkg.go.dev/time#Duration
[methods]: /methods/duration/
