---
title: countrunes
description: Determines the number of runes in a string excluding any whitespace.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [counting, word count]
signature: ["countrunes INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

In contrast with `countwords` function, which counts every word in a string, the `countrunes` function determines the number of runes in the content and excludes any whitespace. This has specific utility if you are dealing with CJK-like languages.

```
{{ "Hello, 世界" | countrunes }}
<!-- outputs a content length of 8 runes. -->
```

[pagevars]: /variables/page/
