---
title: strings.ReplaceRE
description: Returns a copy of INPUT, replacing all occurrences of a regular expression with a replacement pattern.
categories: []
keywords: []
action:
  aliases: [replaceRE]
  related:
    - functions/strings/FindRE
    - functions/strings/FindRESubmatch
    - functions/strings/Replace
  returnType: string
  signatures: ['strings.ReplaceRE PATTERN REPLACEMENT INPUT [LIMIT]']
aliases: [/functions/replacere]
---

{{% include "functions/_common/regular-expressions.md" %}}

```go-html-template
{{ $s := "a-b--c---d" }}
{{ replaceRE `(-{2,})` "-" $s }} → a-b-c-d
```

Limit the number of replacements using the LIMIT argument:

```go-html-template
{{ $s := "a-b--c---d" }}
{{ replaceRE `(-{2,})` "-" $s 1 }} → a-b-c---d
```

Use `$1`, `$2`, etc. within the replacement string to insert the content of each capturing group within the regular expression:

```go-html-template
{{ $s := "http://gohugo.io/docs" }}
{{ replaceRE "^https?://([^/]+).*" "$1" $s }} → gohugo.io
```

{{% note %}}
You can write and test your regular expression using [regex101.com](https://regex101.com/). Be sure to select the Go flavor before you begin.
{{% /note %}}

[RE2]: https://github.com/google/re2/wiki/Syntax
[string literal]: https://go.dev/ref/spec#String_literals
