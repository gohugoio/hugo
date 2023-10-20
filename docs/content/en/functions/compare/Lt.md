---
title: compare.Lt
linkTitle: lt
description: Returns the boolean truth of arg1 < arg2 && arg1 < arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [lt]
  returnType: bool
  signatures: ['compare.Lt ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/lt]
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
