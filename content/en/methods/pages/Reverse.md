---
title: Reverse
description: Returns the given page collection in reverse order.
categories: []
keywords: []
action:
  related: []
  returnType: page.Pages
  signatures: [PAGES.Reverse]
---

```go-html-template
{{ range .Pages.ByDate.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
