---
title: last
linktitle: last
description: "slices an array to only the last <em>N</em>th elements."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: []
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: ["last INDEX COLLECTION"]
workson: [lists, taxonomies, terms, groups]
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---


```
{{ range last 10 .Pages }}
    {{ .Render "summary" }}
{{ end }}
```
