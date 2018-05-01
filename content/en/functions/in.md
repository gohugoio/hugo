---
title: in
linktitle:
description: Checks if an element is in an array or slice--or a substring in a string---and returns a boolean.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["in SET ITEM"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

The elements supported are strings, integers and floats, although only float64 will match as expected.

In addition, `in` can also check if a substring exists in a string.

```
{{ if in .Params.tags "Git" }}Follow me on GitHub!{{ end }}
```


```
{{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}
```
