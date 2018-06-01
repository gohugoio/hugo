---
title: runelen
description: Determines the number of runes in a string.
godocref:
date: 2018-06-01
publishdate: 2018-06-01
lastmod: 2018-06-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [counting, character count, length, rune length, rune count]
signature: ["runelen INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: [/functions/runelen/]
---

In contrast with `countrunes` function, which strips HTML and whitespace before counting runes, `runelen` simply counts all the runes in a string. It relies on the Go [`utf8.RuneCount`] function.

```
{{ "Hello, 世界" | runelen }}
<!-- outputs a content length of 9 runes. -->
```

[`utf8.RuneCount`]: https://golang.org/pkg/unicode/utf8/#RuneCount