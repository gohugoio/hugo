---
title: filepath.ToSlash
description: Returns the result of replacing each separator with a slash.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.ToSlash PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.ToSlash` returns the result of replacing each separator character in path with
a slash ("`/`") character.
Multiple separators are replaced by multiple slashes.

    {{ filepath.ToSlash "on/unix" }} → "on/unix"
    {{ filepath.ToSlash "on\windows" }} → "on/windows"
