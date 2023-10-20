---
title: reflect.IsSlice
description: Reports whether the value is a slice.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: bool
  signatures: [reflect.IsSlice INPUT]
namespace: reflect
relatedFunctions:
  - reflect.IsMap
  - reflect.IsSlice
aliases: [/functions/reflect.isslice]
---

```go-html-template
{{ reflect.IsSlice (slice 1 2 3) }} → true
{{ reflect.IsSlice "yo" }} → false
```
