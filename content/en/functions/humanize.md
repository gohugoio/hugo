---
title: humanize
linktitle:
description: Returns the humanized version of an argument with the first letter capitalized.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,casing]
signature: ["humanize INPUT"]
workson: []
hugoversion:
relatedfuncs: [anchorize]
deprecated: false
aliases: []
---

If the input is either an int64 value or the string representation of an integer, humanize returns the number with the proper ordinal appended.


```
{{humanize "my-first-post"}} → "My first post"
{{humanize "myCamelPost"}} → "My camel post"
{{humanize "52"}} → "52nd"
{{humanize 103}} → "103rd"
```
