---
title: strings.TrimLeft
description: Returns a slice of a given string with all leading characters contained in the cutset removed.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [strings.TrimLeft CUTSET STRING]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/strings.trimleft]
---

Given the string `"abba"`, leading `"a"`'s can be removed a follows:

```go-html-template
{{ strings.TrimLeft "a" "abba" }} → "bba"
```

Numbers can be handled as well:

```go-html-template
{{ strings.TrimLeft 12 1221341221 }} → "341221"
```
