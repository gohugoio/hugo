---
title: strings.TrimPrefix
description: Returns the given string, removing the prefix from the beginning of the string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/strings/Chomp
    - functions/strings/Trim
    - functions/strings/TrimLeft
    - functions/strings/TrimRight
    - functions/strings/TrimSuffix
  returnType: string
  signatures: [strings.TrimPrefix PREFIX STRING]
aliases: [/functions/strings.trimprefix]
---

```go-html-template
{{ strings.TrimPrefix "a" "aabbaa" }} → abbaa
{{ strings.TrimPrefix "aa" "aabbaa" }} → bbaa
{{ strings.TrimPrefix "aaa" "aabbaa" }} → aabbaa
```
