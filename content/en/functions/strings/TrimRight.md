---
title: strings.TrimRight
description: Returns the given string, removing trailing characters specified in the cutset.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/strings/Chomp
    - functions/strings/Trim
    - functions/strings/TrimLeft
    - functions/strings/TrimPrefix
    - functions/strings/TrimSuffix
  returnType: string
  signatures: [strings.TrimRight CUTSET STRING]
aliases: [/functions/strings.trimright]
---

```go-html-template
{{ strings.TrimRight "a" "abba" }} → abb
```

The `strings.TrimRight` function converts the arguments to strings if possible:

```go-html-template
{{ strings.TrimRight 54 12345 }} → 123 (string)
{{ strings.TrimRight "eu" true }} → tr
```
