---
title: strings.CountWords
description: Returns the number of words in a string.
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

The template function works similar to the [.WordCount page variable][pagevars].

```go-html-template
{{ "Hugo is a static site generator." | countwords }} → 6
```

[pagevars]: /variables/page/
