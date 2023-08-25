---
title: replaceRE
description: Returns a string, replacing all occurrences of a regular expression with a replacement pattern.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [regex]
signature:
  - "replaceRE PATTERN REPLACEMENT INPUT [LIMIT]"
  - "strings.ReplaceRE PATTERN REPLACEMENT INPUT [LIMIT]"
relatedfuncs: [findRE, FindRESubmatch, replace]
---
By default, `replaceRE` replaces all matches. You can limit the number of matches with an optional LIMIT parameter.

{{% readfile file="/functions/common/regular-expressions.md" %}}

This example replaces two or more consecutive hyphens with a single hyphen:

```go-html-template
{{ $s := "a-b--c---d" }}
{{ replaceRE `(-{2,})` "-" $s }} → a-b-c-d
```

To limit the number of replacements to one:

```go-html-template
{{ $s := "a-b--c---d" }}
{{ replaceRE `(-{2,})` "-" $s 1 }} → a-b-c---d
```

You can use `$1`, `$2`, etc. within the replacement string to insert the groups captured within the regular expression:

```go-html-template
{{ $s := "http://gohugo.io/docs" }}
{{ replaceRE "^https?://([^/]+).*" "$1" $s }} → gohugo.io
```

{{% note %}}
You can write and test your regular expression using [regex101.com](https://regex101.com/). Be sure to select the Go flavor before you begin.
{{% /note %}}

[RE2]: https://github.com/google/re2/wiki/Syntax
[string literal]: https://go.dev/ref/spec#String_literals
