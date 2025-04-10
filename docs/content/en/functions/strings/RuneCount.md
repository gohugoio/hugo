---
title: strings.RuneCount
description: Returns the number of runes in the given string.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: int
    signatures: [strings.RuneCount INPUT]
aliases: [/functions/strings.runecount]
---

In contrast with the [`strings.CountRunes`] function, which excludes whitespace, `strings.RuneCount` counts every rune in a string.

```go-html-template
{{ "Hello, 世界" | strings.RuneCount }} → 9
```

[`strings.CountRunes`]: /functions/strings/countrunes/
