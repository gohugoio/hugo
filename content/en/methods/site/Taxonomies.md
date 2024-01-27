---
title: Taxonomies
description: Returns a data structure containing the site's taxonomy objects, the terms within each taxonomy object, and the pages to which the terms are assigned.
categories: []
keywords: []
action:
  related: []
  returnType: page.TaxonomyList
  signatures: [SITE.Taxonomies]
---

Conceptually, the `Taxonomies` method on a `Site` object returns a data structure such&nbsp;as:

{{< code-toggle >}}
taxonomy a:
  - term 1:
    - page 1
    - page 2
  - term 2:
    - page 1
taxonomy b:
  - term 1:
    - page 2
  - term 2:
    - page 1
    - page 2
{{< /code-toggle >}}

For example, on a book review site you might create two taxonomies; one for genres and another for authors.

With this site configuration:

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

Conceptually, the taxonomies data structure looks like:

{{< code-toggle >}}
genres:
  - suspense:
    - And Then There Were None
    - Death on the Nile
    - Jamaica Inn
  - romance:
    - Jamaica Inn
    - Pride and Prejudice
authors:
  - achristie:
    - And Then There Were None
    - Death on the Nile
  - ddmaurier:
    - Jamaica Inn
  - jausten:
    - Pride and Prejudice
{{< /code-toggle >}}


To list the "suspense" books:

```go-html-template
<ul>
  {{ range .Site.Taxonomies.genres.suspense }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Hugo renders this to:

```html
<ul>
  <li><a href="/books/and-then-there-were-none/">And Then There Were None</a></li>
  <li><a href="/books/death-on-the-nile/">Death on the Nile</a></li>
  <li><a href="/books/jamaica-inn/">Jamaica Inn</a></li>
</ul>
```

{{% note %}}
Hugo's taxonomy system is powerful, allowing you to classify content and create relationships between pages.

Please see the [taxonomies] section for a complete explanation and examples.

[taxonomies]: content-management/taxonomies/
{{% /note %}}
