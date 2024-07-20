---
title: strings.Substr
description: Returns a substring of the given string, beginning with the start position and ending after the given length.
categories: []
keywords: []
action:
  aliases: [substr]
  related:
    - functions/strings/SliceString
  returnType: string
  signatures: ['strings.Substr STRING [START] [LENGTH]']
aliases: [/functions/substr]
---

The start position is zero-based, where `0` represents the first character of the string. If START is not specified, the substring will begin at position `0`. Specify a negative START position to extract characters from the end of the string. 

If LENGTH is not specified, the substring will include all characters from the START position to the end of the string. If negative, that number of characters will be omitted from the end of string.

```go-html-template
{{ substr "abcdef" 0 }} → abcdef
{{ substr "abcdef" 1 }} → bcdef

{{ substr "abcdef" 0 1 }} → a
{{ substr "abcdef" 1 1 }} → b

{{ substr "abcdef" 0 -1 }} → abcde
{{ substr "abcdef" 1 -1 }} → bcde

{{ substr "abcdef" -1 }} → f
{{ substr "abcdef" -2 }} → ef

{{ substr "abcdef" -1 1 }} → f
{{ substr "abcdef" -2 1 }} → e

{{ substr "abcdef" -3 -1 }} → de
{{ substr "abcdef" -3 -2 }} → d
```
