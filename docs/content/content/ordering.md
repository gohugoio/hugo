---
date: 2014-03-06
linktitle: Ordering
menu:
  main:
    parent: content
next: /content/summaries
prev: /content/archetypes
title: Ordering Content
weight: 60
---

Hugo provides you with all the flexibility you need to organize how your content is ordered.

By default, content is ordered by weight, then by date with the most
recent date first, but alternative sorting (by `title` and `linktitle`) is
also available. The order the content would appear is specified in
the [list template](/templates/list/).

_Both the `date` and `weight` fields are optional._

Unweighted pages appear at the end of the list. If no weights are provided (or
if weights are the same), `date` will be used to sort. If neither is provided,
content will be ordered based on how it's read off the disk, and no order is
guaranteed.

## Assigning weight to content

    +++
    weight = 4
    title = "Three"
    date = "2012-04-06"
    +++
    Front Matter with Ordered Pages 3


## Ordering Content Within Taxonomies

Please see the [Taxonomy Ordering Documentation](/taxonomies/ordering/).
