---
title: reflect.TypeIs
description: Reports whether a value is of a type.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.TypeIs TYPE VALUE"]
workson: []
hugoversion:
relatedfuncs: [reflect.TypeIsLike, reflect.TypeOf]
deprecated: false
aliases: []
---

`reflect.TypeIs` reports whether `VALUE` is of type `TYPE`.

    {{ reflect.TypeIs "int" 1 }} → true
    {{ reflect.TypeIs "int" "1" }} → false
    {{ reflect.TypeIs "[]int" (slice 1 2 3) }} → true
    {{ reflect.TypeIs "map[string]interface {}" (dict "key" "value") }} → true

