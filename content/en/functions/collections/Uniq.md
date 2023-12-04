---
title: collections.Uniq
description: Returns the given collection, removing duplicate elements.
categories: []
keywords: []
action:
  aliases: [uniq]
  returnType: any
  signatures: [collections.Uniq COLLECTION]
related:
  - collections.Reverse
  - collections.Shuffle
  - collections.Sort
  - collections.Uniq
aliases: [/functions/uniq]
---

```go-html-template
{{ slice 1 3 2 1 | uniq }} → [1 3 2]
```
