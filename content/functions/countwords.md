---
title: countwords
linktitle: countwords
description: Counts the number of words in a string that has been passed to it.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [counting, word count]
ns:
signature: []
workson: []
hugoversion:
relatedfuncs: [countrunes]
deprecated: false
aliases: [/functions/countrunes/,/functions/countwords/]
---

`countwords` tries to convert the passed content to a string and counts each word in it. The template function works similar to the [.WordCount page variable][pagevars].

```html
{{ "Hugo is a static site generator." | countwords }}
<!-- outputs a content length of 6 words.  -->
```


[pagevars]: /variables/page/