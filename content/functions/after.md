---
title: after
linktitle: after
description: Slices an array to only the items after the Nth item.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [iteration]
ns:
signature:
workson: []
hugoversion:
relatedfuncs: [last]
deprecated: false
aliases: []
---

`after` slices an array to only the items after the *N*th item. Combining `after` with `first` uses both use both halves of an array split at item *N*.

Works on [lists](/templates/list/), [taxonomies](/taxonomies/displaying/), [terms](/templates/terms/), [groups](/templates/list/)

e.g.

    {{ range after 10 .Data.Pages }}
        {{ .Render "title" }}
    {{ end }}
