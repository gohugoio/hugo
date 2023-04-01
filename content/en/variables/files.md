---
title: File Variables
description: "Use File variables to access file-related values for each page that is backed by a file."
categories: [variables and params]
keywords: [files]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 40
toc: true
weight: 40
aliases: [/variables/file-variables/]
---
## Variables

{{% note %}}
The path separators (slash or backslash) in `.File.Path`, `.File.Dir`, and `.File.Filename` depend on the operating system.
{{% /note %}}

.File.Path
: (`string`) The file path, relative to the `content` directory.

.File.Dir
: (`string`) The file path, excluding the file name, relative to the `content` directory.

.File.LogicalName
: (`string`) The file name.

.File.BaseFileName
: (`string`) The file name, excluding the extension.

.File.TranslationBaseName
: (`string`) The file name, excluding the extension and language identifier.

.File.Ext
: (`string`) The file extension.

.File.Lang
: (`string`) The language associated with the given file.


.File.ContentBaseName
: (`string`) If the page is a branch or leaf bundle, the name of the containing directory, else the `.TranslationBaseName`.

.File.Filename
: (`string`) The absolute file path.

.File.UniqueID
: (`string`) The MD5 hash of `.File.Path`.

## Examples

```text
content/
├── news/
│   ├── b/
│   │   ├── index.de.md   <-- leaf bundle
│   │   └── index.en.md   <-- leaf bundle
│   ├── a.de.md           <-- regular content
│   ├── a.en.md           <-- regular content
│   ├── _index.de.md      <-- branch bundle
│   └── _index.en.md      <-- branch bundle
├── _index.de.md
└── _index.en.md
```

With the content structure above, the `.File` objects for the English pages contain the following properties:

&nbsp;|regular content|leaf bundle|branch bundle
:--|:--|:--|:--
Path|news/a.en.md|news/b/index.en.md|news/_index.en.md
Dir|news/|news/b/|news/
LogicalName|a.en.md|index.en.md|_index.en.md
BaseFileName|a.en|index.en|_index.en
TranslationBaseName|a|index|_index
Ext|md|md|md
Lang|en|en|en
ContentBaseName|a|b|news
Filename|/home/user/...|/home/user/...|/home/user/...
UniqueID|15be14b...|186868f...|7d9159d...

## Defensive coding

Some of the pages on a site may not be backed by a file. For example:

- Top level section pages
- Taxonomy pages
- Term pages

Without a backing file, Hugo will throw a warning if you attempt to access a `.File` property. For example:

```text
WARN .File.ContentBaseName on zero object. Wrap it in if or with...
```

To code defensively:

```go-html-template
{{ with .File }}
  {{ .ContentBaseName }}
{{ end }}
```
