---
title: collections.Delimit
description: Loops through any array, slice, or map and returns a string of all the values separated by a delimiter.
categories: []
keywords: []
action:
  aliases: [delimit]
  related:
    - functions/strings/Split
  returnType: string
  signatures: ['collections.Delimit COLLECTION DELIMITER [LAST]']
aliases: [/functions/delimit]
---

Delimit a slice:

```go-html-template
{{ $s := slice "b" "a" "c" }}
{{ delimit $s ", " }} → b, a, c
{{ delimit $s ", " " and "}} → b, a and c
```

Delimit a map:

{{% note %}}
The `delimit` function sorts maps by key, returning the values.
{{% /note %}}

```go-html-template
{{ $m := dict "b" 2 "a" 1 "c" 3 }}
{{ delimit $m ", " }} → 1, 2, 3
{{ delimit $m ", " " and "}} → 1, 2 and 3
```
