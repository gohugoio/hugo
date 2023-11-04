---
title: os.ReadFile
linkTitle: readFile
description: Returns the contents of a file.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [readFile]
  returnType: string
  signatures: [os.ReadFile PATH]
relatedFunctions:
  - os.FileExists
  - os.Getenv
  - os.ReadDir
  - os.ReadFile
  - os.Stat
aliases: [/functions/readfile]
---

The `os.ReadFile` function attempts to resolve the path relative to the root of your project directory. If a matching file is not found, it will attempt to resolve the path relative to the [`contentDir`](/getting-started/configuration#contentdir). A leading path separator (`/`) is optional.

With a file named README.md in the root of your project directory:

```text
This is **bold** text.
```

This template code:

```go-html-template
{{ os.ReadFile "README.md" }}
```

Produces:

```html
This is **bold** text.
```

Note that `os.ReadFile` returns raw (uninterpreted) content.

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates](/templates/files).
