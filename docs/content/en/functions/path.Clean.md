---
title: path.Clean
description: Replaces path separators with slashes (`/`) and removes extraneous separators.
date: 2021-10-08
# publishdate: 2018-11-28
# lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [path]
signature: ["path.Clean PATH"]
---

`path.Clean` replaces path separators with slashes (`/`) and removes extraneous separators, including trailing separators.

```
{{ path.Clean "foo//bar" }} → "foo/bar"
{{ path.Clean "/foo/bar/" }} → "/foo/bar"
```

On a Windows system, if `.File.Path` is `foo\bar.md`, then:

```
{{ path.Clean .File.Path }} → "foo/bar.md"
```
