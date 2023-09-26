---
title: uniq
description: Takes in a slice or array and returns a slice with duplicate elements removed.
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
  - collections.Uniq COLLECTION
  - uniq COLLECTION
---


```go-html-template
{{ slice 1 3 2 1 | uniq }} → [1 3 2]
```
