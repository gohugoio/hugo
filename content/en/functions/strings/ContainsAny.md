---
title: strings.ContainsAny
description: Reports whether the given string contains any character within the given set.
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
  signatures: [strings.ContainsAny STRING SET]
aliases: [/functions/strings.containsany]
---

```go-html-template
{{ strings.ContainsAny "Hugo" "gm" }} → true
```

The check is case sensitive:

```go-html-template
{{ strings.ContainsAny "Hugo" "Gm" }} → false
```
