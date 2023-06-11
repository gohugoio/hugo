---
title: eq
description: Returns the boolean truth of arg1 == arg2 || arg1 == arg3.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [comparison,operators,logic]
signature: ["eq ARG1 ARG2 [ARG...]"]
relatedfuncs: []
---

```go-html-template
{{ eq 1 1 }} → true
{{ eq 1 2 }} → false

{{ eq 1 1 1 }} → true
{{ eq 1 1 2 }} → true
{{ eq 1 2 1 }} → true
{{ eq 1 2 2 }} → false
```
