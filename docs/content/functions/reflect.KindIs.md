---
title: reflect.KindIs
description: Reports whether a value is of a kind.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.KindIs KIND VALUE"]
workson: []
hugoversion:
relatedfuncs: [reflect.KindOf]
deprecated: false
aliases: []
---

`reflect.KindIs` reports whether `VALUE` is of kind `KIND`.

    {{ reflect.KindIs "int" 1 }} → true
    {{ reflect.KindIs "int" "1" }} → false
    {{ reflect.KindIs "slice" (slice 1 2 3) }} → true
    {{ reflect.KindIs "map" (dict "key" "value") }} → true

`KIND` is one of:

* array
* bool
* chan
* complex128
* complex64
* float32
* float64
* func
* int
* int16
* int32
* int64
* int8
* interface
* invalid
* map
* ptr
* slice
* string
* struct
* uint
* uint16
* uint32
* uint64
* uint8
* uintptr
* unsafePointer
