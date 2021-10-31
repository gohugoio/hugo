---
title: replaceRE
description: Replaces all occurrences of a regular expression with the replacement pattern.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2020-09-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [regex]
signature: ["strings.ReplaceRE PATTERN REPLACEMENT INPUT [LIMIT]", "replaceRE PATTERN REPLACEMENT INPUT [LIMIT]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`strings.ReplaceRE` returns a copy of `INPUT`, replacing all matches of the regular
expression `PATTERN` with the replacement text `REPLACEMENT`.
The number of replacements can be limited with an optional `LIMIT` parameter.

```
{{ replaceRE "^https?://([^/]+).*" "$1" "http://gohugo.io/docs" }}` → "gohugo.io"
{{ "http://gohugo.io/docs" | replaceRE "^https?://([^/]+).*" "$1" }}` → "gohugo.io"
{{ replaceRE "a+b" "X" "aabbaabbab" 1 }} → "Xbaabbab"
```

{{% note %}}
Hugo uses Go's [Regular Expression package](https://golang.org/pkg/regexp/), which is the same general syntax used by Perl, Python, and other languages but with a few minor differences for those coming from a background in PCRE. For a full syntax listing, see the [GitHub wiki for re2](https://github.com/google/re2/wiki/Syntax).

If you are just learning RegEx, or at least Go's flavor, you can practice pattern matching in the browser at <https://regex101.com/>.
{{% /note %}}
