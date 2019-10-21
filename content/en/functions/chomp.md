---
title: chomp
description: Removes any trailing newline characters.
godocref: Removes any trailing newline characters.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [trim]
signature: ["chomp INPUT"]
workson: []
hugoversion:
relatedfuncs: [truncate]
deprecated: false
---

Useful in a pipeline to remove newlines added by other processing (e.g., [`markdownify`](/functions/markdownify/)).

```
{{chomp "<p>Blockhead</p>\n"}} → "<p>Blockhead</p>"
```
