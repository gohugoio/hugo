---
title: index
linktitle: index
description: Looks up the index(es) or key(s) of the data structure passed into it.
godocref: https://golang.org/pkg/text/template/#hdr-Functions
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
needsexample: true
---

`index` looks up the index(es) or key(s) of the data structure passed into it.

From the godocs:

> Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.

In Go templates, you can't access array, slice, or map elements directly the same way you would in Go. For example, `$.Site.Data.authors[.Params.authorkey]` isn't supported syntax.

Instead, you have to use `index`, a function that handles the lookup for you.
