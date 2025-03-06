---
title: Pages
description: Returns a collection of all pages.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [SITE.Pages]
---

This method returns all page [kinds](g) in the current language, in the [default sort order](g). That includes the home page, section pages, taxonomy pages, term pages, and regular pages.

In most cases you should use the [`RegularPages`] method instead.

[`RegularPages`]: /methods/site/regularpages/

```go-html-template
{{ range .Site.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
