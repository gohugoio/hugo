---
aliases:
- /indexes/ordering/
lastmod: 2015-12-23
date: 2013-07-01
linktitle: Ordering
menu:
  main:
    identifier: Ordering Taxonomies
    parent: taxonomy
next: /taxonomies/methods
prev: /taxonomies/templates
title: Ordering Taxonomies
weight: 60
toc: true
---

Hugo provides the ability to both:

 1. Order the way the keys for a taxonomy are displayed
 2. Order the way taxonomyed content appears


## Ordering Taxonomies
Taxonomies can be ordered by either alphabetical key or by the number of content pieces assigned to that key.

### Order Alphabetically Example

    <ul>
    {{ $data := .Data }}
    {{ range $key, $value := .Data.Taxonomy.Alphabetical }}
    <li><a href="{{ .Site.LanguagePrefix }}/{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
    {{ end }}
    </ul>

### Order by Popularity Example

    <ul>
    {{ $data := .Data }}
    {{ range $key, $value := .Data.Taxonomy.ByCount }}
    <li><a href="{{ .Site.LanguagePrefix }}/{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
    {{ end }}
    </ul>


[See Also Taxonomy Lists]({{< relref "templates/list.md" >}})

## Ordering Content within Taxonomies

Hugo uses both **Date** and **Weight** to order content within taxonomies.

Each piece of content in Hugo can optionally be assigned a date.
It can also be assigned a weight for each taxonomy it is assigned to.

When iterating over content within taxonomies the default sort is first by weight then by date. This means that if the weights for two pieces of content are the same, than the more recent content will be displayed first. The default weight for any piece of content is 0.

### Assigning Weight

Content can be assigned weight for each taxonomy that it's assigned to.

```toml
+++
tags = [ "a", "b", "c" ]
tags_weight = 22
categories = ["d"]
title = "foo"
categories_weight = 44
+++
Front Matter with weighted tags and categories
```

The convention is `taxonomyname_weight`.

In the above example, this piece of content has a weight of 22 which applies to the sorting when rendering the pages assigned to the "a", "b" and "c" values of the 'tag' taxonomy.

It has also been assigned the weight of 44 when rendering the 'd' category.

With this the same piece of content can appear in different positions in different taxonomies.

Currently taxonomies only support the default ordering of content which is weight -> date.
