---
title: ByLanguage
description: Returns the given page collection sorted by language in ascending order.
categories: []
keywords: []
action:
  related: []
  returnType: page.Pages
  signatures: [PAGES.ByLanguage]
---

```go-html-template
{{ range .Site.AllPages.ByLanguage }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Site.AllPages.ByLanguage.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
