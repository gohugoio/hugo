---
title: countwords
description: Counts the number of words in a string.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [counting, word count]
signature: ["countwords INPUT"]
workson: []
hugoversion:
relatedfuncs: [countrunes]
deprecated: false
---

The template function works similar to the [.WordCount page variable][pagevars].

```
{{ "Hugo is a static site generator." | countwords }}
<!-- outputs a content length of 6 words.  -->
```


[pagevars]: /variables/page/
