---
title: ByTitle
description: Returns the given page collection sorted by title in ascending order.
categories: []
keywords: []
action:
  related:
    - methods/pages/ByLinkTitle
    - methods/pages/ByParam
  returnType: page.Pages
  signatures: [PAGES.ByTitle]
---

```go-html-template
{{ range .Pages.ByTitle }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByTitle.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```
