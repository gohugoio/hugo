---
title: ByDate
description: Returns the given page collection sorted by date in ascending order.
categories: []
keywords: []
action:
  related:
    - methods/pages/ByExpiryDate
    - methods/pages/ByLastMod
    - methods/pages/ByPublishDate
  returnType: page.Pages
  signatures: [PAGES.ByDate]
---

When sorting by date, the value is determined by your [site configuration], defaulting to the `date` field in front matter.

[site configuration]: /getting-started/configuration/#configure-dates

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
