---
title: lt
description: Returns the boolean truth of arg1 < arg2 && arg1 < arg3.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: compare
relatedFuncs:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
signature:
  - compare.Lt ARG1 ARG2 [ARG...]
  - lt ARG1 ARG2 [ARG...]
---

```go-html-template
{{ lt 1 1 }} → false
{{ lt 1 2 }} → true
{{ lt 2 1 }} → false

{{ lt 1 1 1 }} → false
{{ lt 1 1 2 }} → false
{{ lt 1 2 1 }} → false
{{ lt 1 2 2 }} → true

{{ lt 2 1 1 }} → false
{{ lt 2 1 2 }} → false
{{ lt 2 2 1 }} → false
```
