---
title: chomp
toc: true
description: Removes any trailing newline characters.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
signature:
  - chomp STRING
  - strings.Chomp STRING
---

Useful in a pipeline to remove newlines added by other processing (e.g., [`markdownify`](/functions/markdownify/)).

```go-html-template
{{ chomp "<p>Blockhead</p>\n" }} → "<p>Blockhead</p>"
```
