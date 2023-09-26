---
title: ne
description: Returns the boolean truth of arg1 != arg2 && arg1 != arg3.
categories: [functions]
menu:
  docs:
    parent: functions
namespace: compare
relatedFuncs:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
signature:
  - compare.Ne ARG1 ARG2 [ARG...]
  - ne ARG1 ARG2 [ARG...]
---

```go-html-template
{{ ne 1 1 }} → false
{{ ne 1 2 }} → true

{{ ne 1 1 1 }} → false
{{ ne 1 1 2 }} → false
{{ ne 1 2 1 }} → false
{{ ne 1 2 2 }} → true
```
