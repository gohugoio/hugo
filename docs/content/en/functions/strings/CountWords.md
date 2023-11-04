---
title: strings.CountWords
linkTitle: countwords
description: Counts the number of words in a string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [countwords]
  returnType: int
  signatures: [strings.CountWords INPUT]
relatedFunctions:
  - len
  - strings.Count
  - strings.CountRunes
  - strings.CountWords
  - strings.RuneCount
aliases: [/functions/countwords]
---

The template function works similar to the [.WordCount page variable][pagevars].

```go-html-template
{{ "Hugo is a static site generator." | countwords }}
<!-- outputs a content length of 6 words.  -->
```


[pagevars]: /variables/page/
