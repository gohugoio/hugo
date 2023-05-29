---
title: reflect.IsMap
description: Reports if a value is a map.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [reflect, reflection, kind]
signature: ["reflect.IsMap INPUT"]
relatedfuncs: [reflect.IsSlice]
---

`reflect.IsMap` reports if `VALUE` is a map.  Returns a boolean.

```go-html-template
{{ reflect.IsMap (dict "key" "value") }} → true
{{ reflect.IsMap "yo" }} → false
```
