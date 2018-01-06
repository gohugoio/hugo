---
title: filepath.Split
description: Splits path immediately following the final separator.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Split PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Split` splits path immediately following the final Separator,
separating it into a directory and file name component.
If there is no separator in path, `filepath.Split` returns an empty dir and file set to path.
The returned values have the property that `PATH = dir+file`.

    {{ filepath.Split "on/unix" }} → "on/", "unix"
    {{ filepath.Split "on\windows" }} → "on\", "windows"
