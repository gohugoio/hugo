---
title: Taxonomy Variables
linktitle:
description: Taxonomy pages are of type `Page` and have all page-, site-, and list-level variables available to them. However, taxonomy terms templates have additional variables available to their templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
keywords: [taxonomies,terms]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 30
weight: 30
sections_weight: 30
aliases: []
toc: true
---

## Taxonomy Terms Page Variables

[Taxonomy terms pages][taxonomytemplates] are of the type `Page` and have the following additional variables.

For example, the following fields would be available in `layouts/_defaults/terms.html`, depending on how you organize your [taxonomy templates][taxonomytemplates]:

.Data.Singular
: The singular name of the taxonomy (e.g., `tags => tag`)

.Data.Plural
: The plural name of the taxonomy (e.g., `tags => tags`)

.Data.Pages
: The list of pages in the taxonomy

.Data.Terms
: The taxonomy itself

.Data.Terms.Alphabetical
: The taxonomy terms alphabetized

.Data.Terms.ByCount
: The Terms ordered by popularity

Note that `.Data.Terms.Alphabetical` and `.Data.Terms.ByCount` can also be reversed:

* `.Data.Terms.Alphabetical.Reverse`
* `.Data.Terms.ByCount.Reverse`

## Use `.Site.Taxonomies` Outside of Taxonomy Templates

The `.Site.Taxonomies` variable holds all the taxonomies defined site-wide. `.Site.Taxonomies` is a map of the taxonomy name to a list of its values (e.g., `"tags" -> ["tag1", "tag2", "tag3"]`). Each value, though, is not a string but rather a *Taxonomy variable*.

## The `.Taxonomy` Variable

The `.Taxonomy` variable, available, for example, as `.Site.Taxonomies.tags`, contains the list of tags (values) and, for each tag, their corresponding content pages.

### Example Usage of `.Site.Taxonomies`

The following [partial template][partials] will list all your site's taxonomies, each of their keys, and all the content assigned to each of the keys. For more examples of how to order and render your taxonomies, see  [Taxonomy Templates][taxonomytemplates].

{{< code file="all-taxonomies-keys-and-pages.html" download="all-taxonomies-keys-and-pages.html" >}}
<section>
  <ul>
    {{ range $taxonomyname, $taxonomy := .Site.Taxonomies }}
      <li><a href="{{ "/" | relLangURL}}{{ $taxonomyname | urlize }}">{{ $taxonomyname }}</a>
        <ul>
          {{ range $key, $value := $taxonomy }}
          <li> {{ $key }} </li>
                <ul>
                {{ range $value.Pages }}
                    <li><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                {{ end }}
                </ul>
          {{ end }}
        </ul>
      </li>
    {{ end }}
  </ul>
</section>
{{< /code >}}

[partials]: /templates/partials/
[taxonomytemplates]: /templates/taxonomy-templates/
