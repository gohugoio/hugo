---
title: strings.Replace
linkTitle: replace
description: Replaces all occurrences of the search string with the replacement string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [replace]
  returnType: string
  signatures: ['strings.Replace INPUT OLD NEW [LIMIT]']
relatedFunctions:
  - strings.FindRE
  - strings.FindRESubmatch
  - strings.Replace
  - strings.ReplaceRE
aliases: [/functions/replace]
---

Replace returns a copy of `INPUT` with all occurrences of `OLD` replaced with `NEW`.
The number of replacements can be limited with an optional `LIMIT` argument.

```
{{ replace "Batman and Robin" "Robin" "Catwoman" }}
→ "Batman and Catwoman"

{{ replace "aabbaabb" "a" "z" 2 }} → "zzbbaabb"
```
