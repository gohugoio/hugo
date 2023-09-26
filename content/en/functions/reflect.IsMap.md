---
title: reflect.IsMap
description: Reports if a value is a map.
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
  - reflect.IsMap INPUT
---

`reflect.IsMap` reports if `VALUE` is a map.  Returns a boolean.

```go-html-template
{{ reflect.IsMap (dict "key" "value") }} → true
{{ reflect.IsMap "yo" }} → false
```
