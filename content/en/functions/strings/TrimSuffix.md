---
title: strings.TrimSuffix
description: Returns a given string s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [strings.TrimSuffix SUFFIX STRING]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/strings.trimsuffix]
---

Given the string `"aabbaa"`, the specified suffix is only removed if `"aabbaa"` ends with it:

```go-html-template
{{ strings.TrimSuffix "a" "aabbaa" }} → "aabba"
{{ strings.TrimSuffix "aa" "aabbaa" }} → "aabb"
{{ strings.TrimSuffix "aaa" "aabbaa" }} → "aabbaa"
```
