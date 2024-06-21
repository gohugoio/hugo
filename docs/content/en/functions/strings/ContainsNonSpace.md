---
title: strings.ContainsNonSpace
description: Reports whether the given string contains any non-space characters as defined by Unicode's White Space property.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/strings/Contains
    - functions/strings/ContainsAny
    - functions/strings/HasPrefix
    - functions/strings/HasSuffix
    - functions/collections/In
  returnType: bool
  signatures: [strings.ContainsNonSpace STRING]
aliases: [/functions/strings.containsnonspace]
---

{{< new-in 0.111.0 >}}

```go-html-template
{{ strings.ContainsNonSpace "\n" }} → false
{{ strings.ContainsNonSpace " " }} → false
{{ strings.ContainsNonSpace "\n abc" }} → true
```

Common whitespace characters include:

```text
'\t', '\n', '\v', '\f', '\r', ' '
```

See the [Unicode Character Database] for a complete list.

[Unicode Character Database]: https://www.unicode.org/Public/UCD/latest/ucd/PropList.txt
