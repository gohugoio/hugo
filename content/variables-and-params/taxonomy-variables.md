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
toc: true
needsreview: true
notesforauthors:
---

### Taxonomy Terms Page Variables

[Taxonomy terms pages][taxonomytemplates] are of the type `Page` and have the following additional variables. For example, the following fields would be available in `layouts/_defaults/terms.html`, depending on how you organize your [taxonomy templates][taxonomytemplates]:

`.Data.Singular`
: The singular name of the taxonomy

`.Data.Plural`
: The plural name of the taxonomy

`.Data.Pages`
: The list of pages in the taxonomy

`.Data.Terms`
: The taxonomy itself

`.Data.Terms.Alphabetical`
: The taxonomy terms alphabetized

`.Data.Terms.ByCount`
: The Terms ordered by popularity

Note that `.Data.Terms.Alphabetical` and `.Data.Terms.ByCount` can also be reversed:

* `.Data.Terms.Alphabetical.Reverse`
* `.Data.Terms.ByCount.Reverse`

### Using `.Site.Taxonomies` Outside of Taxonomy Templates

The `.Site.Taxonomies` variable holds all the taxonomies that defined site-wide. It is a map of the taxonomy name to a list of its values (e.g., "tags" -> ["tag1", "tag2", "tag3"]). Each value, though, is not a string but rather a **Taxonomy variable**.

### The Taxonomy Variable

The Taxonomy variable, available, for example, as `.Site.Taxonomies.tags`, contains the list of tags (values) and, for each of tag, their corresponding content pages.

### Example Usage of `.Site.Taxonomies`

**NEEDS EXAMPLE**

[taxonomytemplates]: /templates/taxonomy-templates/