---
title: collections.SymDiff
linkTitle: symdiff
description: Returns the symmetric difference of two collections.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [symdiff]
  returnType: any
  signatures: [COLLECTION | collections.SymDiff COLLECTION]
relatedFunctions:
  - collections.Complement
  - collections.Intersect
  - collections.SymDiff
  - collections.Union
aliases: [/functions/symdiff]
---

Example:

```go-html-template
{{ slice 1 2 3 | symdiff (slice 3 4) }} → [1 2 4]
```

Also see https://en.wikipedia.org/wiki/Symmetric_difference
