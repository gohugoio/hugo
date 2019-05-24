---
title: path.Ext
description: Ext returns the file name extension of a path.
godocref:
date: 2018-11-28
publishdate: 2018-11-28
lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [path, ext, extension]
signature: ["path.Ext PATH"]
workson: []
hugoversion: "0.40"
relatedfuncs: [path.Base, path.Dir, path.Split]
deprecated: false
---

`path.Ext` returns the file name extension `PATH`.

The extension is the suffix beginning at the final dot in the final slash-separated element `PATH`;
it is empty if there is no dot.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```
{{ path.Ext "a/b/c/news.html" }} → ".html"
```
