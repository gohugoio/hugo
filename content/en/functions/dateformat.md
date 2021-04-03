---
title: dateFormat
description: Converts the textual representation of the `datetime` into the specified format.
godocref: https://golang.org/pkg/time/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time,strings]
signature: ["dateFormat LAYOUT INPUT"]
workson: []
hugoversion:
relatedfuncs: [Format,now,Unix,time]
deprecated: false
---

`dateFormat` converts an [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) timestamp string `INPUT` into the format specified by the `LAYOUT` string.

```
{{ dateFormat "Monday, Jan 2, 2006" "2015-01-21" }} → "Wednesday, Jan 21, 2015"
```

{{% warning %}}
As of v0.19 of Hugo, the `dateFormat` function is *not* supported as part of Hugo's [multilingual feature](/content-management/multilingual/).
{{% /warning %}}

See [Go’s Layout String](/functions/format/#gos-layout-string) to learn about how the `LAYOUT` string has to be formatted. There are also some useful examples.

See the [`time` function](/functions/time/) to convert an [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) timestamp string to a Go `time.Time` type value.
