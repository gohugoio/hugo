---
# Do not remove front matter.
---

Before we can use a `Taxonomy` method, we need to capture a `Taxonomy` object.

## Capture a taxonomy object

Consider this site configuration:

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

To capture the "genres" taxonomy object from within any template, use the [`Taxonomies`] method on a `Site` object.

```go-html-template
{{ $taxonomyObject := .Site.Taxonomies.genres }}
```

To capture the "genres" taxonomy object when rendering its page with a taxonomy template, use the [`Terms`] method on the page's [`Data`] object:

{{< code file=layouts/_default/taxonomy.html lang=go-html-template >}}
{{ $taxonomyObject := .Data.Terms }}
{{< /code >}}

To inspect the data structure:

```go-html-template
<pre>{{ jsonify (dict "indent" "  ") $taxonomyObject }}</pre>
```

Although the [`Alphabetical`] and [`ByCount`] methods provide a better data structure for ranging through the taxonomy, you can render the weighted pages by term directly from the `Taxonomy` object:

```go-html-template
{{ range $term, $weightedPages := $taxonomyObject }}
  <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a></h2>
  <ul>
    {{ range $weightedPages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

In the example above, the first anchor element is a link to the term page.


[`Alphabetical`]: /methods/taxonomy/alphabetical
[`ByCount`]: /methods/taxonomy/bycount

[`data`]: /methods/page/data
[`terms`]: /methods/page/data/#in-a-taxonomy-template
[`taxonomies`]: /methods/site/taxonomies
