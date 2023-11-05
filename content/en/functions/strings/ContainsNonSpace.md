---
title: strings.ContainsNonSpace
description: Reports whether a string contains any non-space characters as defined by Unicode’s White Space property.
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

```go-html-template
{{ strings.ContainsNonSpace "\n" }} → false
{{ strings.ContainsNonSpace " " }} → false
{{ strings.ContainsNonSpace "\n abc" }} → true
```

Common white space characters include:

```text
'\t', '\n', '\v', '\f', '\r', ' '
```

See the [Unicode Character Database] for a complete list.

[Unicode Character Database]: https://www.unicode.org/Public/UCD/latest/ucd/PropList.txt
