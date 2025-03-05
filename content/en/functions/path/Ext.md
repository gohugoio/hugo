---
title: path.Ext
description: Replaces path separators with slashes (`/`) and returns the file name extension of the given path.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [path.Ext PATH]
aliases: [/functions/path.ext]
---

The extension is the suffix beginning at the final dot in the final slash-separated element of path; it is empty if there is no dot.

```go-html-template
{{ path.Ext "a/b/c/news.html" }} → .html
```
