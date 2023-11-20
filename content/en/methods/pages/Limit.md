---
title: Limit
description: Returns the first N pages from the given page collection.
categories: []
keywords: []
action:
  related: []
  returnType: page.Pages
  signatures: [PAGES.Limit NUMBER]
---

```go-html-template
{{ range .Pages.Limit 3 }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
