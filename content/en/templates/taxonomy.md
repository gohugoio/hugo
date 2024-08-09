---
title: Taxonomy templates
description: Create a taxonomy template to render a list of terms.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 90
weight: 90
toc: true
aliases: [/taxonomies/displaying/,/templates/terms/,/indexes/displaying/,/taxonomies/templates/,/indexes/ordering/, /templates/taxonomies/, /templates/taxonomy-templates/]
---

The [taxonomy] template below inherits the site's shell from the [base template], and renders a list of [terms] in the current taxonomy.

[taxonomy]: /getting-started/glossary/#taxonomy
[terms]: /getting-started/glossary/#term
[base template]: /templates/types/

{{< code file=layouts/_default/taxonomy.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
{{< /code >}}

Review the [template lookup order] to select a template path that provides the desired level of specificity.

[template lookup order]: /templates/lookup-order/#taxonomy-templates

In the example above, the taxonomy and term will be capitalized if their respective pages are not backed by files. You can disable this in your site configuration:

{{< code-toggle file=hugo >}}
capitalizeListTitles = false
{{< /code-toggle >}}

## Data object

Use these methods on the `Data` object within a taxonomy template.

Singular
: (`string`) Returns the singular name of the taxonomy.

```go-html-template
{{ .Data.Singular }} → tag
```

Plural
: (`string`) Returns the plural name of the taxonomy.

```go-html-template
{{ .Data.Plural }} → tags
```

Terms
: (`page.Taxonomy`) Returns the `Taxonomy` object, consisting of a map of terms and the [weighted pages] associated with each term.

[weighted pages]: /getting-started/glossary/#weighted-page

```go-html-template
{{ $taxonomyObject := .Data.Terms }} 
```

Once we have the `Taxonomy` object, we can call any of its [methods], allowing us to sort alphabetically or by term count.

[methods]: /methods/taxonomy/

## Sort alphabetically

The taxonomy template below inherits the site's shell from the base template, and renders a list of terms in the current taxonomy. Hugo sorts the list alphabetically by term, and displays the number of pages associated with each term.

{{< code file=layouts/_default/taxonomy.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Data.Terms.Alphabetical }}
    <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a> ({{ .Count }})</h2>
  {{ end }}
{{ end }}
{{< /code >}}

## Sort by term count

The taxonomy template below inherits the site's shell from the base template, and renders a list of terms in the current taxonomy. Hugo sorts the list by the number of pages associated with each term, and displays the number of pages associated with each term.

{{< code file=layouts/_default/taxonomy.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Data.Terms.ByCount }}
    <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a> ({{ .Count }})</h2>
  {{ end }}
{{ end }}
{{< /code >}}

## Include content links

The [`Alphabetical`] and [`ByCount`] methods used in the previous examples return an [ordered taxonomy], so we can also list the content to which each term is assigned.

[ordered taxonomy]: /getting-started/glossary/#ordered-taxonomy
[`Alphabetical`]: /methods/taxonomy/alphabetical/
[`ByCount`]: /methods/taxonomy/bycount/

The taxonomy template below inherits the site's shell from the base template, and renders a list of terms in the current taxonomy. Hugo sorts the list by the number of pages associated with each term, displays the number of pages associated with each term, then lists the content to which each term is assigned.

{{< code file=layouts/_default/taxonomy.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Data.Terms.ByCount }}
    <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a> ({{ .Count }})</h2>
    <ul>
      {{ range .WeightedPages }}
        <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
      {{ end }}
    </ul>
  {{ end }}
{{ end }}
{{< /code >}}

## Display metadata

Display metadata about each term by creating a corresponding branch bundle in the content directory.

For example, create an "authors" taxonomy:

{{< code-toggle file=hugo >}}
[taxonomies]
author = 'authors'
{{< /code-toggle >}}

Then create content with one [branch bundle] for each term:

[branch bundle]: /getting-started/glossary/#branch-bundle

```text
content/
└── authors/
    ├── jsmith/
    │   ├── _index.md
    │   └── portrait.jpg
    └── rjones/
        ├── _index.md
        └── portrait.jpg
```

Then add front matter to each term page:

{{< code-toggle file=content/authors/jsmith/_index.md fm=true >}}
title = "John Smith"
affiliation = "University of Chicago"
{{< /code-toggle >}}

Then create a taxonomy template specific to the "authors" taxonomy:

{{< code file=layouts/authors/taxonomy.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Data.Terms.Alphabetical }}
    <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a></h2>
    <p>Affiliation: {{ .Page.Params.Affiliation }}</p>
    {{ with .Page.Resources.Get "portrait.jpg" }}
      {{ with .Fill "100x100" }}
        <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="portrait">
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
{{< /code >}}

In the example above we list each author including their affiliation and portrait.
