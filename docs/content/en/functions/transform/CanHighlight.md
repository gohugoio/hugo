---
title: transform.CanHighlight
description: Reports whether the given code language is supported by the Chroma highlighter.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: bool
  signatures: [transform.CanHighlight LANGUAGE]
relatedFunctions:
  - transform.CanHighlight
  - transform.Highlight
  - transform.HighlightCodeBlock
---

```go-html-template
{{ transform.CanHighlight "go" }} → true
{{ transform.CanHighlight "klingon" }} → false
```
