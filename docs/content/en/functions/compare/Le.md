---
title: compare.Le
linkTitle: le
description: Returns the boolean truth of arg1 <= arg2 && arg1 <= arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [le]
  returnType: bool
  signatures: ['compare.Le ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/le]
---

```go-html-template
{{ le 1 1 }} → true
{{ le 1 2 }} → true
{{ le 2 1 }} → false

{{ le 1 1 1 }} → true
{{ le 1 1 2 }} → true
{{ le 1 2 1 }} → true
{{ le 1 2 2 }} → true

{{ le 2 1 1 }} → false
{{ le 2 1 2 }} → false
{{ le 2 2 1 }} → false
```
