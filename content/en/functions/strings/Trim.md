---
title: strings.Trim
description: Returns the given string, removing leading and trailing characters specified in the cutset.
categories: []
keywords: []
action:
  aliases: [trim]
  related:
    - functions/strings/Chomp
    - functions/strings/TrimSpace
    - functions/strings/TrimLeft
    - functions/strings/TrimPrefix
    - functions/strings/TrimRight
    - functions/strings/TrimSuffix
  returnType: string
  signatures: [strings.Trim INPUT CUTSET]
aliases: [/functions/trim]
---

```go-html-template
{{ trim "++foo--" "+-" }} → foo
```
