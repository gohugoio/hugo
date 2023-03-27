---
title: Taxonomy Variables
linktitle:
description: Hugo's taxonomy system exposes variables to taxonomy and term templates.
categories: [variables and params]
keywords: [taxonomy,term]
menu:
  docs:
    parent: "variables"
    weight: 30
toc: true
weight: 30
aliases: []
---

## Taxonomy templates

Pages rendered by taxonomy templates have `.Kind` set to `taxonomy` and `.Type` set to the taxonomy name.

In taxonomy templates you may access `.Site`, `.Page`. `.Section`, and `.File` variables, as well as the following _taxonomy_ variables:

.Data.Singular
: The singular name of the taxonomy (e.g., `tags => tag`).

.Data.Plural
: The plural name of the taxonomy (e.g., `tags => tags`).

.Data.Pages
: The collection of term pages related to this taxonomy. Aliased by `.Pages`.

.Data.Terms
: A map of terms and weighted pages related to this taxonomy.

.Data.Terms.Alphabetical
: A map of terms and weighted pages related to this taxonomy, sorted alphabetically in ascending order. Reverse the sort order with`.Data.Terms.Alphabetical.Reverse`.

.Data.Terms.ByCount
: A map of terms and weighted pages related to this taxonomy, sorted by count in ascending order. Reverse the sort order with`.Data.Terms.ByCount.Reverse`.

## Term templates

Pages rendered by term templates have `.Kind` set to `term` and `.Type` set to the taxonomy name.

In term templates you may access `.Site`, `.Page`. `.Section`, and `.File` variables, as well as the following _term_ variables:

.Data.Singular
: The singular name of the taxonomy (e.g., `tags => tag`).

.Data.Plural
: The plural name of the taxonomy (e.g., `tags => tags`).

.Data.Pages
: The collection of content pages related to this taxonomy. Aliased by `.Pages`.

.Data.Term
: The term itself (e.g., `tag-one`).

## Access taxonomy data from any template

Access the entire taxonomy data structure from any template with `site.Taxonomies`. This returns a map of taxonomies, terms, and a collection of weighted content pages related to each term. For example:

```json
{
  "categories": {
    "news": [
      {
        "Weight": 0,
        "Page": {
          "Title": "Post 1",
          "Date": "2022-12-18T15:13:35-08:00"
          ...
          }
      },
      {
        "Weight": 0,
        "Page": {
          "Title": "Post 2",
          "Date": "2022-12-18T15:13:46-08:00",
          ...
        }
      }
    ]
  },
  "tags": {
    "international": [
      {
        "Weight": 0,
        "Page": {
          "Title": "Post 1",
          "Date": "2021-01-01T00:00:00Z"
          ... 
        }
      }
    ]
  }
}
```

Access a subset of the taxonomy data structure by chaining one or more identifiers, or by using the [`index`] function with one or more keys. For example, to access the collection of weighted content pages related to the news category, use either of the following:

[`index`]: /functions/index-function/

```go-html-template
{{ $pages := site.Taxonomies.categories.news }}
{{ $pages := index site.Taxonomies "categories" "news" }}
```

For example, to render the entire taxonomy data structure as a nested unordered list:

```go-html-template
<ul>
  {{ range $taxonomy, $terms := site.Taxonomies }}
    <li>
      {{ with site.GetPage $taxonomy }}
        <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
      {{ end }}
      <ul>
        {{ range $term, $weightedPages := $terms }}
        <li>
          {{ with site.GetPage (path.Join $taxonomy $term) }}
            <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
          {{ end }}
        </li>
          <ul>
            {{ range $weightedPages }}
              <li>
                <a href="{{ .RelPermalink }}"> {{ .LinkTitle }}</a>
              </li>
            {{ end }}
          </ul>
        {{ end }}
      </ul>
    </li>
  {{ end }}
</ul>
```

See [Taxonomy Templates] for more examples.

[Taxonomy Templates]: /templates/taxonomy-templates/
