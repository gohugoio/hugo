---
title: reflect.IsMap
description: Reports if a value is a map.
godocref:
date: 2018-11-28
publishdate: 2018-11-28
lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.IsMap INPUT"]
workson: []
hugoversion: "v0.53"
relatedfuncs: [reflect.IsSlice]
deprecated: false
---

`reflect.IsMap` reports if `VALUE` is a map.  Returns a boolean.

```
{{ reflect.IsMap (dict "key" "value") }} → true
{{ reflect.IsMap "yo" }} → false
```
