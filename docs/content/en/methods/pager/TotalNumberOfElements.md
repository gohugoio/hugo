---
title: TotalNumberOfElements
description: Returns the number of pages in the pager collection.
categories: []
keywords: []
action:
  related:
    - methods/pager/NumberOfElements
    - methods/page/Paginate
  returnType: int
  signatures: [PAGER.TotalNumberOfElements]
---

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ with $paginator }}
  {{ .TotalNumberOfElements }}
{{ end }}
```
