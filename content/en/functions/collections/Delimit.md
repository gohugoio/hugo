---
title: collections.Delimit
description: Returns a string by joining the values of the given slice or map with a delimiter.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [delimit]
    returnType: string
    signatures: ['collections.Delimit SLICE|MAP DELIMITER [LAST]']
aliases: [/functions/delimit]
---

Delimit a slice:

```go-html-template
{{ $s := slice "b" "a" "c" }}
{{ delimit $s ", " }} → b, a, c
{{ delimit $s ", " " and "}} → b, a and c
```

Delimit a map:

> [!note]
> The `delimit` function sorts maps by key, returning the values.

```go-html-template
{{ $m := dict "b" 2 "a" 1 "c" 3 }}
{{ delimit $m ", " }} → 1, 2, 3
{{ delimit $m ", " " and "}} → 1, 2 and 3
```
