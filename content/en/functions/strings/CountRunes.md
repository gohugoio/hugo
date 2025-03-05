---
title: strings.CountRunes
description: Returns the number of runes in the given string excluding whitespace.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [countrunes]
    returnType: int
    signatures: [strings.CountRunes INPUT]
aliases: [/functions/countrunes]
---

In contrast with the [`strings.RuneCount`] function, which counts every rune in a string, `strings.CountRunes` excludes whitespace.

```go-html-template
{{ "Hello, 世界" | strings.CountRunes }} → 8
```

[`strings.RuneCount`]: /functions/strings/runecount/
