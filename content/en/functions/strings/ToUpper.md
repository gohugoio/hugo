---
title: strings.ToUpper
linkTitle: upper
description: Converts all characters in a string to uppercase
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [upper]
  returnType: string
  signatures: [strings.ToUpper INPUT]
relatedFunctions:
  - strings.FirstUpper
  - strings.Title
  - strings.ToLower
  - strings.ToUpper
aliases: [/functions/upper]
---

Note that `upper` can be applied in your templates in more than one way:

```go-html-template
{{ upper "BatMan" }} → "BATMAN"
{{ "BatMan" | upper }} → "BATMAN"
```
