---
title: "symdiff"
description: "`collections.SymDiff` (alias `symdiff`) returns the symmetric difference of two collections."
date: 2018-11-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [collections,intersect,union, complement]
signature: ["COLLECTION | symdiff COLLECTION" ]
hugoversion: "0.51"
aliases: []
---

Example:

```go-html-template
{{ slice 1 2 3 | symdiff (slice 3 4) }}
```

The above will print `[1 2 4]`.

Also see https://en.wikipedia.org/wiki/Symmetric_difference





