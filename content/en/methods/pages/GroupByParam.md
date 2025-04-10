---
title: GroupByParam
description: Returns the given page collection grouped by the given parameter in ascending order.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.PagesGroup
    signatures: ['PAGES.GroupByParam PARAM [SORT]']
---

{{% include "/_common/methods/pages/group-sort-order.md" %}}

```go-html-template
{{ range .Pages.GroupByParam "color" }}
  <p>{{ .Key | title }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

To sort the groups in descending order:

```go-html-template
{{ range .Pages.GroupByParam "color" "desc" }}
  <p>{{ .Key | title }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```
