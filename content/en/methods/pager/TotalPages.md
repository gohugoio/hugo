---
title: TotalPages
description: Returns the number of pagers in the pager collection.
categories: []
keywords: []
action:
  related:
    - methods/pager/PageNumber
    - methods/page/Paginate
  returnType: int
  signatures: [PAGER.TotalPages]
---

Use the `TotalPages` method to build navigation between pagers.

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ with $paginator }}
  <p>Pager {{ .PageNumber }} of {{ .TotalPages }}</p>
  <ul>
    {{ with .First }}
      <li><a href="{{ .URL }}">First</a></li>
    {{ end }}
    {{ with .Prev }}
      <li><a href="{{ .URL }}">Previous</a></li>
    {{ end }}
    {{ with .Next }}
      <li><a href="{{ .URL }}">Next</a></li>
    {{ end }}
    {{ with .Last }}
      <li><a href="{{ .URL }}">Last</a></li>
    {{ end }}
  </ul>
{{ end }}
```
