---
title: strings.ToLower
linkTitle: lower
description: Converts all characters in the provided string to lowercase.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [lower]
  returnType: string
  signatures: [strings.ToLower INPUT]
relatedFunctions:
  - strings.FirstUpper
  - strings.Title
  - strings.ToLower
  - strings.ToUpper
aliases: [/functions/lower]
---


Note that `lower` can be applied in your templates in more than one way:

```go-html-template
{{ lower "BatMan" }} → "batman"
{{ "BatMan" | lower }} → "batman"
```
