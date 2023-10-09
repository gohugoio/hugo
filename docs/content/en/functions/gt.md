---
title: gt
description: Returns the boolean truth of arg1 > arg2 && arg1 > arg3.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [comparison,operators,logic]
signature: ["gt ARG1 ARG2 [ARG...]"]
relatedfuncs: []
---

```go-html-template
{{ gt 1 1 }} → false
{{ gt 1 2 }} → false
{{ gt 2 1 }} → true

{{ gt 1 1 1 }} → false
{{ gt 1 1 2 }} → false
{{ gt 1 2 1 }} → false
{{ gt 1 2 2 }} → false

{{ gt 2 1 1 }} → true
{{ gt 2 1 2 }} → false
{{ gt 2 2 1 }} → false
```
