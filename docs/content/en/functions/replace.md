---
title: replace
description: Replaces all occurrences of the search string with the replacement string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [replace]
signature: 
  - "replace INPUT OLD NEW [LIMIT]"
  - "strings.Replace INPUT OLD NEW [LIMIT]"
relatedfuncs: [replaceRE]
---

Replace returns a copy of `INPUT` with all occurrences of `OLD` replaced with `NEW`.
The number of replacements can be limited with an optional `LIMIT` parameter.

```
`{{ replace "Batman and Robin" "Robin" "Catwoman" }}`
→ "Batman and Catwoman"

{{ replace "aabbaabb" "a" "z" 2 }} → "zzbbaabb"
```
