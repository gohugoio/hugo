---
title: path.BaseName
description: Replaces path separators with slashes (`/`) and returns the last element of the given path, removing the extension if present.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [path.BaseName PATH]
aliases: [/functions/path.basename]
---

```go-html-template
{{ path.BaseName "a/news.html" }} → news
{{ path.BaseName "news.html" }} → news
{{ path.BaseName "a/b/c" }} → c
{{ path.BaseName "/x/y/z/" }} → z
{{ path.BaseName "" }} → .
```
