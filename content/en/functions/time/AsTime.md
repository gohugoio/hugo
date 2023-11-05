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
{{ $t := "2023-10-15T14:20:28-07:00" }}
{{ time.AsTime $t }} → 2023-10-15 14:20:28 -0700 PDT (time.Time)
```

## Parsable strings

As shown above, the first argument must be a *parsable* string representation of a date/time value. For example:

{{% include "functions/time/_common/parsable-date-time-strings.md" %}}

## Time zones

When the parsable string does not contain a time zone offset, you can do either of the following to assign a time zone other than Etc/UTC:

- Provide a second argument to the `time.AsTime` function

  ```go-html-template
  {{ time.AsTime "15 Oct 2023" "America/Chicago" }}
  ```

- Set the default time zone in your site configuration

  {{< code-toggle file=hugo >}}
  timeZone = 'America/New_York'
  {{< /code-toggle >}}

The order of precedence for determining the time zone is:

1. The time zone offset in the date/time string
2. The time zone provide as the second argument to the `time.AsTime` function
3. The time zone specified in your site configuration

The list of valid time zones may be system dependent, but should include `UTC`, `Local`, or any location in the [IANA Time Zone database].

[`time.Time`]: https://pkg.go.dev/time#Time
[functions]: /functions/time/
[iana time zone database]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
[methods]: /methods/time/
