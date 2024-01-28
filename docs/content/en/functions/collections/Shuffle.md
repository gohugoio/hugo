---
title: collections.Shuffle
description: Returns a random permutation of a given array or slice.
categories: []
keywords: []
action:
  aliases: [shuffle]
  related:
    - functions/collections/Reverse
    - functions/collections/Sort
    - functions/collections/Uniq
  returnType: any
  signatures: [collections.Shuffle COLLECTION]
aliases: [/functions/shuffle]
---

```go-html-template
{{ shuffle (seq 1 2 3) }} → [3 1 2] 
{{ shuffle (slice "a" "b" "c") }} → [b a c] 
```

The result will vary from one build to the next.
