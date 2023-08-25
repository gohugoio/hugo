---
title: path.Clean
description: Replaces path separators with slashes (`/`) and removes extraneous separators.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [path, clean]
signature: ["path.Clean PATH"]
relatedfuncs: [path.Base, path.BaseName, path.Dir, path.Ext, path.Join, path.Split]
---

`path.Clean` replaces path separators with slashes (`/`) and removes extraneous separators, including trailing separators.

```go-html-template
{{ path.Clean "foo//bar" }} → "foo/bar"
{{ path.Clean "/foo/bar/" }} → "/foo/bar"
```

On a Windows system, if `.File.Path` is `foo\bar.md`, then:

```go-html-template
{{ path.Clean .File.Path }} → "foo/bar.md"
```
