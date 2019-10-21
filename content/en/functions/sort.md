---
title: sort
# linktitle: sort
description: Sorts maps, arrays, and slices and returns a sorted slice.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [ordering,sorting,lists]
signature: []
workson: [lists,taxonomies,terms,groups]
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

A sorted array of map values will be returned with the keys eliminated. There are two optional arguments: `sortByField` and `sortAsc`. If left blank, sort will sort by keys (for maps) in ascending order as its default behavior.

```
+++
keywords: [ "tag3", "tag1", "tag2" ]
+++

// Site config
+++
[params.authors]
  [params.authors.Derek]
    "firstName"  = "Derek"
    "lastName"   = "Perkins"
  [params.authors.Joe]
    "firstName"  = "Joe"
    "lastName"   = "Bergevin"
  [params.authors.Tanner]
    "firstName"  = "Tanner"
    "lastName"   = "Linsley"
+++
```

```
// Use default sort options - sort by key / ascending
Tags: {{ range sort .Params.tags }}{{ . }} {{ end }}

→ Outputs Tags: tag1 tag2 tag3

// Sort by value / descending
Tags: {{ range sort .Params.tags "value" "desc" }}{{ . }} {{ end }}

→ Outputs Tags: tag3 tag2 tag1

// Use default sort options - sort by value / descending
Authors: {{ range sort .Site.Params.authors }}{{ .firstName }} {{ end }}

→ Outputs Authors: Derek Joe Tanner

// Use default sort options - sort by value / descending
Authors: {{ range sort .Site.Params.authors "lastName" "desc" }}{{ .lastName }} {{ end }}

→ Outputs Authors: Perkins Linsley Bergevin
```

