---
title: PageSize
description: Returns the number of pages per pager.
categories: []
keywords: []
action:
  related:
    - methods/page/Paginate
  returnType: int
  signatures: [PAGER.PageSize]
expiryDate: 2025-06-09 # deprecated 2024-06-09
---

{{% deprecated-in 0.128.0 %}}
Use [`PAGER.PagerSize`] instead.

[`PAGER.PagerSize`]: /methods/pager/pagersize/
{{% /deprecated-in %}}

The number of pages per pager is determined by the optional second argument passed to the [`Paginate`] method, falling back to the `pagerSize` as defined in your [site configuration].

[`Paginate`]: /methods/page/paginate/
[site configuration]: /templates/pagination/#configuration

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
