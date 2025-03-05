---
title: path.Clean
description: Replaces path separators with slashes (`/`) and returns the shortest path name equivalent to the given path.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [path.Clean PATH]
aliases: [/functions/path.clean]
---

See Go's [`path.Clean`] documentation for details.

[`path.Clean`]: https://pkg.go.dev/path#Clean

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
