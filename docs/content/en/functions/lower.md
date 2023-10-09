---
title: lower
description: Converts all characters in the provided string to lowercase.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings,casing]
signature:
  - "lower INPUT"
  - "strings.ToLower INPUT"
relatedfuncs: []
---


Note that `lower` can be applied in your templates in more than one way:

```go-html-template
{{ lower "BatMan" }} → "batman"
{{ "BatMan" | lower }} → "batman"
```
