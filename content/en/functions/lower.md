---
title: lower
description: Converts all characters in the provided string to lowercase.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.FirstUpper
  - strings.Title
  - strings.ToLower
  - strings.ToUpper
signature:
  - strings.ToLower INPUT
  - lower INPUT
---


Note that `lower` can be applied in your templates in more than one way:

```go-html-template
{{ lower "BatMan" }} → "batman"
{{ "BatMan" | lower }} → "batman"
```
