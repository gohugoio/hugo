---
title: chomp
toc: true
description: Removes any trailing newline characters.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [trim]
signature:
  - "chomp INPUT"
  - "strings.Chomp INPUT"
relatedfuncs: [truncate]
---

Useful in a pipeline to remove newlines added by other processing (e.g., [`markdownify`](/functions/markdownify/)).

```go-html-template
{{ chomp "<p>Blockhead</p>\n"}} → "<p>Blockhead</p>"
```
