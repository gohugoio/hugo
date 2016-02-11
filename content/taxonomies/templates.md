---
aliases:
- /indexes/templates/
lastmod: 2014-05-29
date: 2013-07-01
linktitle: Templates
menu:
  main:
    parent: taxonomy
next: /taxonomies/ordering
prev: /templates/displaying
title: Taxonomy Templates
weight: 30
---

Taxonomy templates should be placed in the folder `layouts/taxonomy`.
When Taxonomy term template is provided for a Taxonomy, a section is rendered for it at `/SINGULAR/`. (eg. `/tag/` or `/category/`)

There are two different templates that the use of taxonomies will require you to provide:

### All content attached to taxonomy

A [taxonomy terms template](/templates/terms/) is a template which has access to all the full Taxonomy structure.
This Template is commonly used to generate the list of terms for a given template.


#### layouts/taxonomy/SINGULAR.terms.html

For example: `tag.terms.html`, `category.terms.html`, or your custom Taxonomy: `actor.terms.html`

### All content attached to term

A [list template](/templates/list/) is used to automatically generate pages for each unique term found.


#### layouts/taxonomy/SINGULAR.html

For example: `tag.html`, `category.html`, or your custom Taxonomy: `actor.html`

Terms are rendered at `SINGULAR/TERM/`. (eg. `/tag/book/` or `/category/news/`)
