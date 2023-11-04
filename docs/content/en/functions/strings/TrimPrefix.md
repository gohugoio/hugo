---
title: strings.TrimPrefix
description: Returns a given string s without the provided leading prefix string. If s doesn't start with prefix, s is returned unchanged.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [strings.TrimPrefix PREFIX STRING]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/strings.trimprefix]
---

Given the string `"aabbaa"`, the specified prefix is only removed if `"aabbaa"` starts with it:

```go-html-template
{{ strings.TrimPrefix "a" "aabbaa" }} → "abbaa"
{{ strings.TrimPrefix "aa" "aabbaa" }} → "bbaa"
{{ strings.TrimPrefix "aaa" "aabbaa" }} → "aabbaa"
```
