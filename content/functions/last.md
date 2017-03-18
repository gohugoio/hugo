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
ns:
signature:
workson: [lists, taxonomies, terms, groups]
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`last` slices an array to only the last <em>N</em>th elements.

```
{{ range last 10 .Data.Pages }}
    {{ .Render "summary" }}
{{ end }}
```

