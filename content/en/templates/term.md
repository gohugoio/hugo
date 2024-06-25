---
title: Term templates
description: Create a term template to render a list of pages associated with the current term.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 100
weight: 100
toc: true
---

The [term] template below inherits the site's shell from the [base template], and renders a list of pages associated with the current term.

[term]: /getting-started/glossary/#term
[base template]: /templates/types/

{{< code file=layouts/_default/term.html >}}
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

In the example above, the term will be capitalized if its respective page is not backed by a file. You can disable this in your site configuration:

{{< code-toggle file=hugo >}}
capitalizeListTitles = false
{{< /code-toggle >}}

## Data object

Use these methods on the `Data` object within a term template.

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

Term
: (`string`) Returns the name of the term.

```go-html-template
{{ .Data.Term }} → fiction
```

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

Then create a term template specific to the "authors" taxonomy:

{{< code file=layouts/authors/term.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  <p>Affiliation: {{ .Params.affiliation }}</p>
  {{ with .Resources.Get "portrait.jpg" }}
    {{ with .Fill "100x100" }}
      <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="portrait">
    {{ end }}
  {{ end }}
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
{{< /code >}}

In the example above we display the author with their affiliation and portrait, then a list of associated content.
