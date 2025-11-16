---
title: collections.SymDiff
description: Returns the symmetric difference of two collections.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [symdiff]
    returnType: any
    signatures: [COLLECTION | collections.SymDiff COLLECTION]
aliases: [/functions/symdiff]
---

Example:

```go-html-template
{{ slice 1 2 3 | symdiff (slice 3 4) }} → [1 2 4]
```

Also see <https://en.wikipedia.org/wiki/Symmetric_difference>.
