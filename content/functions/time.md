---
title: time
linktitle:
description: Converts a timestamp string into a `time.Time` structure.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [dates,time]
signature:
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