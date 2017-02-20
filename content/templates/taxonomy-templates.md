---
title: Taxonomy Templates
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [taxonomies,metadata,front matter,terms]
categories: [templates]
weight:
draft: false
slug:
aliases: [/taxonomies/displaying/,/indexes/displaying/,/taxonomies/templates/,/indexes/ordering/]
needsreview: true
toc: true
notesforauthors:
---

<!-- Begin /taxonomies/methods/ -->

Hugo makes a set of values and methods available on the various Taxonomy structures.

## Taxonomy Methods

A Taxonomy is a `map[string]WeightedPages`.

.Get(term)
: Returns the WeightedPages for a term.

.Count(term)
: The number of pieces of content assigned to this term.

.Alphabetical
: Returns an OrderedTaxonomy (slice) ordered by Term.

.ByCount
: Returns an OrderedTaxonomy (slice) ordered by number of entries.

## OrderedTaxonomy

Since Maps are unordered, an OrderedTaxonomy is a special structure that has a defined order.

```go
[]struct {
    Name          string
    WeightedPages WeightedPages
}
```

Each element of the slice has:

.Term
: The Term used.

.WeightedPages
: A slice of Weighted Pages.

.Count
: The number of pieces of content assigned to this term.

.Pages
: All Pages assigned to this term. All [list methods](/templates/list/) are available to this.

## WeightedPages

WeightedPages is simply a slice of WeightedPage.

```go
type WeightedPages []WeightedPage
```

.Count(term)
: The number of pieces of content assigned to this term.

.Pages
: Returns a slice of pages, which then can be ordered using any of the [list methods](/templates/list/).

<!-- Begin /taxonomies/ordering/ -->

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


[See Also Taxonomy Lists](/templates/list/)

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

<!-- Begin /taxonomies/templates/ -->

There are two different templates that the use of taxonomies will require you to provide.

Both templates are covered in detail in the templates section.

A [list template](/templates/list/) is any template that will be used to render multiple pieces of
content in a single html page. This template will be used to generate
all the automatically created taxonomy pages.

A [taxonomy terms template](/templates/terms/) is a template used to
generate the list of terms for a given template.

<!-- Begin /taxonomies/displaying/ -->

There are four common ways you can display the data in your
taxonomies in addition to the automatic taxonomy pages created by hugo
using the [list templates](/templates/list/):

1. For a given piece of content, you can list the terms attached
2. For a given piece of content, you can list other content with the same
   term
3. You can list all terms for a taxonomy
4. You can list all taxonomies (with their terms)

## 1. Displaying taxonomy terms assigned to this content

Within your content templates, you may wish to display
the taxonomies that that piece of content is assigned to.

Because we are leveraging the front matter system to
define taxonomies for content, the taxonomies assigned to
each content piece are located in the usual place
(.Params.`plural`).

### Example

    <ul id="tags">
      {{ range .Params.tags }}
        <li><a href="{{ "/tags/" | relLangURL }}{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>

If you want to list taxonomies inline, you will have to take
care of optional plural ending in the title (if multiple taxonomies),
as well as commas. Let's say we have a taxonomy "directors" such as
`directors: [ "Joel Coen", "Ethan Coen" ]` in the TOML-format front matter.
To list such taxonomy use the following:

### Example

    {{ if .Params.directors }}
      <strong>Director{{ if gt (len .Params.directors) 1 }}s{{ end }}:</strong>
      {{ range $index, $director := .Params.directors }}{{ if gt $index 0 }}, {{ end }}<a href="{{ "/directors/" | relURL }}{{ . | urlize }}">{{ . }}</a>{{ end }}
    {{ end }}

Alternatively, you may use the [delimit](/templates/functions/) template function as a shortcut if the taxonomies should just be listed with a separator.  See {{< gh 2143 >}} on GitHub for discussion.

## 2. Listing content with the same taxonomy term

First, you may be asking why you would use this. If you are using a taxonomy for something like a series of posts, this is exactly how you would do it. It’s also an quick and dirty way to show some related content.

### Example

```html
<ul>
  {{ range .Site.Taxonomies.series.golang }}
    <li><a href="{{ .Page.RelPermalink }}">{{ .Page.Title }}</a></li>
  {{ end }}
</ul>
```

## 3. Listing all content in a given taxonomy

This would be very useful in a sidebar as “featured content”. You could
even have different sections of “featured content” by assigning
different terms to the content.

### Example

    <section id="menu">
        <ul>
            {{ range $key, $taxonomy := .Site.Taxonomies.featured }}
            <li> {{ $key }} </li>
            <ul>
                {{ range $taxonomy.Pages }}
                <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                {{ end }}
            </ul>
            {{ end }}
        </ul>
    </section>


## 4. Rendering a Site's Taxonomies

If you wish to display the list of all keys for a taxonomy, you can find retrieve
them from the `.Site` variable which is available on every page.

This may take the form of a tag cloud, a menu or simply a list.

The following example displays all tag keys:

### Example

    <ul id="all-tags">
      {{ range $name, $taxonomy := .Site.Taxonomies.tags }}
        <li><a href="{{ "/tags/" | relLangURL }}{{ $name | urlize }}">{{ $name }}</a></li>
      {{ end }}
    </ul>

### Complete Example
This example will list all taxonomies, each of their keys and all the content assigned to each key.

    <section>
      <ul>
        {{ range $taxonomyname, $taxonomy := .Site.Taxonomies }}
          <li><a href="{{ "/" | relLangURL}}{{ $taxonomyname | urlize }}">{{ $taxonomyname }}</a>
            <ul>
              {{ range $key, $value := $taxonomy }}
              <li> {{ $key }} </li>
                    <ul>
                    {{ range $value.Pages }}
                        <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                    {{ end }}
                    </ul>
              {{ end }}
            </ul>
          </li>
        {{ end }}
      </ul>
    </section>
