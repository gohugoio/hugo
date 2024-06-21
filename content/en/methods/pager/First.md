---
title: First
description: Returns the first pager in the pager collection.
categories: []
keywords: []
action:
  related:
    - methods/pager/Last
    - methods/pager/Prev
    - methods/pager/Next
    - methods/pager/HasPrev
    - methods/pager/HasNext
    - methods/page/Paginate
  returnType: page.Pager
  signatures: [PAGER.First]
---

Use the `First` method to build navigation between pagers.

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ with $paginator }}
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
