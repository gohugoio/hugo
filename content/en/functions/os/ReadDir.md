---
title: os.ReadDir
linkTitle: readDir
description: Returns an array of FileInfo structures sorted by file name, one element for each directory entry.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [readDir]
  returnType: FileInfo
  signatures: [os.ReadDir PATH]
relatedFunctions:
  - os.FileExists
  - os.Getenv
  - os.ReadDir
  - os.ReadFile
  - os.Stat
aliases: [/functions/readdir]
---

The `os.ReadDir` function resolves the path relative to the root of your project directory. A leading path separator (`/`) is optional.

With this directory structure:

```text
content/
├── about.md
├── contact.md
└── news/
    ├── article-1.md
    └── article-2.md
```

This template code:

```go-html-template
{{ range os.ReadDir "content" }}
  {{ .Name }} → {{ .IsDir }}
{{ end }}
```

Produces:

```html
about.md → false
contact.md → false
news → true
```

Note that `os.ReadDir` is not recursive.

Details of the `FileInfo` structure are available in the [Go documentation](https://pkg.go.dev/io/fs#FileInfo).

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates](/templates/files).
