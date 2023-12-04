---
title: collections.Shuffle
description: Returns a random permutation of a given array or slice.
categories: []
keywords: []
action:
  aliases: [shuffle]
  returnType: any
  signatures: [collections.Shuffle COLLECTION]
related:
  - collections.Reverse
  - collections.Shuffle
  - collections.Sort
  - collections.Uniq
aliases: [/functions/shuffle]
---

```go-html-template
{{ shuffle (seq 1 2 3) }} → [3 1 2] 
{{ shuffle (slice "a" "b" "c") }} → [b a c] 
```

The result will vary from one build to the next.
