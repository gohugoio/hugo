---
title: shuffle
description: Returns a random permutation of a given array or slice.
keywords: [ordering]
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs:
  - collections.Reverse
  - collections.Shuffle
  - collections.Sort
  - collections.Uniq
signature:
  - collections.Shuffle COLLECTION
  - shuffle COLLECTION
---


```go-html-template
{{ shuffle (seq 1 2 3) }} → [3 1 2] 
{{ shuffle (slice "a" "b" "c") }} → [b a c] 
```

The result will vary from one build to the next.
