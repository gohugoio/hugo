---
title: ByLastmod
description: Returns the given page collection sorted by last modification date in ascending order.
categories: []
keywords: []
action:
  related:
    - methods/pages/ByDate
    - methods/pages/ByExpiryDate
    - methods/pages/ByPublishDate
  returnType: page.Pages
  signatures: [PAGES.ByLastmod]
---

When sorting by last modification date, the value is determined by your [site configuration], defaulting to the `lastmod` field in front matter.

[site configuration]: /getting-started/configuration/#configure-dates

```go-html-template
{{ range .Pages.ByLastmod }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByLastmod.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
