---
title: Related content
description: List related content in "See Also" sections.
categories: []
keywords: []
aliases: [/content/related/,/related/,/content-management/related/]
---

Hugo uses a set of factors to identify a page's related content based on front matter parameters. This can be tuned to the desired set of indices and parameters or left to Hugo's default [related content configuration](/configuration/related-content/).

## List related content

To list up to 5 related pages (which share the same _date_ or _keyword_ parameters) is as simple as including something similar to this partial in your template:

```go-html-template {file="layouts/_partials/related.html" copy=true}
{{ with site.RegularPages.Related . | first 5 }}
  <p>Related content:</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

The `Related` method takes one argument which may be a `Page` or an options map. The options map has these options:

indices
: (`slice`) The indices to search within.

document
: (`page`) The page for which to find related content. Required when specifying an options map.

namedSlices
: (`slice`) The keywords to search for, expressed as a slice of `KeyValues` using the [`keyVals`] function.

fragments
: (`slice`) A list of special keywords that is used for indices configured as type "fragments". This will match the [fragment](g) identifiers of the documents.

A fictional example using all of the above options:

```go-html-template
{{ $page := . }}
{{ $opts := dict
  "indices" (slice "tags" "keywords")
  "document" $page
  "namedSlices" (slice (keyVals "tags" "hugo" "rocks") (keyVals "date" $page.Date))
  "fragments" (slice "heading-1" "heading-2")
}}
```

> [!note]
> We improved and simplified this feature in Hugo 0.111.0. Before this we had 3 different methods: `Related`, `RelatedTo` and `RelatedIndices`. Now we have only one method: `Related`. The old methods are still available but deprecated. Also see [this blog article](https://regisphilibert.com/blog/2018/04/hugo-optmized-relashionships-with-related-content/) for a great explanation of more advanced usage of this feature.

## Index content headings

Hugo can index the headings in your content and use this to find related content. You can enable this by adding a index of type `fragments` to your `related` configuration:

{{< code-toggle file=hugo >}}
[related]
threshold    = 20
includeNewer = true
toLower      = false
[[related.indices]]
name        = "fragmentrefs"
type        = "fragments"
applyFilter = true
weight      = 80
{{< /code-toggle >}}

- The `name` maps to a optional front matter slice attribute that can be used to link from the page level down to the fragment/heading level.
- If `applyFilter` is enabled, the `.HeadingsFiltered` on each page in the result will reflect the filtered headings. This is useful if you want to show the headings in the related content listing:

```go-html-template
{{ $related := .Site.RegularPages.Related . | first 5 }}
{{ with $related }}
  <h2>See Also</h2>
  <ul>
    {{ range $i, $p := . }}
      <li>
        <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
        {{ with .HeadingsFiltered }}
          <ul>
            {{ range . }}
              {{ $link := printf "%s#%s" $p.RelPermalink .ID | safeURL }}
              <li>
                <a href="{{ $link }}">{{ .Title }}</a>
              </li>
            {{ end }}
          </ul>
        {{ end }}
      </li>
    {{ end }}
  </ul>
{{ end }}
```

## Configuration

See [configure related content](/configuration/related-content/).

[`keyVals`]: /functions/collections/keyvals/
