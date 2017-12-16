---
title: filepath.Ext
description: Returns the file name extension used by a path.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Ext PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Ext` returns the file name extension used by `PATH`.
The extension is the suffix beginning at the final dot in the final element of path;
it is empty if there is no dot.

    {{ filepath.Ext "abc/def.json" }} → ".json"
    {{ filepath.Ext "abc" }} → ""
