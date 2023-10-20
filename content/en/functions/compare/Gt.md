---
title: compare.Gt
linkTitle: gt
description: Returns the boolean truth of arg1 > arg2 && arg1 > arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [gt]
  returnType: bool
  signatures: ['compare.Gt ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/gt]
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
