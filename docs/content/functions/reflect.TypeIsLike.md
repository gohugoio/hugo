---
title: reflect.TypeIsLike
description: Reports whether `VALUE` is of a type or a pointer to a type.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.TypeIsLike TYPE VALUE"]
workson: []
hugoversion:
relatedfuncs: [reflect.TypeIs, reflect.TypeOf]
deprecated: false
aliases: []
---

`reflect.TypeIsLike` reports whether `VALUE` is of type `TYPE` or a pointer to type `TYPE`.

    {{ reflect.TypeIsLike "hugolib.Page" . }} → true

    {{ reflect.TypeIsLike "int" 1 }} → true
    {{ reflect.TypeIsLike "int" "1" }} → false
    {{ reflect.TypeIsLike "[]int" (slice 1 2 3) }} → true
    {{ reflect.TypeIsLike "map[string]interface {}" (dict "key" "value") }} → true

