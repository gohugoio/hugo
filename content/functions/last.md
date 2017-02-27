---
title: last
linktitle: last
description: Slices an array to only the last Nth elements.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: []
categories: [functions]
toc:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`last` slices an array to only the last _N_th elements.

Works on [lists](/templates/list/), [taxonomies](/taxonomies/displaying/), [terms](/templates/terms/), [groups](/templates/list/)

```
{{ range last 10 .Data.Pages }}
    {{ .Render "summary" }}
{{ end }}
```

