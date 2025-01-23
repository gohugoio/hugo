---
title: time.AsTime
description: Returns the given string representation of a date/time value as a time.Time value.
categories: []
keywords: []
action:
  aliases: [time]
  related:
    - functions/time/Duration
    - functions/time/Format
    - functions/time/Now
    - functions/time/ParseDuration
  returnType: time.Time
  signatures: ['time.AsTime INPUT [TIMEZONE]']
aliases: [/functions/time]
toc: true
---

## Overview

Hugo provides [functions] and [methods] to format, localize, parse, compare, and manipulate date/time values. Before you can do any of these with string representations of date/time values, you must first convert them to [`time.Time`] values using the `time.AsTime` function.

```go-html-template
{{ $t := "2023-10-15T13:18:50-07:00" }}
{{ time.AsTime $t }} → 2023-10-15 13:18:50 -0700 PDT (time.Time)
```

## Parsable strings

As shown above, the first argument must be a parsable string representation of a date/time value. For example:

{{% include "functions/time/_common/parsable-date-time-strings.md" %}}

To override the default time zone, set the [`timeZone`] in your site configuration or provide a second argument to the `time.AsTime` function. For example:

```go-html-template
{{ time.AsTime "15 Oct 2023" "America/Los_Angeles" }}
```

The list of valid time zones may be system dependent, but should include `UTC`, `Local`, or any location in the [IANA Time Zone database].

The order of precedence for determining the time zone is:

1. The time zone offset in the date/time string
1. The time zone provided as the second argument to the `time.AsTime` function
1. The time zone specified in your site configuration
1. The `Etc/UTC` time zone

[IANA Time Zone database]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
[`time.Time`]: https://pkg.go.dev/time#Time
[`timeZone`]: /getting-started/configuration/#timezone
[functions]: /functions/time/
[methods]: /methods/time/
