---
title: reflect.IsSlice
description: Reports if a value is a slice.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: reflect
relatedFuncs:
  - reflect.IsMap
  - reflect.IsSlice
signature: 
  - reflect.IsSlice INPUT
---

`reflect.IsSlice` reports if `VALUE` is a slice.  Returns a boolean.

```go-html-template
{{ reflect.IsSlice (slice 1 2 3) }} → true
{{ reflect.IsSlice "yo" }} → false
```
