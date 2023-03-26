---
title: "fileExists"
linktitle: "fileExists"
date: 2017-08-31T22:38:22+02:00
description: Checks for file or directory existence.
publishdate: 2017-08-31T22:38:22+02:00
lastmod: 2021-11-26
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["os.FileExists PATH","fileExists PATH"]
workson: []
hugoversion:
relatedfuncs: ['os.ReadDir','os.ReadFile','os.Stat']
deprecated: false
aliases: []
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
{{ os.FileExists "content" }} --> true
{{ os.FileExists "content/news" }} --> true
{{ os.FileExists "content/news/article-1" }} --> false
{{ os.FileExists "content/news/article-1.md" }} --> true
{{ os.FileExists "news" }} --> true
{{ os.FileExists "news/article-1" }} --> false
{{ os.FileExists "news/article-1.md" }} --> true
```
