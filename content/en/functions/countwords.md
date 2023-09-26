---
title: countwords
description: Counts the number of words in a string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - len
  - strings.Count
  - strings.CountRunes
  - strings.CountWords
  - strings.RuneCount
signature:
  - strings.CountWords INPUT
  - countwords INPUT
---

The template function works similar to the [.WordCount page variable][pagevars].

```go-html-template
{{ "Hugo is a static site generator." | countwords }}
<!-- outputs a content length of 6 words.  -->
```


[pagevars]: /variables/page/
