---
title: shuffle
description: Returns a random permutation of a given array or slice.
keywords: [ordering]
categories: [functions]
menu:
  docs:
    parent: functions
signature: ["shuffle COLLECTION"]
relatedfuncs: [seq]
---


```go-html-template
{{ shuffle (seq 1 2 3) }} → [3 1 2] 
{{ shuffle (slice "a" "b" "c") }} → [b a c] 
```

The result will vary from one build to the next.
