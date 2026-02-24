---
title: AllPages
description: Returns a collection of all pages in all languages.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [SITE.AllPages]
expiryDate: '2028-02-18' # deprecated 2026-02-18 in v0.156.0
---

{{< deprecated-in 0.156.0 >}}
See [details](https://discourse.gohugo.io/t/56732).
{{< /deprecated-in >}}

This method returns all page [kinds](g) in all languages, in the [default sort order](g). That includes the home page, section pages, taxonomy pages, term pages, and regular pages.

In most cases you should use the [`RegularPages`] method instead.

[`RegularPages`]: /methods/site/regularpages/

```go-html-template
{{ range .Site.AllPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
