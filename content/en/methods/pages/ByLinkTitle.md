---
title: ByLinkTitle
description: Returns the given page collection sorted by link title in ascending order, falling back to title if link title is not defined.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGES.ByLinkTitle]
---

```go-html-template
{{ range .Pages.ByLinkTitle }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByLinkTitle.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
