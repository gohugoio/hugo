---
title: PageNumber
description: Returns the current pager's number within the pager collection.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [PAGER.PageNumber]
---

Use the `PageNumber` method to build navigation between pagers.

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ with $paginator }}
  <ul>
    {{ range .Pagers }}
      <li><a href="{{ .URL }}">{{ .PageNumber }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```
