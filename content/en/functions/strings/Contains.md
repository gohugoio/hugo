---
title: strings.Contains
description: Reports whether the given string contains the given substring.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/strings/ContainsAny
    - functions/strings/ContainsNonSpace
    - functions/strings/HasPrefix
    - functions/strings/HasSuffix
    - functions/collections/In
  returnType: bool
  signatures: [strings.Contains STRING SUBSTRING]
aliases: [/functions/strings.contains]
---

```go-html-template
{{ strings.Contains "Hugo" "go" }} → true
```

The check is case sensitive:

```go-html-template
{{ strings.Contains "Hugo" "Go" }} → false
```
