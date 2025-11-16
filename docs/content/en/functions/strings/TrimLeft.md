---
title: strings.TrimLeft
description: Returns the given string, removing leading characters specified in the cutset.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [strings.TrimLeft CUTSET STRING]
aliases: [/functions/strings.trimleft]
---

```go-html-template
{{ strings.TrimLeft "a" "abba" }} → bba
```

The `strings.TrimLeft` function converts the arguments to strings if possible:

```go-html-template
{{ strings.TrimLeft 21 12345 }} → 345 (string)
{{ strings.TrimLeft "rt" true }} → ue
```
