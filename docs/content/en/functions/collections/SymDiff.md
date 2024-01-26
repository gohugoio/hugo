---
title: collections.SymDiff
description: Returns the symmetric difference of two collections.
categories: []
keywords: []
action:
  aliases: [symdiff]
  related:
    - functions/collections/Complement
    - functions/collections/Intersect
    - functions/collections/SymDiff
    - functions/collections/Union
  returnType: any
  signatures: [COLLECTION | collections.SymDiff COLLECTION]
aliases: [/functions/symdiff]
---

Example:

```go-html-template
{{ slice 1 2 3 | symdiff (slice 3 4) }} → [1 2 4]
```

Also see <https://en.wikipedia.org/wiki/Symmetric_difference>.
