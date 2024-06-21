---
title: Pages
description: Returns a collection of all pages.
categories: []
keywords: []
action:
  related:
    - methods/site/AllPages
    - methods/site/RegularPages
    - methods/site/Sections
  returnType: page.Pages
  signatures: [SITE.Pages]
---

This method returns all page [kinds] in the current language. That includes the home page, section pages, taxonomy pages, term pages, and regular pages.

In most cases you should use the [`RegularPages`] method instead.

[`RegularPages`]: /methods/site/regularpages/
[kinds]: /getting-started/glossary/#page-kind

```go-html-template
{{ range .Site.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
