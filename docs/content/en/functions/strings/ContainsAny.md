---
title: strings.ContainsAny
description: Reports whether a string contains any character from a given string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: bool
  signatures: [strings.ContainsAny STRING CHARACTERS]
relatedFunctions:
  - strings.Contains
  - strings.ContainsAny
  - strings.ContainsNonSpace
  - strings.HasPrefix
  - strings.HasSuffix
aliases: [/functions/strings.containsany]
---

```go-html-template
{{ strings.ContainsAny "Hugo" "gm" }} → true
```

The check is case sensitive: 

```go-html-template
{{ strings.ContainsAny "Hugo" "Gm" }} → false
```
