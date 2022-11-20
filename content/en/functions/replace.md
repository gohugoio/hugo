---
title: replace
description: Replaces all occurrences of the search string with the replacement string.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2020-09-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [replace]
signature: ["strings.Replace INPUT OLD NEW [LIMIT]", "replace INPUT OLD NEW [LIMIT]"]
workson: []
hugoversion:
relatedfuncs: [replaceRE]
deprecated: false
aliases: []
---

Replace returns a copy of `INPUT` with all occurrences of `OLD` replaced with `NEW`.
The number of replacements can be limited with an optional `LIMIT` parameter.

```
`{{ replace "Batman and Robin" "Robin" "Catwoman" }}`
→ "Batman and Catwoman"

{{ replace "aabbaabb" "a" "z" 2 }} → "zzbbaabb"
```
