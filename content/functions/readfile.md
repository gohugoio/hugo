---
title: readFile
description: Reads a file from disk relative to the current project working directory and returns a string.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-04-30
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [files]
signature: ["readFile PATH"]
workson: []
hugoversion:
relatedfuncs: [readDir]
deprecated: false
aliases: []
---

Note that the filename must be relative to the current project working directory.

So, if you have a file with the name `README.txt` in the root of your project with the content `Hugo Rocks!`:

```
{{readFile "README.txt"}} → "Hugo Rocks!"
```

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates][local].

[local]: /templates/files/
