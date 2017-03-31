---
title: Taxonomy Templates
linktitle:
description: Taxonomy templating includes taxonomy list pages, taxonomy terms pages, and the ability to iterate through taxonomies in your single page templates as well.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [taxonomies,metadata,front matter,terms]
menu:
  main:
    parent: "Templates"
    weight: 50
weight: 50	#rem
draft: false
aliases: [/taxonomies/displaying/,/templates/terms/,/indexes/displaying/,/taxonomies/templates/,/indexes/ordering/, /templates/taxonomies/, /templates/taxonomy/]
toc: true
---

<!-- NOTE! Check on https://github.com/spf13/hugo/issues/2826 for shifting of terms' pages to .Data.Pages AND
https://discuss.gohugo.io/t/how-to-specify-category-slug/4856/15 for original discussion.-->

Hugo includes support for user-defined groupings of content called **taxonomies**. Taxonomies are classifications that demonstrate logical relationships between content. See [Taxonomies under Content Management](/content-management/taxonomies) if you are unfamiliar with how Hugo leverages this powerful feature.

Hugo provides multiple ways to use taxonomies throughout your project templates:

* Order the way the terms for a taxonomy are displayed in a [taxonomy terms template](#taxonomy-terms-template)
* Order the way content associated with a taxonomy term is displayed in a [taxonomy list template](#taxonomy-list-template)
* List a single content's taxonomy terms within a [single page template][]

## Taxonomy List Templates

Taxonomy list page templates are lists and therefore have all the variables and methods available to [list pages][lists].

### Taxonomy List Template Lookup Order

A Taxonomy will be rendered at /`PLURAL`/`TERM`/ (e.g., http://spf13.com/topics/golang/) from according to the following lookup order:

1. `/layouts/taxonomy/<SINGULAR>.html`
2. `/layouts/_default/taxonomy.html`
3. `/layouts/_default/list.html`
4. `/themes/<THEME>/layouts/taxonomy/<SINGULAR>.html`
5. `/themes/<THEME>/layouts/_default/taxonomy.html`
6. `/themes/<THEME>/layouts/_default/list.html`

## Taxonomy Terms Template

### Taxonomy Terms Templates Lookup Order

{{% warning "The Taxonomy Terms Template has a Unique Lookup Order" %}}
Compared to taxonomy list pages and other list templates such as [sections](/templates/section-templates/), a terms template lookup has only two options. If Hugo does not find a terms template in `layout/` or `/themes/<THEME>/layouts/`, Hugo will *not* render a taxonomy terms page.
{{% /warning %}}

<!-- Begin /taxonomies/methods/ -->
Hugo makes a set of values and methods available on the various Taxonomy structures.

### Taxonomy Methods

A Taxonomy is a `map[string]WeightedPages`.

.Get(term)
: Returns the WeightedPages for a term.

.Count(term)
: The number of pieces of content assigned to this term.

.Alphabetical
: Returns an OrderedTaxonomy (slice) ordered by Term.

.ByCount
: Returns an OrderedTaxonomy (slice) ordered by number of entries.

### OrderedTaxonomy

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
: All Pages assigned to this term. All [list methods][renderlists] are available to this.

## WeightedPages

WeightedPages is simply a slice of WeightedPage.

```go
type WeightedPages []WeightedPage
```

.Count(term)
: The number of pieces of content assigned to this term.

.Pages
: Returns a slice of pages, which then can be ordered using any of the [list methods][renderlists].

<!-- Begin /taxonomies/ordering/ -->

## Ordering Taxonomies

Taxonomies can be ordered by either alphabetical key or by the number of content pieces assigned to that key.

### Order Alphabetically Example

```html
<ul>
  {{ $data := .Data }}
  {{ range $key, $value := .Data.Taxonomy.Alphabetical }}
  <li><a href="{{ .Site.LanguagePrefix }}/{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
  {{ end }}
</ul>
```

### Order by Popularity Example

```html
<ul>
  {{ $data := .Data }}
  {{ range $key, $value := .Data.Taxonomy.ByCount }}
  <li><a href="{{ .Site.LanguagePrefix }}/{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
  {{ end }}
</ul>
```

<!-- [See Also Taxonomy Lists](/templates/list/) -->

## Ordering Content within Taxonomies

Hugo uses both `date` and `weight` to order content within taxonomies.

Each piece of content in Hugo can optionally be assigned a date. It can also be assigned a weight for each taxonomy it is assigned to.

When iterating over content within taxonomies, the default sort is the same as that used for [section and list pages]() first by weight then by date. This means that if the weights for two pieces of content are the same, than the more recent content will be displayed first. The default weight for any piece of content is 0.

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

A [list template](/templates/list/) is any template that will be used to render multiple pieces of content in a single html page. This template will be used to generate all the automatically created taxonomy pages.

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

### Displaying a Single Piece of Content's Taxonomies

Within your content templates, you may wish to display the taxonomies that piece of content is assigned to.

Because we are leveraging the front matter system to define taxonomies for content, the taxonomies assigned to each content piece are located in the usual place (i.e., `.Params.<TAXONOMYPLURAL>`).

### Example: List Tags in a Single Page Template

```html
<ul id="tags">
  {{ range .Params.tags }}
    <li><a href="{{ "/tags/" | relLangURL }}{{ . | urlize }}">{{ . }}</a> </li>
  {{ end }}
</ul>
```

If you want to list taxonomies inline, you will have to take care of optional plural endings in the title (if multiple taxonomies), as well as commas. Let's say we have a taxonomy "directors" such as `directors: [ "Joel Coen", "Ethan Coen" ]` in the TOML-format front matter.

To list such taxonomies, use the following:

### Example: Comma-delimit Tags in a Single Page Template

```html
{{ if .Params.directors }}
  <strong>Director{{ if gt (len .Params.directors) 1 }}s{{ end }}:</strong>
  {{ range $index, $director := .Params.directors }}{{ if gt $index 0 }}, {{ end }}<a href="{{ "directors/" | relURL }}{{ . | urlize }}">{{ . }}</a>{{ end }}
{{ end }}
```

Alternatively, you may use the [delimit template function][delimit] as a shortcut if the taxonomies should just be listed with a separator. See {{< gh 2143 >}} on GitHub for discussion.

## Listing Content with the Same Taxonomy Term

If you are using a taxonomy for something like a series of posts, you can list individual pages associated with the same taxonomy. This is also a quick and dirty method for showing related content:

### Example: Showing Content in Same Series

```html
<ul>
  {{ range .Site.Taxonomies.series.golang }}
    <li><a href="{{ .Page.RelPermalink }}">{{ .Page.Title }}</a></li>
  {{ end }}
</ul>
```

## Listing All content in a Given taxonomy

This would be very useful in a sidebar as “featured content”. You could even have different sections of “featured content” by assigning different terms to the content.

### Example: Grouping "Featured" Content

```html
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
```

## Rendering a Site's Taxonomies

If you wish to display the list of all keys for your site's taxonomy, you can retrieve them from the [`.Site` variable][sitevars] available on every page.

This may take the form of a tag cloud, a menu, or simply a list.

The following example displays all terms in a site's tags taxonomy:

### Example: List All Site Tags

```html
<ul id="all-tags">
  {{ range $name, $taxonomy := .Site.Taxonomies.tags }}
    <li><a href="{{ "/tags/" | relLangURL }}{{ $name | urlize }}">{{ $name }}</a></li>
  {{ end }}
</ul>
```

### Example: List All Taxonomies, Terms, and Assigned Content

This example will list all taxonomies and their terms, as well as all the content assigned to each of the terms.

{{% code file="layouts/partials/all-taxonomies.html" download="all-taxonomies.html" download="all-taxonomies.html" %}}
```html
<section>
  <ul id="all-taxonomies">
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
```
{{% /code %}}

<!-- ## `.Site.GetPage` for Taxonomies

### `.Site.GetPage` Taxonomy List Example

### `.Site.GetPage` Taxonomy Terms Example -->


[delimit]: /functions/delimit/
[lists]: /templates/lists/
[renderlists]: /templates/lists/
[single page template]: /templates/single-page-templates/
[sitevars]: /variables/site/