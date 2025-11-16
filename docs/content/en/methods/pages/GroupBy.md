---
title: GroupBy
description: Returns the given page collection grouped by the given field in ascending order.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.PagesGroup
    signatures: ['PAGES.GroupBy FIELD [SORT]']
---

{{% include "/_common/methods/pages/group-sort-order.md" %}}

```go-html-template
{{ range .Pages.GroupBy "Section" }}
  <p>{{ .Key }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

To sort the groups in descending order:

```go-html-template
{{ range .Pages.GroupBy "Section" "desc" }}
  <p>{{ .Key }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```
