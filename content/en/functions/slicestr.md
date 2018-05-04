---
title: slicestr
# linktitle:
description: Creates a slice of a half-open range, including start and end indices.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["slicestr STRING START [END]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

For example, 1 and 4 creates a slice including elements 1 through 3.
The `end` index can be omitted; it defaults to the string's length.

* `{{slicestr "BatMan" 3}}` → "Man"
* `{{slicestr "BatMan" 0 3}}` → "Bat"

