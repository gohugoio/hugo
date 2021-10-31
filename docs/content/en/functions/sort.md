---
title: sort
# linktitle: sort
description: Sorts maps, arrays, and slices and returns a sorted slice.
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
---
tags: ["tag3", "tag1", "tag2"]
---

// Site config
+++
[params.authors]
  [params.authors.Joe]
    firstName = "Joe"
    lastName  = "Bergevin"
  [params.authors.Derek]
    firstName = "Derek"
    lastName  = "Perkins"
  [params.authors.Tanner]
    firstName = "Tanner"
    lastName  = "Linsley"
+++
```

```
// Sort by value, ascending (default for lists)
Tags: {{ range sort .Params.tags }}{{ . }} {{ end }}

→ Outputs Tags: tag1 tag2 tag3

// Sort by value, descending
Tags: {{ range sort .Params.tags "value" "desc" }}{{ . }} {{ end }}

→ Outputs Tags: tag3 tag2 tag1

// Sort by key, ascending (default for maps)
Authors: {{ range sort .Site.Params.authors }}{{ .firstName }} {{ end }}

→ Outputs Authors: Derek Joe Tanner

// Sort by field, descending
Authors: {{ range sort .Site.Params.authors "lastName" "desc" }}{{ .lastName }} {{ end }}

→ Outputs Authors: Perkins Linsley Bergevin
```
