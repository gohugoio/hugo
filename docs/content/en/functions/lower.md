---
title: lower
linktitle: lower
description: Converts all characters in the provided string to lowercase.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,casing]
signature:
  - "lower INPUT"
  - "strings.ToLower INPUT"
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---


Note that `lower` can be applied in your templates in more than one way:

```go-html-template
{{ lower "BatMan" }} → "batman"
{{ "BatMan" | lower }} → "batman"
```
