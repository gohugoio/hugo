---
title: isset
linktitle: isset
description: Returns true if the parameter is set.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`isset` returns true if the parameter is set.
Takes either a slice, array or channel and an index or a map and a key as input.

```
{{ if isset .Params "project_url" }} {{ index .Params "project_url" }}{{ end }}
```

