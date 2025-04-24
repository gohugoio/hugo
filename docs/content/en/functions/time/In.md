---
title: time.In
description: Returns the given date/time as represented in the specified IANA time zone.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: time.Time
    signatures: [time.In TIMEZONE INPUT]
---

{{< new-in 0.146.0 />}}

The `time.In` function returns the given date/time as represented in the specified [IANA](g) time zone.

- If the time zone is an empty string or `UTC`, the time is returned in [UTC](g).
- If the time zone is `Local`, the time is returned in the system's local time zone.
- Otherwise, the time zone must be a valid IANA [time zone name].

[time zone name]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List

```go-html-template
{{ $layout := "2006-01-02T15:04:05-07:00" }}
{{ $t := time.AsTime "2025-03-31T14:45:00-00:00" }}

{{ $t | time.In "America/Denver" | time.Format $layout }}     → 2025-03-31T08:45:00-06:00
{{ $t | time.In "Australia/Adelaide" | time.Format $layout }} → 2025-04-01T01:15:00+10:30
{{ $t | time.In "Europe/Oslo" | time.Format $layout }}        → 2025-03-31T16:45:00+02:00
```
