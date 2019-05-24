---
title: path.Dir
description: Dir returns all but the last element of a path.
godocref:
date: 2018-11-28
publishdate: 2018-11-28
lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [path, dir]
signature: ["path.Dir PATH"]
workson: []
hugoversion: "0.40"
relatedfuncs: [path.Base, path.Ext, path.Split]
deprecated: false
---

`path.Dir` returns all but the last element of `PATH`, typically `PATH`'s directory.

The returned path will never end in a slash.
If `PATH` is empty, `.` is returned.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```
{{ path.Dir "a/news.html" }} → "a"
{{ path.Dir "news.html" }} → "."
{{ path.Dir "a/b/c" }} → "a/b"
{{ path.Dir "/x/y/z" }} → "/x/y"
```
