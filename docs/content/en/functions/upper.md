---
title: upper
description: Converts all characters in a string to uppercase
keywords: []
categories: [functions]
menu:
  docs:
    parent: functions
toc:
signature:
  - "upper INPUT"
  - "strings.ToUpper INPUT"
relatedfuncs: []
---

Note that `upper` can be applied in your templates in more than one way:

```go-html-template
{{ upper "BatMan" }} → "BATMAN"
{{ "BatMan" | upper }} → "BATMAN"
```
