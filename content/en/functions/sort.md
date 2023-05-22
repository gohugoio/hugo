---
title: sort
description: Sorts slices, maps, and page collections.
categories: [functions]
signature: ["sort COLLECTION [KEY] [ORDER]"]
menu:
  docs:
    parent: functions
keywords: [ordering,sorting,lists]
toc: true
---

The `KEY` is optional when sorting slices in ascending order, otherwise it is required. When sorting slices, use the literal `value` in place of the `KEY`. See examples below.

The `ORDER` may be either `asc` (ascending) or `desc` (descending). The default sort order is ascending.

## Sort a slice

The examples below assume this site configuration:

{{< code-toggle file="config" copy=false >}}
[params]
grades = ['b','a','c']
{{< /code-toggle >}}

### Ascending order {#slice-ascending-order}

Sort slice elements in ascending order using either of these constructs:

{{< code file="layouts/_default/single.html" copy=false >}}
{{ sort site.Params.grades }} → [a b c]
{{ sort site.Params.grades "value" "asc" }} → [a b c]
{{< /code >}}

In the examples above, `value` is the `KEY` representing the value of the slice element.

### Descending order {#slice-descending-order}

Sort slice elements in descending order:

{{< code file="layouts/_default/single.html" copy=false >}}
{{ sort site.Params.grades "value" "desc" }} → [c b a]
{{< /code >}}

In the example above, `value` is the `KEY` representing the value of the slice element.

## Sort a map

The examples below assume this site configuration:

{{< code-toggle file="config" copy=false >}}
[params.authors.a]
firstName = "Marius"
lastName  = "Pontmercy"
[params.authors.b]
firstName = "Victor"
lastName  = "Hugo"
[params.authors.c]
firstName = "Jean"
lastName  = "Valjean"
{{< /code-toggle >}}

{{% note %}}
When sorting maps, the `KEY` argument must be lowercase.
{{% /note %}}

### Ascending order {#map-ascending-order}

Sort map objects in ascending order using either of these constructs:

{{< code file="layouts/_default/single.html" copy=false >}}
{{ range sort site.Params.authors "firstname" }}
  {{ .firstName }}
{{ end }}

{{ range sort site.Params.authors "firstname" "asc" }}
  {{ .firstName }}
{{ end }}
{{< /code >}}

These produce:

```text
Jean Marius Victor
```

### Descending order {#map-descending-order}

Sort map objects in descending order:

{{< code file="layouts/_default/single.html" copy=false >}}
{{ range sort site.Params.authors "firstname" "desc" }}
  {{ .firstName }}
{{ end }}
{{< /code >}}

This produces:

```text
Victor Marius Jean
```

## Sort a page collection

Although you can use the `sort` function to sort a page collection, Hugo provides [built-in methods for sorting page collections] by:

- weight
- linktitle
- title
- front matter parameter
- date
- expiration date
- last modified date
- publish date
- length

In this contrived example, sort the site's regular pages by `.Type` in descending order:

{{< code file="layouts/_default/home.html" copy=false >}}
{{ range sort site.RegularPages "Type" "desc" }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
{{< /code >}}


[built-in methods for sorting page collections]: /templates/lists/#order-content
