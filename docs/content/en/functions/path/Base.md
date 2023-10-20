---
title: path.Base
description: Base returns the last element of a path.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [path.Base PATH]
relatedFunctions:
  - path.Base
  - path.BaseName
  - path.Clean
  - path.Dir
  - path.Ext
  - path.Join
  - path.Split
aliases: [/functions/path.base]
---

`path.Base` returns the last element of `PATH`.

If `PATH` is empty, `.` is returned.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```go-html-template
{{ path.Base "a/news.html" }} → "news.html"
{{ path.Base "news.html" }} → "news.html"
{{ path.Base "a/b/c" }} → "c"
{{ path.Base "/x/y/z/" }} → "z"
```
