---
title: compare.Eq
linkTitle: eq
description: Returns the boolean truth of arg1 == arg2 || arg1 == arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [eq]
  returnType: bool
  signatures: ['compare.Eq ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/eq]
---

```go-html-template
{{ eq 1 1 }} → true
{{ eq 1 2 }} → false

{{ eq 1 1 1 }} → true
{{ eq 1 1 2 }} → true
{{ eq 1 2 1 }} → true
{{ eq 1 2 2 }} → false
```
