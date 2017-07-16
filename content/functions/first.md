---
title: first
linktitle: first
description: Slices an array to only the first Nth elements.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
#tags: [iteration]
signature: ["first LIMIT COLLECTION"]
workson: [lists,taxonomies,terms,groups]
hugoversion:
relatedfuncs: [after,last]
deprecated: false
aliases: []
---

`first` slices an array to only the first _N_th elements.

```golang
{{ range first 10 .Data.Pages }}
    {{ .Render "summary" }}
{{ end }}
```

