---
title: path.Clean
description: Returns the shortest path name equivalent to the given path, after replacing path separators with slashes (`/`).
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [path.Clean PATH]
aliases: [/functions/path.clean]
---

See Go's [`path.Clean`][] documentation for details.

```go-html-template
{{ path.Clean "foo/bar" }} → foo/bar
{{ path.Clean "/foo/bar" }} → /foo/bar
{{ path.Clean "/foo/bar/" }} → /foo/bar
{{ path.Clean "/foo//bar/" }} → /foo/bar
{{ path.Clean "/foo/./bar/" }} → /foo/bar
{{ path.Clean "/foo/../bar/" }} → /bar
{{ path.Clean "/../foo/../bar/" }} → /bar
{{ path.Clean "" }} → .
```

[`path.Clean`]: https://pkg.go.dev/path#Clean
