---
title: strings.SliceString
linkTitle: slicestr
description: Creates a slice of a half-open range, including start and end indices.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [slicestr]
  returnType: string
  signatures: ['strings.SliceString STRING START [END]']
relatedFunctions: []
aliases: [/functions/slicestr]
---

For example, 1 and 4 creates a slice including elements 1 through 3.
The `end` index can be omitted; it defaults to the string's length.

```go-html-template
{{ slicestr "BatMan" 3 }}` → "Man"
{{ slicestr "BatMan" 0 3 }}` → "Bat"
```
