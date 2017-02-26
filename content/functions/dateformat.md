---
title: dateFormat
linktitle:
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [dates,time,strings]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

`dateFormat` converts the textual representation of the `datetime` into the specified format or returns it as a Go `time.Time` type value. These are formatted with the layout string.

```
{{ dateFormat "Monday, Jan 2, 2006" "2015-01-21" }} → "Wednesday, Jan 21, 2015"
```