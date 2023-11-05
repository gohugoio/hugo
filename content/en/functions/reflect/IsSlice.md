---
title: reflect.IsSlice
description: Reports whether the given value is a slice.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/reflect/IsMap
  returnType: bool
  signatures: [reflect.IsSlice INPUT]
aliases: [/functions/reflect.isslice]
---

```go-html-template
{{ reflect.IsSlice (slice 1 2 3) }} → true
{{ reflect.IsSlice "yo" }} → false
```
