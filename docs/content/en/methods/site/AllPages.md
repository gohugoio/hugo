---
title: AllPages
description: Returns a collection of all pages in all languages.
categories: []
keywords: []
action:
  related:
    - methods/site/Pages
    - methods/site/RegularPages
    - methods/site/Sections
  returnType: page.Pages
  signatures: [SITE.AllPages]
---

This method returns all page [kinds](g) in all languages. That includes the home page, section pages, taxonomy pages, term pages, and regular pages.

In most cases you should use the [`RegularPages`] method instead.

[`RegularPages`]: /methods/site/regularpages/

```go-html-template
{{ range .Site.AllPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
