---
title: transform.CanHighlight
description: Reports whether the given code language is supported by the Chroma highlighter.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    related:
      - functions/transform/Highlight
      - functions/transform/HighlightCodeBlock
    returnType: bool
    signatures: [transform.CanHighlight LANGUAGE]
---

```go-html-template
{{ transform.CanHighlight "go" }} → true
{{ transform.CanHighlight "klingon" }} → false
```
