---
title: ByLength
description: Returns the given page collection sorted by content length in ascending order.
categories: []
keywords: []
action:
  related: []
  returnType: page.Pages
  signatures: [PAGES.ByLength]
---

```go-html-template
{{ range .Pages.ByLength }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByLength.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
