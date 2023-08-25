---
title: countrunes
description: Determines the number of runes in a string excluding any whitespace.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [counting, word count]
signature:
  - "countrunes INPUT"
  - "strings.CountRunes INPUT"
relatedfuncs: []
---

In contrast with `countwords` function, which counts every word in a string, the `countrunes` function determines the number of runes in the content and excludes any whitespace. This has specific utility if you are dealing with CJK-like languages.

```go-html-template
{{ "Hello, 世界" | countrunes }}
<!-- outputs a content length of 8 runes. -->
```

[pagevars]: /variables/page/
