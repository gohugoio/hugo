---
title: strings.Contains
description: Reports whether the string contains a substring.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: bool
  signatures: [strings.Contains STRING SUBSTRING]
relatedFunctions:
  - strings.Contains
  - strings.ContainsAny
  - strings.ContainsNonSpace
  - strings.HasPrefix
  - strings.HasSuffix
aliases: [/functions/strings.contains]
---

```go-html-template
{{ strings.Contains "Hugo" "go" }} → true
```
The check is case sensitive: 

```go-html-template
{{ strings.Contains "Hugo" "Go" }} → false
```
