---
title: readDir
description: Gets a directory listing from a directory relative to the current working directory.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [files]
signature: ["readDir PATH"]
workson: []
hugoversion:
relatedfuncs: [readFile]
deprecated: false
aliases: []
---

If your current project working directory has a single file named `README.txt`:

```
{{ range (readDir ".") }}{{ .Name }}{{ end }} → "README.txt"
```

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates][local].

[local]: /templates/files/

