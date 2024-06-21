---
title: File
description: For pages backed by a file, returns file information for the given page.
categories: []
keywords: []
action:
  related: []
  returnType: hugolib.fileInfo
  signatures: [PAGE.File]
toc: true
---

By default, not all pages are backed by a file, including top level [section] pages, [taxonomy] pages, and [term] pages. By definition, you cannot retrieve file information when the file does not exist.

To back one of the pages above with a file, create an _index.md file in the corresponding directory. For example:

```text
content/
└── books/
    ├── _index.md  <-- the top level section page
    ├── book-1.md
    └── book-2.md
```

{{% note %}}
Code defensively by verifying file existence as shown in the examples below.
{{% /note %}}

## Methods

{{% note %}}
The path separators (slash or backslash) in `Path`, `Dir`, and `Filename` depend on the operating system.
{{% /note %}}

###### BaseFileName

(`string`) The file name, excluding the extension.

```go-html-template
{{ with .File }}
  {{ .BaseFileName }}
{{ end }}
```

###### ContentBaseName

(`string`) If the page is a branch or leaf bundle, the name of the containing directory, else the `TranslationBaseName`.

```go-html-template
{{ with .File }}
  {{ .ContentBaseName }}
{{ end }}
```

###### Dir

(`string`) The file path, excluding the file name, relative to the `content` directory.

```go-html-template
{{ with .File }}
  {{ .Dir }}
{{ end }}
```

###### Ext

(`string`) The file extension.

```go-html-template
{{ with .File }}
  {{ .Ext }}
{{ end }}
```

###### Filename

(`string`) The absolute file path.

```go-html-template
{{ with .File }}
  {{ .Filename }}
{{ end }}
```

###### IsContentAdapter

{{< new-in 0.126.0 >}}

(`bool`) Reports whether the file is a [content adapter].

[content adapter]: /content-management/content-adapters/

```go-html-template
{{ with .File }}
  {{ .IsContentAdapter }}
{{ end }}
```

###### LogicalName

(`string`) The file name.

```go-html-template
{{ with .File }}
  {{ .LogicalName }}
{{ end }}
```

###### Path

(`string`) The file path, relative to the `content` directory.

```go-html-template
{{ with .File }}
  {{ .Path }}
{{ end }}
```

###### Section

(`string`) The name of the top level section in which the file resides.

```go-html-template
{{ with .File }}
  {{ .Section }}
{{ end }}
```

###### TranslationBaseName

(`string`) The file name, excluding the extension and language identifier.

```go-html-template
{{ with .File }}
  {{ .TranslationBaseName }}
{{ end }}
```

###### UniqueID

(`string`) The MD5 hash of `.File.Path`.

```go-html-template
{{ with .File }}
  {{ .UniqueID }}
{{ end }}
```

## Examples

Consider this content structure in a multilingual project:

```text
content/
├── news/
│   ├── b/
│   │   ├── index.de.md   <-- leaf bundle
│   │   └── index.en.md   <-- leaf bundle
│   ├── a.de.md           <-- regular content
│   ├── a.en.md           <-- regular content
│   ├── _index.de.md      <-- branch bundle
│   └── _index.en.md      <-- branch bundle
├── _index.de.md
└── _index.en.md
```

With the English language site:

&nbsp;|regular content|leaf bundle|branch bundle
:--|:--|:--|:--
BaseFileName|a.en|index.en|_index.en
ContentBaseName|a|b|news
Dir|news/|news/b/|news/
Ext|md|md|md
Filename|/home/user/...|/home/user/...|/home/user/...
IsContentAdapter|false|false|false
LogicalName|a.en.md|index.en.md|_index.en.md
Path|news/a.en.md|news/b/index.en.md|news/_index.en.md
Section|news|news|news
TranslationBaseName|a|index|_index
UniqueID|15be14b...|186868f...|7d9159d...

## Defensive coding

Some of the pages on a site may not be backed by a file. For example:

- Top level section pages
- Taxonomy pages
- Term pages

Without a backing file, Hugo will throw an error if you attempt to access a `.File` property. To code defensively, first check for file existence:

```go-html-template
{{ with .File }}
  {{ .ContentBaseName }}
{{ end }}
```

[section]: /getting-started/glossary/#section
[taxonomy]: /getting-started/glossary/#taxonomy
[term]: /getting-started/glossary/#term
