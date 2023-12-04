---
title: compare.Ge
description: Returns the boolean truth of arg1 >= arg2 && arg1 >= arg3.
categories: []
keywords: []
action:
  aliases: [ge]
  related:
    - functions/compare/Eq
    - functions/compare/Gt
    - functions/compare/Le
    - functions/compare/Lt
    - functions/compare/Ne
  returnType: bool
  signatures: ['compare.Ge ARG1 ARG2 [ARG...]']
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
