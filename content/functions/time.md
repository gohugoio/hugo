---
title: time
linktitle:
description: Converts a timestamp string into a `time.Time` structure.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time]
signature: ["time INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`time` converts a timestamp string into a [`time.Time`](https://godoc.org/time#Time) structure so you can access its fields:

```
{{ time "2016-05-28" }} → "2016-05-28T00:00:00Z"
{{ (time "2016-05-28").YearDay }} → 149
{{ mul 1000 (time "2016-05-28T10:30:00.00+10:00").Unix }} → 1464395400000, or Unix time in milliseconds
```

## Example: Using `time` to get Month Index

The following example takes a UNIX timestamp---set as `utimestamp: "1489276800"` in a content's front matter---converts the timestamp (string) to an integer using the [`int` function][int], and then uses [`printf`][] to convert the `Month` property of `time` into an index. 

The following example may be useful when setting up [multilingual sites][multilingual]:

{{< code file="unix-to-month-integer.html" >}}
{{$time := time (int .Params.addDate)}}
=> $time = 1489276800
{{$time.Month}}
=> "March"
{{$monthindex := printf "%d" $time.Month }}
=> $monthindex = 3
{{< /code >}}


[int]: /functions/int/
[multilingual]: /content-management/multilingual/
[`printf`]: /functions/printf/
