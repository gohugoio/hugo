---
title: substr
description: Extracts parts of a string from a specified character's position and returns the specified number of characters.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs: []
signature:
  - substr STRING START [LENGTH]
  - strings.Substr STRING START [LENGTH]
---

It normally takes two argument: `start` and `length`. It can also take one argument: `start`, i.e. `length` is omitted, in which case the substring starting from start until the end of the string will be returned.

To extract characters from the end of the string, use a negative start number.

If `length` is given and is negative, that number of characters will be omitted from the end of string.

```go-html-template
{{ substr "abcdef" 0 }} → "abcdef"
{{ substr "abcdef" 1 }} → "bcdef"

{{ substr "abcdef" 0 1 }} → "a"
{{ substr "abcdef" 1 1 }} → "b"

{{ substr "abcdef" 0 -1 }} → "abcde"
{{ substr "abcdef" 1 -1 }} → "bcde"

{{ substr "abcdef" -1 }} → "f"
{{ substr "abcdef" -2 }} → "ef"

{{ substr "abcdef" -1 1 }} → "f"
{{ substr "abcdef" -2 1 }} → "e"

{{ substr "abcdef" -3 -1 }} → "de"
{{ substr "abcdef" -3 -2 }} → "d"
```
