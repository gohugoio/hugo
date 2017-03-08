---
title: readfile
linktitle: readFile
description: Reads a file from disk relative to the current project working directory and converts it into a string.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [files]
signature:
workson: []
hugoversion:
relatedfuncs: [readDir]
deprecated: false
aliases: []
---

`readFile` reads a file from disk and converts it into a string. Note that the filename must be relative to the current project working directory.

So, if you have a file with the name `README.txt` in the root of your project with the content `Hugo Rocks!`:

```
{{readFile "README.txt"}} → "Hugo Rocks!"
```

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates][local].

[local]: /templates/files/