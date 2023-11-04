---
title: strings.CountRunes
linkTitle: countrunes
description: Returns the number of runes in a string excluding whitespace.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [countrunes]
  returnType: int
  signatures: [strings.CountRunes INPUT]
relatedFunctions:
  - len
  - strings.Count
  - strings.CountRunes
  - strings.CountWords
  - strings.RuneCount
aliases: [/functions/countrunes]
---

In contrast with the [`strings.RuneCount`] function, which counts every rune in a string, `strings.CountRunes` excludes whitespace.

```go-html-template
{{ "Hello, 世界" | strings.CountRunes }} → 8
```

[`strings.RuneCount`]: /functions/strings/runecount
