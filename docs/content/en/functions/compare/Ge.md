---
title: compare.Ge
linkTitle: ge
description: Returns the boolean truth of arg1 >= arg2 && arg1 >= arg3.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [ge]
  returnType: bool
  signatures: ['compare.Ge ARG1 ARG2 [ARG...]']
relatedFunctions:
  - compare.Eq
  - compare.Ge
  - compare.Gt
  - compare.Le
  - compare.Lt
  - compare.Ne
aliases: [/functions/ge]
---

```go-html-template
{{ ge 1 1 }} → true
{{ ge 1 2 }} → false
{{ ge 2 1 }} → true

{{ ge 1 1 1 }} → true
{{ ge 1 1 2 }} → false
{{ ge 1 2 1 }} → false
{{ ge 1 2 2 }} → false

{{ ge 2 1 1 }} → true
{{ ge 2 1 2 }} → true
{{ ge 2 2 1 }} → true
```
