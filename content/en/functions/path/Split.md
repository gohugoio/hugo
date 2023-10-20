---
title: path.Split
description: Split path immediately following the final slash.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: DirFile
  signatures: [path.Split PATH]
relatedFunctions:
  - path.Base
  - path.BaseName
  - path.Clean
  - path.Dir
  - path.Ext
  - path.Join
  - path.Split
aliases: [/functions/path.split]
---

`path.Split` splits `PATH` immediately following the final slash, separating it into a directory and a base component.

The returned values have the property that `PATH` = `DIR`+`BASE`.
If there is no slash in `PATH`, it returns an empty directory and the base is set to `PATH`.

**Note:** On Windows, `PATH` is converted to slash (`/`) separators.

```go-html-template
{{ $dirFile := path.Split "a/news.html" }}
{{ $dirFile.Dir }} → "a/" 
{{ $dirFile.File }} → "news.html"

{{ $dirFile := path.Split "news.html" }}
{{ $dirFile.Dir }} → ""
{{ $dirFile.File }} → "news.html"

{{ $dirFile := path.Split "a/b/c" }}
{{ $dirFile.Dir }} → "a/b/"
{{ $dirFile.File }} → "c"
```
