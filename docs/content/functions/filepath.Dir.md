---
title: filepath.Dir
description: Returns all but the last element of a file system path, typically the path's directory.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Dir PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Dir` returns all but the last element of file system path `PATH`, typically the path's directory.
After dropping the final element, `filepath.Dir` calls `filepath.Clean` on the path and trailing slashes are removed.
If the path is empty, `filepath.Dir` returns "`.`".
If the path consists entirely of separators, `filepath.Dir` returns a single separator.
The returned path does not end in a separator unless it is the root directory.

    {{ filepath.Dir "abc/" }} → "abc"
    {{ filepath.Dir "///abc/" }} → "/abc"
    {{ filepath.Dir "///" }} → "/"
    {{ filepath.Dir "/abc/def/../../" }} → "/"
    {{ filepath.Dir "" }} → "."
