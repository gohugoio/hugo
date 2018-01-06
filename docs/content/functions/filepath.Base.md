---
title: filepath.Base
description: Returns the last element of a path.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Base PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Base` returns the last element of `PATH`.
Trailing path separators are removed before extracting the last element.
If the path is empty, `filepath.Base` returns "`.`".
If the path consists entirely of separators, `filepath.Base` returns a single separator.

    {{ filepath.Base "abc/def.json" }} → "def.json"
    {{ filepath.Base "abc/def/" }} → "def"
    {{ filepath.Base "///" }} → "/"
    {{ filepath.Base "" }} → "."
