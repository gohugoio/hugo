---
title: os.FileExists
description: Reports whether the file or directory exists.
categories: []
keywords: []
action:
  aliases: [fileExists]
  related:
    - functions/os/Getenv
    - functions/os/ReadDir
    - functions/os/ReadFile
    - functions/os/Stat
  returnType: bool
  signatures: [os.FileExists PATH]
aliases: [/functions/fileexists]
---

The `os.FileExists` function attempts to resolve the path relative to the root of your project directory. If a matching file or directory is not found, it will attempt to resolve the path relative to the [`contentDir`](/getting-started/configuration#contentdir). A leading path separator (`/`) is optional.

With this directory structure:

```text
content/
├── about.md
├── contact.md
└── news/
    ├── article-1.md
    └── article-2.md
```

The function returns these values:

```go-html-template
{{ fileExists "content" }} → true
{{ fileExists "content/news" }} → true
{{ fileExists "content/news/article-1" }} → false
{{ fileExists "content/news/article-1.md" }} → true
{{ fileExists "news" }} → true
{{ fileExists "news/article-1" }} → false
{{ fileExists "news/article-1.md" }} → true
```
