---
title: filepath.Separator
description: Returns the OS-specific path separator.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Separator PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Separator` returns the OS-specific path separator.

    Unix:
    {{ filepath.Separator }} → "/"

    Windows:
    {{ filepath.Separator }} → "\"
