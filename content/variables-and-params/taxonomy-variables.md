---
title: Taxonomy Variables
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
tags: [taxonomies,terms]
draft: false
weight: 30
aliases: []
toc: false
needsreview: true
notesforauthors:
---

### Taxonomy Terms Page Variables

[Taxonomy terms](/templates/taxonomy-templates/) pages are of the type `Page` and have the following additional variables. These are available in `layouts/_defaults/terms.html` for example.

* `.Data.Singular` The singular name of the taxonomy<br>
* `.Data.Plural` The plural name of the taxonomy<br>
* `.Data.Pages` the list of pages in this taxonomy<br>
* `.Data.Terms` The taxonomy itself<br>
* `.Data.Terms.Alphabetical` The Terms alphabetized<br>
* `.Data.Terms.ByCount` The Terms ordered by popularity<br>

The last two can also be reversed: `.Data.Terms.Alphabetical.Reverse`, `.Data.Terms.ByCount.Reverse`.

### Taxonomies elsewhere

The `.Site.Taxonomies` variable holds all taxonomies defines site-wide. It is a map of the taxonomy name to a list of its values. For example: "tags" -> ["tag1", "tag2", "tag3"]. Each value, though, is not a string but rather a [Taxonomy variable](#the-taxonomy-variable).

#### The Taxonomy variable

The Taxonomy variable, available as `.Site.Taxonomies.tags` for example, contains the list of tags (values) and, for each of tag, their corresponding content pages.
