---
title: path.Dir
description: Dir returns all but the last element of a path.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [path.Dir PATH]
relatedFunctions:
  - path.Base
  - path.BaseName
  - path.Clean
  - path.Dir
  - path.Ext
  - path.Join
  - path.Split
aliases: [/functions/path.dir]
---

`path.Dir` returns all but the last element of `PATH`, typically `PATH`'s directory.

The returned path will never end in a slash.
If `PATH` is empty, `.` is returned.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```go-html-template
{{ path.Dir "a/news.html" }} → "a"
{{ path.Dir "news.html" }} → "."
{{ path.Dir "a/b/c" }} → "a/b"
{{ path.Dir "/x/y/z" }} → "/x/y"
```
