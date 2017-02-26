---
title: countrunes and countwords
linktitle: countrunes and countwords
description: countrunes and countwords both serve as a means to quantify the total the length of your content.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [counting, word count]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: [/functions/countrunes/,/functions/countwords/]
---

`countwords` tries to convert the passed content to a string and counts each word in it. The template function works similar to the [.WordCount page variable][pagevars].

```html
{{ "Hugo is a static site generator." | countwords }}
<!-- outputs a content length of 6 words.  -->
```


In contrast with counting every word, the `countrunes` function determines the number of runes in the content and excludes any whitespace. This has specific utility if you are dealing with CJK-like languages.

```html
{{ "Hello, 世界" | countrunes }}
<!-- outputs a content length of 8 runes. -->
```

[pagevars]: /variables/page-variables/