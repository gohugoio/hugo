---
title: reflect.KindOf
description: Reports a value's kind.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [reflect, reflection, kind]
signature: ["reflect.KindOf VALUE"]
workson: []
hugoversion:
relatedfuncs: [reflect.KindIs]
deprecated: false
aliases: []
---

`reflect.KindOf` reports `VALUE`'s kind.

    {{ reflect.KindOf 1 }} → "int"
    {{ reflect.KindOf "1" }} → "string"
    {{ reflect.KindOf (slice 1 2 3) }} → "slice"
    {{ reflect.KindOf (dict "key" "value") }} → "map"

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
