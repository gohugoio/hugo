---
title: strings.Replace
description: Returns the given string, replacing all occurrences of OLD with NEW.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [replace]
    returnType: string
    signatures: ['strings.Replace STRING OLD NEW [LIMIT]']
aliases: [/functions/replace]
---

```go-html-template
{{ $s := "Batman and Robin" }}
{{ replace $s "Robin" "Catwoman" }} → Batman and Catwoman
```

Limit the number of replacements using the `LIMIT` argument:

```go-html-template
{{ replace "aabbaabb" "a" "z" 2 }} → zzbbaabb
```
