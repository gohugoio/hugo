---
title: Data
description: Returns a unique data object for each page kind.
categories: []
keywords: []
action:
  related: []
  returnType: page.Data
  signatures: [PAGE.Data]
toc: true
---

The `Data` method on a `Page` object returns a unique data object for each [page kind].

[page kind]: /getting-started/glossary/#page-kind

{{% note %}}
The `Data` method is only useful within [taxonomy] and [term] templates.

Themes that are not actively maintained may still use `.Data.Pages` in list templates. Although that syntax remains functional, use one of these methods instead: [`Pages`], [`RegularPages`], or [`RegularPagesRecursive`]

[`Pages`]: /methods/page/pages/
[`RegularPages`]: /methods/page/regularpages/
[`RegularPagesRecursive`]: /methods/page/regularpagesrecursive/
[term]: /getting-started/glossary/#term
[taxonomy]: /getting-started/glossary/#taxonomy
{{% /note %}}

The examples that follow are based on this site configuration:

{{< code-toggle file=hugo >}}
[taxonomies]
genre = 'genres'
author = 'authors'
{{< /code-toggle >}}

And this content structure:

```text
content/
├── books/
│   ├── and-then-there-were-none.md --> genres: suspense
│   ├── death-on-the-nile.md        --> genres: suspense
│   └── jamaica-inn.md              --> genres: suspense, romance
│   └── pride-and-prejudice.md      --> genres: romance
└── _index.md
```

## In a taxonomy template

Use these methods on the `Data` object within a taxonomy template.

Singular
: (`string`) Returns the singular name of the taxonomy.

```go-html-template
{{ .Data.Singular }} → genre
```

Plural
: (`string`) Returns the plural name of the taxonomy.

```go-html-template
{{ .Data.Plural }} → genres
```

Terms
: (`page.Taxonomy`) Returns the `Taxonomy` object, consisting of a map of terms and the [weighted pages] associated with each term.

```go-html-template
{{ $taxonomyObject := .Data.Terms }} 
```

{{% note %}}
Once you have captured the `Taxonomy` object, use any of the [taxonomy methods] to sort, count, or capture a subset of its weighted pages.

[taxonomy methods]: /methods/taxonomy/
{{% /note %}}

Learn more about [taxonomy templates].

## In a term template

Use these methods on the `Data` object within a term template.

Singular
: (`string`) Returns the singular name of the taxonomy.

```go-html-template
{{ .Data.Singular }} → genre
```

Plural
: (`string`) Returns the plural name of the taxonomy.

```go-html-template
{{ .Data.Plural }} → genres
```

Term
: (`string`) Returns the name of the term.

```go-html-template
{{ .Data.Term }} → suspense
```

Learn more about [term templates].

[taxonomy templates]: /templates/types/#taxonomy
[term templates]: /templates/types/#term
[weighted pages]: /getting-started/glossary/#weighted-page
