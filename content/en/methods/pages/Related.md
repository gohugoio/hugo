---
title: Related
description: Returns a collection of pages related to the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/HeadingsFiltered
    - functions/collections/KeyVals
  returnType: page.Pages
  signatures:
    - PAGES.Related PAGE
    - PAGES.Related OPTIONS
---

Based on front matter, Hugo uses several factors to identify content related to the given page. Use the default [related content configuration], or tune the results to the desired indices and parameters. See&nbsp;[details].

The argument passed to the `Related` method may be a `Page` or an options map. For example, to pass the current page:

{{< code file=layouts/_default/single.html  >}}
{{ with .Site.RegularPages.Related . | first 5 }}
  <p>Related pages:</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
{{< /code >}}

To pass an options map:

{{< code file=layouts/_default/single.html  >}}
{{ $opts := dict
  "document" .
  "indices" (slice "tags" "keywords")
}}
{{ with .Site.RegularPages.Related $opts | first 5 }}
  <p>Related pages:</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
{{< /code >}}


## Options

indices
: (`slice`) The indices to search within.

document
: (`page`) The page for which to find related content. Required when specifying an options map.

namedSlices
: (`slice`) The keywords to search for, expressed as a slice of `KeyValues` using the [`keyVals`] function.

[`keyVals`]: /functions/collections/keyvals/

fragments
: (`slice`) A list of special keywords that is used for indices configured as type "fragments". This will match the [fragment] identifiers of the documents.

A contrived example using all of the above:

```go-html-template
{{ $page := . }}
{{ $opts := dict
  "indices" (slice "tags" "keywords")
  "document" $page
  "namedSlices" (slice (keyVals "tags" "hugo" "rocks") (keyVals "date" $page.Date))
  "fragments" (slice "heading-1" "heading-2")
}}
```

[details]: /content-management/related/
[fragment]: /getting-started/glossary/#fragment
[related content configuration]: /content-management/related/
