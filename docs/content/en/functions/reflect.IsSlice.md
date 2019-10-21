---
title: reflect.IsSlice
description: Reports if a value is a slice.
godocref:
date: 2018-11-28
publishdate: 2018-11-28
lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.IsSlice INPUT"]
workson: []
hugoversion: "0.53"
relatedfuncs: [reflect.IsMap]
deprecated: false
---

`reflect.IsSlice` reports if `VALUE` is a slice.  Returns a boolean.

```
{{ reflect.IsSlice (slice 1 2 3) }} → true
{{ reflect.IsSlice "yo" }} → false
```
