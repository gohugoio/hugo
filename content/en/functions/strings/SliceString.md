---
title: strings.SliceString
description: Returns a substring of the given string, beginning with the start position and ending before the end position.
categories: []
keywords: []
action:
  aliases: [slicestr]
  related:
    - functions/strings/Substr
  returnType: string
  signatures: ['strings.SliceString STRING [START] [END]']
aliases: [/functions/slicestr]
---

The START and END positions are zero-based, where `0` represents the first character of the string. If START is not specified, the substring will begin at position `0`. If END is not specified, the substring will end after the last character.

```go-html-template
{{ slicestr "BatMan" }} → BatMan
{{ slicestr "BatMan" 3 }} → Man
{{ slicestr "BatMan" 0 3 }} → Bat
```

The START and END arguments represent the endpoints of a [half-open interval], a concept that may be difficult to grasp when first encountered. You may find that the [`strings.Substr`] function is easier to understand.

[half-open interval]: /getting-started/glossary/#interval
[`strings.Substr`]: /functions/strings/substr/
