---
title: path.Dir
description: Dir returns all but the last element of a path.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: path
relatedFuncs:
  - path.Base
  - path.BaseName
  - path.Clean
  - path.Dir
  - path.Ext
  - path.Join
  - path.Split
signature:
  - path.Dir PATH
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
