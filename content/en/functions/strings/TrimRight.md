---
title: strings.TrimRight
description: Returns a slice of a given string with all trailing characters contained in the cutset removed.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [strings.TrimRight CUTSET STRING]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/strings.trimright]
---

Given the string `"abba"`, trailing `"a"`'s can be removed a follows:

```go-html-template
{{ strings.TrimRight "a" "abba" }} → "abb"
```

Numbers can be handled as well:

```go-html-template
{{ strings.TrimRight 12 1221341221 }} → "122134"
```
