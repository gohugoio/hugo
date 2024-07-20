---
title: strings.CountRunes
description: Returns the number of runes in the given string excluding whitespace.
categories: []
keywords: []
action:
  aliases: [countrunes]
  related:
    - functions/go-template/len
    - functions/strings/Count
    - functions/strings/CountWords
    - functions/strings/RuneCount
  returnType: int
  signatures: [strings.CountRunes INPUT]
aliases: [/functions/countrunes]
---

In contrast with the [`strings.RuneCount`] function, which counts every rune in a string, `strings.CountRunes` excludes whitespace.

```go-html-template
{{ "Hello, 世界" | strings.CountRunes }} → 8
```

[`strings.RuneCount`]: /functions/strings/runecount/
