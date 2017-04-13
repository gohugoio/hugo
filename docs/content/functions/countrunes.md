---
title: countrunes
linktitle: countrunes
description: Determines the number of runes in a string and excludes any whitespace.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [counting, word count]
ns:
signature: ["countrunes INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: [/functions/countrunes/,/functions/countwords/]
---

In contrast with `countwords` function, which counts every word in a string, the `countrunes` function determines the number of runes in the content and excludes any whitespace. This has specific utility if you are dealing with CJK-like languages.

```html
{{ "Hello, 世界" | countrunes }}
<!-- outputs a content length of 8 runes. -->
```

[pagevars]: /variables/page/
