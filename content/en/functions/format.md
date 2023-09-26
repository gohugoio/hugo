---
title: .Format
description: Formats built-in Hugo dates---`.Date`, `.PublishDate`, and `.Lastmod`---according to Go's layout string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace:
relatedFuncs: []
signature:
  - .Format FORMAT
toc: true
---

`.Format` will format date values defined in your front matter and can be used as a property on the following [page variables][pagevars]:

* `.PublishDate`
* `.Date`
* `.Lastmod`

Assuming a key-value of `date: 2017-03-03` in a content file's front matter, your can run the date through `.Format` followed by a layout string for your desired output at build time:

```go-html-template
{{ .PublishDate.Format "January 2, 2006" }} => March 3, 2017
```

For formatting *any* string representations of dates defined in your front matter, see the [`dateFormat` function][dateFormat], which will still leverage the Go layout string explained below but uses a slightly different syntax.

## Go's layout string

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

### Hugo date and time templating reference

The following examples show the layout string followed by the rendered output.

The examples were rendered and tested in [CST] and all point to the same field in a content file's front matter:

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

### Cardinal s

Spelled-out cardinal numbers (e.g. "one", "two", and "three") are not currently supported.

Use the [`humanize`](/functions/humanize) function to render the day of the month as an ordinal number:

```go-html-template
{{ humanize .Date.Day }} of {{ .Date.Format "January 2006" }}
```

This will output:

```
5th of March 2017
```


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
