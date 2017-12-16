---
title: filepath.FromSlash
description: Returns the result of replacing each slash with a separator.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.FromSlash PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.FromSlash` returns the result of replacing each slash ("`/`")
character in path with a separator character.
Multiple slashes are replaced by multiple separators.

    {{ filepath.FromSlash "on/unix" }} → "on/unix"
    {{ filepath.FromSlash "on/windows" }} → "on\windows"
