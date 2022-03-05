---
title: readDir
description: Returns an array of FileInfo structures sorted by filename, one element for each directory entry.
publishdate: 2017-02-01
lastmod: 2021-11-26
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [files]
signature: ["os.ReadDir PATH", "readDir PATH", "readDir PATH DIRECTORY"]
workson: []
hugoversion:
relatedfuncs: ['os.FileExists','os.ReadFile','os.Stat']
deprecated: false
aliases: []
---
The `os.ReadDir` function resolves the path relative to the root of your project directory unless the optional directory parameter is provided. A leading path separator (`/`) is optional.

The following DIRECTORY values are supported
* assets
* archetypes
* data
* content
* layouts
* i18n

For more information on the hugo directory structure see [Directory Structure]({{< relref "/getting-started/directory-structure" >}}).

{{% note %}}
Note that `os.ReadDir` is not recursive.
{{% /note %}}

Details of the `FileInfo` structure are available in the [Go documentation](https://pkg.go.dev/io/fs#FileInfo).

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates]({{< relref "/templates/files" >}}).

## Example with path parameter
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

## Example with path and directory parameters

With this directory structure:

```text
.
├── config.toml
├── themes/
    ├── mytheme/
        ├── mytheme.txt
        └── archetypes/
          └── mytheme-default.md
└── archetypes/
    ├── default.md
    └── post-bundle/
        ├── bio.md
        ├── index.md
        └── images/
            └── featured.jpg
```

This template code:

```go-html-template
{{ range os.ReadDir "." "archetypes" }}
  {{ .Name }} --> {{ .IsDir }}
{{ end }}
```

Produces:

```html
default.md --> false
post-bundle --> true
mytheme-default.md --> false
```
