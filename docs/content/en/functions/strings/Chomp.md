---
title: chomp
linkTitle: chomp
description: Removes any trailing newline characters.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [chomp]
  returnType: any
  signatures: [strings.Chomp STRING]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/chomp]
---

If the argument is of type template.HTML, returns template.HTML, else returns a string.


Useful in a pipeline to remove newlines added by other processing (e.g., [`markdownify`](/functions/transform/markdownify)).

```go-html-template
{{ chomp "<p>Blockhead</p>\n" }} → "<p>Blockhead</p>"
```
