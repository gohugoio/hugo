---
title: readdir and readfile
linktitle:
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [files]
categories: [functions]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: [/functions/readdir/,/functions/readfile/]
---

## `readDir`

`readDir` gets a directory listing from a directory relative to the current project working directory.

If your current project working directory has a single file named `README.txt`:

```
{{ range (readDir ".") }}{{ .Name }}{{ end }} → "README.txt"
```

## `readFile`

Reads a file from disk and converts it into a string. Note that the filename must be relative to the current project working directory.

So, if you have a file with the name `README.txt` in the root of your project with the content `Hugo Rocks!`:

```
{{readFile "README.txt"}} → "Hugo Rocks!"
```

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates][local].


[local]: /templates/local-file-templates/

