---
title: ByExpiryDate
description: Returns the given page collection sorted by expiration date in ascending order.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGES.ByExpiryDate]
---

When sorting by expiration date, the value is determined by your [site configuration], defaulting to the `expiryDate` field in front matter.

[site configuration]: /configuration/front-matter/#dates

```go-html-template
{{ range .Pages.ByExpiryDate }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByExpiryDate.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
