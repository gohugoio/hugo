---
title: ByDate
description: Returns the given page collection sorted by date in ascending order.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGES.ByDate]
---

When sorting by date, the value is determined by your [site configuration], defaulting to the `date` field in front matter.

[site configuration]: /configuration/front-matter/#dates

```go-html-template
{{ range .Pages.ByDate }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByDate.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
