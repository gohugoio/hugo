---
title: upper
description: Converts all characters in a string to uppercase

categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.FirstUpper
  - strings.Title
  - strings.ToLower
  - strings.ToUpper
signature:
  - strings.ToUpper INPUT
  - upper INPUT
---

Note that `upper` can be applied in your templates in more than one way:

```go-html-template
{{ upper "BatMan" }} → "BATMAN"
{{ "BatMan" | upper }} → "BATMAN"
```
