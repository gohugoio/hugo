---
title: replace
description: Replaces all occurrences of the search string with the replacement string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.FindRE
  - strings.FindRESubmatch
  - strings.Replace
  - strings.ReplaceRE
signature: 
  - strings.Replace INPUT OLD NEW [LIMIT]
  - replace INPUT OLD NEW [LIMIT]
---

Replace returns a copy of `INPUT` with all occurrences of `OLD` replaced with `NEW`.
The number of replacements can be limited with an optional `LIMIT` argument.

```
`{{ replace "Batman and Robin" "Robin" "Catwoman" }}`
→ "Batman and Catwoman"

{{ replace "aabbaabb" "a" "z" 2 }} → "zzbbaabb"
```
