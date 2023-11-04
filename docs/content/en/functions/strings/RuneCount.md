---
title: strings.RuneCount
description: Returns the number of runes in a string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: int
  signatures: [strings.RuneCount INPUT]
relatedFunctions:
  - len
  - strings.Count
  - strings.CountRunes
  - strings.CountWords
  - strings.RuneCount
aliases: [/functions/strings.runecount]
---

In contrast with the [`strings.CountRunes`] function, which excludes whitespace, `strings.RuneCount` counts every rune in a string.

```go-html-template
{{ "Hello, 世界" | strings.RuneCount }} → 9
```

[`strings.CountRunes`]: /functions/strings/countrunes
