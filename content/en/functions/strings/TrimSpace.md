---
title: strings.TrimSpace
description: Returns the given string, removing leading and trailing whitespace as defined by Unicode.
categories: []
keywords: []
action:
  related:
    - functions/strings/Chomp
    - functions/strings/Trim
    - functions/strings/TrimLeft
    - functions/strings/TrimPrefix
    - functions/strings/TrimRight
    - functions/strings/TrimSuffix
  returnType: string
  signatures: [strings.TrimSpace INPUT]
---

{{< new-in 0.136.3 >}}

Whitespace characters include `\t`, `\n`, `\v`, `\f`, `\r`, and characters in the [Unicode Space Separator] category.

[Unicode Space Separator]: https://www.compart.com/en/unicode/category/Zs

```go-html-template
{{ strings.TrimSpace "\n\r\t   foo   \n\r\t" }} → foo
```
