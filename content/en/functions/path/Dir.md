---
title: path.Dir
description: Replaces path separators with slashes (/) and returns all but the last element of the given path.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/path/Base
    - functions/path/BaseName
    - functions/path/Clean
    - functions/path/Ext
    - functions/path/Join
    - functions/path/Split
  returnType: string
  signatures: [path.Dir PATH]
aliases: [/functions/path.dir]
---

```go-html-template
{{ path.Dir "a/news.html" }} → a
{{ path.Dir "news.html" }} → .
{{ path.Dir "a/b/c" }} → a/b
{{ path.Dir "/a/b/c" }} → /a/b
{{ path.Dir "/a/b/c/" }} → /a/b/c
{{ path.Dir "" }} → .
```
