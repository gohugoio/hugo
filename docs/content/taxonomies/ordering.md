---
title: "Ordering Taxonomies"
date: "2013-07-01"
linktitle: "Ordering"
aliases: ["/indexes/ordering/"]
weight: 60
menu:
  main:
    identifier: "Ordering Taxonomies"
    parent: 'taxonomy'
---

Hugo provides the ability to both:

 1. Order the way the keys for an index are displayed
 2. Order the way indexed content appears


## Ordering Indexes
Indexes can be ordered by either alphabetical key or by the number of content pieces assigned to that key.

### Order Alphabetically Example:

    <ul>
    {{ $data := .Data }}
    {{ range $key, $value := .Data.Index.Alphabetical }}
    <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
    {{ end }}
    </ul>

### Order by Popularity Example:

    <ul>
    {{ $data := .Data }}
    {{ range $key, $value := .Data.Index.ByCount }}
    <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
    {{ end }}
    </ul>


[See Also Index Lists](/indexes/lists/)

## Ordering Content within Indexes

Hugo uses both **Date** and **Weight** to order content within indexes.

Each piece of content in Hugo can optionally be assigned a date.
It can also be assigned a weight for each index it is assigned to.

When iterating over content within indexes the default sort is first by weight then by date. This means that if the weights for two pieces of content are the same, than the more recent content will be displayed first. The default weight for any piece of content is 0.

### Assigning Weight

Content can be assigned weight for each index that it's assigned to.

    +++
    tags = [ "a", "b", "c" ]
    tags_weight = 22
    categories = ["d"]
    title = "foo"
    categories_weight = 44
    +++
    Front Matter with weighted tags and categories


The convention is `indexname_weight`.

In the above example, this piece of content has a weight of 22 which applies to the sorting when rendering the pages assigned to the "a", "b" and "c" values of the 'tag' index.

It has also been assigned the weight of 44 when rendering the 'd' category.

With this the same piece of content can appear in different positions in different indexes.

Currently indexes only support the default ordering of content which is weight -> date.
