---
title: os.FileExists
linkTitle: fileExists
description: Reports whether the file or directory exists.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [fileExists]
  returnType: bool
  signatures: [os.FileExists PATH]
relatedFunctions:
  - os.FileExists
  - os.Getenv
  - os.ReadDir
  - os.ReadFile
  - os.Stat
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
{{ os.FileExists "content" }} → true
{{ os.FileExists "content/news" }} → true
{{ os.FileExists "content/news/article-1" }} → false
{{ os.FileExists "content/news/article-1.md" }} → true
{{ os.FileExists "news" }} → true
{{ os.FileExists "news/article-1" }} → false
{{ os.FileExists "news/article-1.md" }} → true
```
