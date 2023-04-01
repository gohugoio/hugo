---
title: readDir
description: Returns an array of FileInfo structures sorted by filename, one element for each directory entry.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [files]
signature: ["os.ReadDir PATH", "readDir PATH"]
relatedfuncs: ['os.FileExists','os.ReadFile','os.Stat']
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
  {{ .Name }} --> {{ .IsDir }}
{{ end }}
```

Produces:

```html
about.md --> false
contact.md --> false
news --> true
```

Note that `os.ReadDir` is not recursive.

Details of the `FileInfo` structure are available in the [Go documentation](https://pkg.go.dev/io/fs#FileInfo).

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates](/templates/files).
