---
title: strings.ContainsNonSpace
description: Reports whether the given string contains any non-space characters as defined by Unicode.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [strings.ContainsNonSpace STRING]
aliases: [/functions/strings.containsnonspace]
---

Whitespace characters include `\t`, `\n`, `\v`, `\f`, `\r`, and characters in the [Unicode Space Separator] category.

[Unicode Space Separator]: https://www.compart.com/en/unicode/category/Zs

```go-html-template
{{ strings.ContainsNonSpace "\n" }} → false
{{ strings.ContainsNonSpace " " }} → false
{{ strings.ContainsNonSpace "\n abc" }} → true
```
