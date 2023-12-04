---
title: path.Base
description: Replaces path separators with slashes (`/`) and returns the last element of the given path. 
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/path/BaseName
    - functions/path/Clean
    - functions/path/Dir
    - functions/path/Ext
    - functions/path/Join
    - functions/path/Split
  returnType: string
  signatures: [path.Base PATH]
aliases: [/functions/path.base]
---

```go-html-template
{{ path.Base "a/news.html" }} → news.html
{{ path.Base "news.html" }} → news.html
{{ path.Base "a/b/c" }} → c
{{ path.Base "/x/y/z/" }} → z
{{ path.Base "" }} → .
```
