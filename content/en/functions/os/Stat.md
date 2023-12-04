---
title: os.Stat
description: Returns a FileInfo structure describing a file or directory.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/os/FileExists
    - functions/os/Getenv
    - functions/os/ReadDir
    - functions/os/ReadFile
  returnType: os.FileInfo
  signatures: [os.Stat PATH]
aliases: [/functions/os.stat]
---

The `os.Stat` function attempts to resolve the path relative to the root of your project directory. If a matching file or directory is not found, it will attempt to resolve the path relative to the [`contentDir`](/getting-started/configuration#contentdir). A leading path separator (`/`) is optional.

```go-html-template
{{ $f := os.Stat "README.md" }}
{{ $f.IsDir }}    → false (bool)
{{ $f.ModTime }}  → 2021-11-25 10:06:49.315429236 -0800 PST (time.Time)
{{ $f.Name }}     → README.md (string)
{{ $f.Size }}     → 241 (int64)

{{ $d := os.Stat "content" }}
{{ $d.IsDir }}    → true (bool)
```

Details of the `FileInfo` structure are available in the [Go documentation](https://pkg.go.dev/io/fs#FileInfo).
