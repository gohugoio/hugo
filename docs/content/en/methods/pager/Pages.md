---
title: Pages
description: Returns the pages in the current pager.
categories: []
keywords: []
action:
  related:
    - methods/page/Paginate
  returnType: page.Pages
  signatures: [PAGER.Pages]
---

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ template "_internal/pagination.html" . }}
```
