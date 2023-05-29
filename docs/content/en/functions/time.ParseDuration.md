---
title: time.ParseDuration
description: Parses a given duration string into a `time.Duration` structure.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [time parse duration]
signature: ["time.ParseDuration DURATION"]
---

`time.ParseDuration` parses a duration string into a [`time.Duration`](https://pkg.go.dev/time#Duration) structure so you can access its fields.
A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as `300ms`, `-1.5h` or `2h45m`. Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

You can perform [time operations](https://pkg.go.dev/time#Duration) on the returned `time.Duration` value:

    {{ printf "There are %.0f seconds in one day." (time.ParseDuration "24h").Seconds }}
    <!-- Output: There are 86400 seconds in one day. -->
