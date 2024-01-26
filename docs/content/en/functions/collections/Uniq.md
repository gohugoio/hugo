---
title: collections.Uniq
description: Returns the given collection, removing duplicate elements.
categories: []
keywords: []
action:
  aliases: [uniq]
  related:
    - functions/collections/Reverse
    - functions/collections/Shuffle
    - functions/collections/Sort
    - functions/collections/Uniq
  returnType: any
  signatures: [collections.Uniq COLLECTION]
aliases: [/functions/uniq]
---

```go-html-template
{{ slice 1 3 2 1 | uniq }} → [1 3 2]
```
