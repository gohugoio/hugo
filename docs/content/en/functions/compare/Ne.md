---
title: compare.Ne
linkTitle: ne
description: Returns the boolean truth of arg1 != arg2 && arg1 != arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [ne]
  returnType: bool
  signatures: ['compare.Ne ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/ne]
---

```go-html-template
{{ ne 1 1 }} → false
{{ ne 1 2 }} → true

{{ ne 1 1 1 }} → false
{{ ne 1 1 2 }} → false
{{ ne 1 2 1 }} → false
{{ ne 1 2 2 }} → true
```
