---
title: collections.Uniq
linkTitle: uniq
description: Takes in a slice or array and returns a slice with duplicate elements removed.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [uniq]
  returnType: any
  signatures: [collections.Uniq COLLECTION]
relatedFunctions:
  - collections.Reverse
  - collections.Shuffle
  - collections.Sort
  - collections.Uniq
aliases: [/functions/uniq]
---


```go-html-template
{{ slice 1 3 2 1 | uniq }} → [1 3 2]
```
