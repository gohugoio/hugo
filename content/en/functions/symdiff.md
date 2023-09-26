---
title: symdiff
description: "`collections.SymDiff` (alias `symdiff`) returns the symmetric difference of two collections."
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs:
  - collections.Complement
  - collections.Intersect
  - collections.SymDiff
  - collections.Union
signature: 
  - COLLECTION | collections.SymDiff COLLECTION
  - COLLECTION | symdiff COLLECTION
---

Example:

```go-html-template
{{ slice 1 2 3 | symdiff (slice 3 4) }}
```

The above will print `[1 2 4]`.

Also see https://en.wikipedia.org/wiki/Symmetric_difference
