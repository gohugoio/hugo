---
title: .Format
description: Formats built-in Hugo dates---`.Date`, `.PublishDate`, and `.Lastmod`---according to Go's layout string.
godocref: https://golang.org/pkg/time/#example_Time_Format
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time]
signature: [".Format FORMAT"]
workson: [times]
hugoversion:
relatedfuncs: [dateFormat,now,Unix,time]
deprecated: false
aliases: []
toc: true
---

`.Format` will format date values defined in your front matter and can be used as a property on the following [page variables][pagevars]:

* `.PublishDate`
* `.Date`
* `.Lastmod`

Assuming a key-value of `date: 2017-03-03` in a content file's front matter, your can run the date through `.Format` followed by a layout string for your desired output at build time:

```
{{ .PublishDate.Format "January 2, 2006" }} => March 3, 2017
```

For formatting *any* string representations of dates defined in your front matter, see the [`dateFormat` function][dateFormat], which will still leverage the Go layout string explained below but uses a slightly different syntax.

## Go's Layout String

Hugo templates [format your dates][time] via layout strings that point to a specific reference time:

```
Mon Jan 2 15:04:05 MST 2006
```

While this may seem arbitrary, the numerical value of `MST` is `07`, thus making the layout string a sequence of numbers.

Here is a visual explanation [taken directly from the Go docs][gdex]:

```
 Jan 2 15:04:05 2006 MST
=> 1 2  3  4  5    6  -7
```

### Hugo Date and Time Templating Reference

The following examples show the layout string followed by the rendered output.

The examples were rendered and tested in [CST][] and all point to the same field in a content file's front matter:

```
date: 2017-03-03T14:15:59-06:00
```

`.Date` (i.e. called via [page variable][pagevars])
: **Returns**: `2017-03-03 14:15:59 -0600 CST`

`"Monday, January 2, 2006"`
: **Returns**: `Friday, March 3, 2017`

`"Mon Jan 2 2006"`
: **Returns**: `Fri Mar 3 2017`

`"January 2006"`
: **Returns**: `March 2017`

`"2006-01-02"`
: **Returns**: `2017-03-03`

`"Monday"`
: **Returns**: `Friday`

`"02 Jan 06 15:04 MST"` (RFC822)
: **Returns**: `03 Mar 17 14:15 CST`

`"02 Jan 06 15:04 -0700"` (RFC822Z)
: **Returns**: `03 Mar 17 14:15 -0600`

`"Mon, 02 Jan 2006 15:04:05 MST"` (RFC1123)
: **Returns**: `Fri, 03 Mar 2017 14:15:59 CST`

`"Mon, 02 Jan 2006 15:04:05 -0700"` (RFC1123Z)
: **Returns**: `Fri, 03 Mar 2017 14:15:59 -0600`

More examples can be found in Go's [documentation for the time package][timeconst].

### Cardinal Numbers and Ordinal Abbreviations

Spelled-out cardinal numbers (e.g. "one", "two", and "three") and ordinal abbreviations (i.e., with shorted suffixes like "1st", "2nd", and "3rd") are not currently supported:

```
{{.Date.Format "Jan 2nd 2006"}}
```

Hugo assumes you want to append `nd` as a string to the day of the month and outputs the following:

```
Mar 3nd 2017
```

<!-- Content idea: see https://discourse.gohugo.io/t/formatting-a-date-with-suffix-2nd/5701 -->

### Use `.Local` and `.UTC`

In conjunction with the [`dateFormat` function][dateFormat], you can also convert your dates to `UTC` or to local timezones:

`{{ dateFormat "02 Jan 06 15:04 MST" .Date.UTC }}`
: **Returns**: `03 Mar 17 20:15 UTC`

`{{ dateFormat "02 Jan 06 15:04 MST" .Date.Local }}`
: **Returns**: `03 Mar 17 14:15 CST`

[CST]: https://en.wikipedia.org/wiki/Central_Time_Zone
[dateFormat]: /functions/dateformat/
[gdex]: https://golang.org/pkg/time/#example_Time_Format
[pagevars]: /variables/page/
[time]: https://golang.org/pkg/time/
[timeconst]: https://golang.org/pkg/time/#ANSIC
