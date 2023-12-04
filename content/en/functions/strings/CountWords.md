---
title: strings.CountWords
description: Returns the number of words in the given string.
categories: []
keywords: []
action:
  aliases: [countwords]
  related:
    - functions/go-template/len
    - functions/strings/Count
    - functions/strings/CountRunes
    - functions/strings/RuneCount
  returnType: int
  signatures: [strings.CountWords INPUT]
aliases: [/functions/countwords]
---

```go-html-template
{{ "Hugo is a static site generator." | countwords }} → 6
```
