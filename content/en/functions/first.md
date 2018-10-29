---
title: first
linktitle: first
description: "Slices an array to only the first _N_ elements."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [iteration]
signature: ["first LIMIT COLLECTION"]
workson: [lists,taxonomies,terms,groups]
hugoversion:
relatedfuncs: [after,last]
deprecated: false
aliases: []
---


```
{{ range first 10 .Pages }}
    {{ .Render "summary" }}
{{ end }}
```

*Note: Exclusive to `first`, LIMIT can be '0' to return an empty array.*
