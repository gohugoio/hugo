---
title: eq
description: Returns the boolean truth of arg1 == arg2 || arg1 == arg3.
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
  - compare.Eq ARG1 ARG2 [ARG...]
  - eq ARG1 ARG2 [ARG...]
---

```go-html-template
{{ eq 1 1 }} → true
{{ eq 1 2 }} → false

{{ eq 1 1 1 }} → true
{{ eq 1 1 2 }} → true
{{ eq 1 2 1 }} → true
{{ eq 1 2 2 }} → false
```
