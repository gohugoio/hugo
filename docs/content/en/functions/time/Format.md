---
title: time.Format
description: Returns the given date/time as a formatted and localized string.
categories: []
keywords: []
action:
  aliases: [dateFormat]
  related:
    - functions/time/AsTime
    - functions/time/Duration
    - functions/time/Now
    - functions/time/ParseDuration
  returnType: string
  signatures: [time.Format LAYOUT INPUT]
aliases: [/functions/dateformat]
toc: true
---

Use the `time.Format` function with `time.Time` values:

```go-html-template
{{ $t := time.AsTime "2023-02-27T23:44:58-08:00" }}
{{ time.Format "2 Jan 2006" $t }} → 27 Feb 2023
```

Or use `time.Format` with a *parsable* string representation of a date/time value:

```go-html-template
{{ $t := "27 Feb 2023" }}
{{ time.Format "January 2, 2006" $t }} → February 27, 2023
```

Examples of parsable string representations:

{{% include "functions/time/_common/parsable-date-time-strings.md" %}}

## Layout string

{{% include "functions/_common/time-layout-string.md" %}}

## Localization

Use the `time.Format` function to localize `time.Time` values for the current language and region.

{{% include "functions/_common/locales.md" %}}

Use the layout string as described above, or one of the tokens below. For example:

```go-html-template
{{ .Date | time.Format ":date_medium" }} → Jan 27, 2023
```

Localized to en-US:

Token|Result
:--|:--
`:date_full`|`Friday, January 27, 2023`
`:date_long`|`January 27, 2023`
`:date_medium`|`Jan 27, 2023`
`:date_short`|`1/27/23`
`:time_full`|`11:44:58 pm Pacific Standard Time`
`:time_long`|`11:44:58 pm PST`
`:time_medium`|`11:44:58 pm`
`:time_short`|`11:44 pm`

Localized to de-DE:

Token|Result
:--|:--
`:date_full`|`Freitag, 27. Januar 2023`
`:date_long`|`27. Januar 2023`
`:date_medium`|`27.01.2023`
`:date_short`|`27.01.23`
`:time_full`|`23:44:58 Nordamerikanische Westküsten-Normalzeit`
`:time_long`|`23:44:58 PST`
`:time_medium`|`23:44:58`
`:time_short`|`23:44`
