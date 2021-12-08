---
title: readFile
description: Returns the contents of a file.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2021-11-26
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [files]
signature: ["os.ReadFile PATH", "readFile PATH"]
workson: []
hugoversion:
relatedfuncs: ['os.FileExists','os.ReadDir','os.Stat']
deprecated: false
aliases: []
---
The `os.ReadFile` function attempts to resolve the path relative to the root of your project directory. If a matching file is not found, it will attempt to resolve the path relative to the [`contentDir`]({{< relref "getting-started/configuration#contentdir">}}). A leading path separator (`/`) is optional.

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

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates]({{< relref "/templates/files" >}}).
