---
title: strings.SliceString
description: Creates a slice of a half-open range, including start and end indices.
categories: []
keywords: []
action:
  aliases: [slicestr]
  related: []
  returnType: string
  signatures: ['strings.SliceString STRING START [END]']
aliases: [/functions/slicestr]
---

For example, 1 and 4 creates a slice including elements 1 through&nbsp;3.
The `end` index can be omitted; it defaults to the string's length.

```go-html-template
{{ slicestr "BatMan" 3 }}` → Man
{{ slicestr "BatMan" 0 3 }}` → Bat
```
