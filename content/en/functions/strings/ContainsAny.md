---
title: strings.ContainsAny
description: Reports whether a string contains any character from a given string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/strings/Contains
    - functions/strings/ContainsNonSpace
    - functions/strings/HasPrefix
    - functions/strings/HasSuffix
    - functions/collections/In
  returnType: bool
  signatures: [strings.ContainsAny STRING CHARACTERS]
aliases: [/functions/strings.containsany]
---

```go-html-template
{{ strings.ContainsAny "Hugo" "gm" }} → true
```

The check is case sensitive:

```go-html-template
{{ strings.ContainsAny "Hugo" "Gm" }} → false
```
