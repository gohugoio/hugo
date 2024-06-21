---
title: PageGroups
description: Returns the page groups in the current pager.
categories: []
keywords: []
action:
  related:
    - methods/page/Paginate
  returnType: page.PagesGroup
  signatures: [PAGER.PageGroups]
---

Use the `PageGroups` method with any of the [grouping methods].

[grouping methods]: /quick-reference/page-collections/#group

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate ($pages.GroupByDate "Jan 2006") }}

{{ range $paginator.PageGroups }}
  <h2>{{ .Key }}</h2>
  {{ range .Pages }}
    <h3><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h3>
  {{ end }}
{{ end }}

{{ template "_internal/pagination.html" . }}
```
