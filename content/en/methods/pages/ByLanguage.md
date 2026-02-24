---
title: ByLanguage
description: Returns the given page collection sorted by language.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGES.ByLanguage]
---

When sorting by language, Hugo orders the page collection using the following priority:

1. Language weight (ascending)
1. Date (descending)
1. LinkTitle (ascending)

This method is rarely, if ever, needed. Page collections that already contain multiple languages, such as those returned by the [`Rotate`][], [`Translations`][], or [`AllTranslations`][] methods on a `Page` object, are already sorted by language weight.

This contrived example aggregates pages from all sites and then sorts them by language:

```go-html-template
{{ $p := slice }}
{{ range hugo.Sites }}
  {{ range .Pages }}
    {{ $p = $p | append . }}
  {{ end }}
{{ end }}

{{ range $p.ByLanguage }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range $p.ByLanguage.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

[`AllTranslations`]: /methods/page/alltranslations/
[`Rotate`]: /methods/page/rotate/
[`Translations`]: /methods/page/translations/
