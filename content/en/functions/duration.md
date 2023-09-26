---
title: duration
description: Returns a `time.Duration` structure, using the given time unit and duration number.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: time
relatedFuncs:
  - time.AsTime
  - time.Duration
  - time.Format
  - time.Now
  - time.ParseDuration
signature:
  - time.Duration TIME_UNIT DURATION_NUMBER
  - duration TIME_UNIT DURATION_NUMBER
---

`time.Duration` converts a given number into a [`time.Duration`](https://pkg.go.dev/time#Duration) structure so you can access its fields. E.g. you can perform [time operations](https://pkg.go.dev/time#Duration) on the returned `time.Duration` value:

```go-html-template
{{ printf "There are %.0f seconds in one day." (duration "hour" 24).Seconds }}
<!-- Output: There are 86400 seconds in one day. -->
```

Make your code simpler to understand by using a [chained pipeline](https://pkg.go.dev/text/template#hdr-Pipelines):

```go-html-template
{{ mul 7.75 60 | duration "minute" }} → 7h45m0s
{{ mul 120 60 | mul 1000 | duration "millisecond" }} → 2h0m0s
```

You have to specify a time unit for the number given to the function. Valid time units are:

Duration|Valid time units
:--|:--
hours|`hour`, `h`
minutes|`minute`, `m`
seconds|`second`, `s`
milliseconds|`millisecond`, `ms`
microseconds|`microsecond`, `us`, `µs`
nanoseconds|`nanosecond`, `ns`
