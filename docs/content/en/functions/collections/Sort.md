---
title: collections.Sort
description: Sorts slices, maps, and page collections.
categories: []
keywords: []
action:
  aliases: [sort]
  related:
    - functions/collections/Reverse
    - functions/collections/Shuffle
    - functions/collections/Uniq
  returnType: any
  signatures: ['collections.Sort COLLECTION [KEY] [ORDER]']
toc: true
aliases: [/functions/sort]
---

The `KEY` is optional when sorting slices in ascending order, otherwise it is required. When sorting slices, use the literal `value` in place of the `KEY`. See examples below.

The `ORDER` may be either `asc` (ascending) or `desc` (descending). The default sort order is ascending.

## Sort a slice

The examples below assume this site configuration:

{{< code-toggle file=hugo >}}
[params]
grades = ['b','a','c']
{{< /code-toggle >}}

### Ascending order {#slice-ascending-order}

Sort slice elements in ascending order using either of these constructs:

```go-html-template
{{ sort site.Params.grades }} → [a b c]
{{ sort site.Params.grades "value" "asc" }} → [a b c]
```

In the examples above, `value` is the `KEY` representing the value of the slice element.

### Descending order {#slice-descending-order}

Sort slice elements in descending order:

```go-html-template
{{ sort site.Params.grades "value" "desc" }} → [c b a]
```

In the example above, `value` is the `KEY` representing the value of the slice element.

## Sort a map

The examples below assume this site configuration:

{{< code-toggle file=hugo >}}
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

```go-html-template
{{ range sort site.Params.authors "firstname" }}
  {{ .firstName }}
{{ end }}

{{ range sort site.Params.authors "firstname" "asc" }}
  {{ .firstName }}
{{ end }}
```

These produce:

```text
Jean Marius Victor
```

### Descending order {#map-descending-order}

Sort map objects in descending order:

```go-html-template
{{ range sort site.Params.authors "firstname" "desc" }}
  {{ .firstName }}
{{ end }}
```

This produces:

```text
Victor Marius Jean
```

### First level key removal

Hugo removes the first level keys when sorting a map.

Original map:

```json
{
  "felix": {
    "breed": "malicious",
    "type": "cat"
  },
  "spot": {
    "breed": "boxer",
    "type": "dog"
  }
}
```

After sorting:

```json
[
  {
    "breed": "malicious",
    "type": "cat"
  },
  {
    "breed": "boxer",
    "type": "dog"
  }
]
```

## Sort a page collection

{{% note %}}
Although you can use the `sort` function to sort a page collection, Hugo provides [sorting and grouping methods] as well.

[sorting and grouping methods]: /methods/pages/
{{% /note %}}

In this contrived example, sort the site's regular pages by `.Type` in descending order:

```go-html-template
{{ range sort site.RegularPages "Type" "desc" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
