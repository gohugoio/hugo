---
title: path.Ext
description: Ext returns the file name extension of a path.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [path, ext, extension]
signature: ["path.Ext PATH"]
relatedfuncs: [path.Base, path.BaseName, path.Clean, path.Dir, path.Join, path.Split]
---

`path.Ext` returns the file name extension `PATH`.

The extension is the suffix beginning at the final dot in the final slash-separated element `PATH`;
it is empty if there is no dot.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```go-html-template
{{ path.Ext "a/b/c/news.html" }} → ".html"
```
