Go templates [format your dates][time] according to a single reference time:

```
Mon Jan 2 15:04:05 MST 2006
```

You can think of `MST` as `07`, thus making the reference format string a sequence of numbers. The following is [taken directly from the Go docs][gdex]:

```
Jan 2 15:04:05 2006 MST
  1 2  3  4  5    6  -7
```

### Hugo Date Templating Reference

Each of the following examples show the reference formatting string followed by the string Hugo will output in your HTML.

Note that the examples were rendered and tested in [CST][] and pull from a single example date you might have in your content's front matter:

```
date: 2017-03-03T14:15:59-06:00
```

`.Date` (i.e. called via [page variable][pagevars])
: **Returns**: `2017-03-03 14:15:59 -0600 CST`

`"Monday, January 2, 2006"`
: **Returns**: `Friday, March 3, 2017`

`"Mon Jan 2 2006"`
: **Returns**: `Fri Mar 3 2017`

`"January 2nd"`
: **Returns**: `March 3rd`

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

`"Mon, 02 Jan 2006 15:04:05 -0700"` (RFC339)
: **Returns**: `Fri, 03 Mar 2017 14:15:59 -0600`

### Cardinal Numbers and Ordinal Abbreviations

Spelled-out cardinal numbers (e.g. "one", "two", and "three") and ordinal abbreviations (e.g. "1st", "2nd", and "3rd") are not currently supported.

To continue with the example above:

```
{{.Date.Format "Jan 2nd 2006"}}
```

Hugo assumes you want to append `nd` as a string to the day of the month and outputs the following:

```
Mar 3nd 2017
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