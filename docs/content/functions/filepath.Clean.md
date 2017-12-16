---
title: filepath.Clean
description: Returns the shortest path name equivalent of a given path.
godocref:
date: 2017-12-15
publishdate: 2017-12-15
lastmod: 2017-12-15
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filepath, base]
signature: ["filepath.Clean PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`filepath.Clean` returns the shortest path name equivalent to `PATH` by purely lexical processing.
It applies the following rules iteratively until no further processing can be done:

1. Replace multiple `filepath.Separator` elements with a single one.
2. Eliminate each "`.`" path name element (the current directory).
3. Eliminate each inner "`..`" path name element (the parent directory)
   along with the non-"`..`" element that precedes it.
4. Eliminate "`..`" elements that begin a rooted path:
   that is, replace "`/..`" by "`/`" at the beginning of a path,
   assuming Separator is "`/`".

The returned path ends in a slash only if it represents a root directory,
such as "`/`" on Unix or "`C:\`" on Windows.

Finally, any occurrences of slash are replaced by `filepath.Separator`.

If the result of this process is an empty string, `filepath.Clean` returns the string "`.`".

    {{ filepath.Clean "abc/" }} → "abc"
    {{ filepath.Clean "//./abc/def" }} → "/abc/def"
    {{ filepath.Clean "" }} → "."
