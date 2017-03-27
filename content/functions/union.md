---
title: union
linktitle: union
description: Given two arrays or slices, returns a new array that contains the elements or objects that belong to either or both arrays/slices.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-12
categories: [functions]
tags: [filtering,lists]
ns:
signature: []
workson: []
hugoversion: 0.20
relatedfuncs: [intersect,where]
deprecated: false
aliases: []
---

Given two arrays (or slices) A and B, this function will return a new array that contains the elements or objects that belong to either A or to B or to both. The elements supported are strings, integers, and floats (only float64).

```golang
{{ union (slice 1 2 3) (slice 3 4 5) }}
<!-- returns [1 2 3 4 5] -->

{{ union (slice 1 2 3) nil }}
<!-- returns [1 2 3] -->

{{ union nil (slice 1 2 3) }}
<!-- returns [1 2 3] -->

{{ union nil nil }}
<!-- returns an error because both arrays/slices have to be of the same type -->
```