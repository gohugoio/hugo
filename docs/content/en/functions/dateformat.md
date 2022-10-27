---
title: time.Format
description: Converts a date/time to a localized string.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time,strings]
signature: ["time.Format LAYOUT INPUT [LANGUAGE]"]
workson: []
hugoversion:
relatedfuncs: [Format,now,Unix,time]
deprecated: false
---

`time.Format` (alias `dateFormat`) converts either a `time.Time` object (e.g. `.Date`) or a timestamp string `INPUT` into the format specified by the `LAYOUT` string and an optional `LANGUAGE` code.

```go-html-template
{{ time.Format "Monday, Jan 2, 2006" "2015-12-21" }} → "Monday, Dec 21, 2015"
```

```go-html-template
{{ time.Format "Monday, Jan 2, 2006" "2015-12-21" "pt" }} → "lunes, dic. 21, 2015"
```

Note that since Hugo 0.87.0, `time.Format` will return a localized string for the current language. {{< new-in "0.87.0" >}}

The `LAYOUT` string can be either:

* [Go’s Layout String](/functions/format/#gos-layout-string) to learn about how the `LAYOUT` string has to be formatted. There are also some useful examples.
* A custom Hugo layout identifier (see full list below)

See the [`time` function](/functions/time/) to convert a timestamp string to a Go `time.Time` type value.


## Date/time formatting layouts

{{< new-in "0.87.0" >}}

Go's date layout strings can be hard to reason about, especially with multiple languages. Since Hugo 0.87.0 you can alternatively use some predefined layout identifiers that will output localized dates or times:

```go-html-template
{{ .Date | time.Format ":date_long" }}
```

The full list of custom layouts with examples for English:

* `:date_full` => `Wednesday, June 6, 2018`
* `:date_long` => `June 6, 2018`
* `:date_medium` => `Jun 6, 2018`
* `:date_short` => `6/6/18`

* `:time_full` => `2:09:37 am UTC`
* `:time_long` => `2:09:37 am UTC`
* `:time_medium` => `2:09:37 am`
* `:time_short` => `2:09 am`
