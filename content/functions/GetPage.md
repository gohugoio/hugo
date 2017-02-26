---
title: getpage
linktitle: GetPage
description:
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

Every `Page` has a `Kind` attribute that shows what kind of page it is. While this attribute can be used to list pages of a certain `kind` using `where`, often it can be useful to fetch a single page by its path.

`GetPage` looks up an index page of a given `Kind` and `path`. This method may support regular pages in the future, but currently it is a convenient way of getting the index pages, such as the home page or a section, from a template:

    {{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}

This method wil return `nil` when no page could be found, so the above will not print anything if the blog section isn't found.

The valid page kinds are: *home, section, taxonomy and taxonomyTerm.*

**The following is a good example of `.GetPage`**:

<https://discuss.gohugo.io/t/problem-with-loop-and-scratch/5597

