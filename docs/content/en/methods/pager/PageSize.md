---
title: PageSize
description: Returns the maximum number of pages per pager.
categories: []
keywords: []
action:
  related:
    - methods/page/Paginate
  returnType: int
  signatures: [PAGER.PageSize]
---

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ with $paginator }}
  {{ .PageSize }}
{{ end }}
```
