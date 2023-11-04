---
title: time.Format
description: Returns a formatted and localized time.Time value.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [dateFormat]
  returnType: string
  signatures: [time.Format LAYOUT INPUT]
relatedFunctions:
  - time.AsTime
  - time.Duration
  - time.Format
  - time.Now
  - time.ParseDuration
aliases: [/functions/dateformat]
toc: true
---

```go-template
{{ $t := "2023-01-27T23:44:58-08:00" }}
{{ $format := "2 Jan 2006" }}

{{ $t | time.Format $format }} → 27 Jan 2023

{{ $t = time.AsTime $t }}
{{ $t | time.Format $format }} → 27 Jan 2023
```

## Layout string

{{% readfile file="/functions/_common/time-layout-string.md" %}}

## Localization

Use the `time.Format` function to localize `time.Time` values for the current language and region.

{{% note %}}
{{% readfile file="/functions/_common/locales.md" %}}
{{% /note %}}


Use the layout string as described above, or one of the tokens below. For example:

```go-template
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
