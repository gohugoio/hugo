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

Note that the filename must be relative to the current project working directory, or the project's `/content` folder. 

So, if you have a file with the name `README.txt` in the root of your project with the content `Hugo Rocks!`:

```
{{readFile "README.txt"}} → "Hugo Rocks!"
```

If you receive a "file doesn't exist" error with a path listed, do take note that the path is the last one checked by the function, and may not accurately reflect your target. You should generally double-check your path for mistakes.

Note that there is a 1 MB file size limit.

For more information on using `readDir` and `readFile` in your templates, see [Local File Templates][local].

[local]: /templates/files/
