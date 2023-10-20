---
title: collections.NewScratch
linkTitle: newScratch
description: Creates a new Scratch which can be used to store values in a thread safe way.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [newScratch]
  returnType: Scratch
  signatures: [collections.NewScratch ]
relatedFunctions: []
---

```go-html-template
{{ $scratch := newScratch }}
{{ $scratch.Add "b" 2 }}
{{ $scratch.Add "b" 2 }}
{{ $scratch.Get "b" }} → 4
```
