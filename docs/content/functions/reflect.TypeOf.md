---
title: reflect.TypeOf
description: Reports a value's type.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.TypeOf VALUE"]
workson: []
hugoversion:
relatedfuncs: [reflect.TypeIs, reflect.TypeIsLike]
deprecated: false
aliases: []
---

`reflect.TypeOf` reports `VALUE`'s type.

    {{ reflect.TypeOf 1 }} → "int"
    {{ reflect.TypeOf "1" }} → "int"
    {{ reflect.TypeOf (slice 1 2 3) }} → "[]int"
    {{ reflect.TypeOf (dict "key" "value") }} → "map[string]interface {}"

