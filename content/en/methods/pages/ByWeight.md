---
title: ByWeight
description: Returns the given page collection sorted by weight in ascending order.
categories: []
keywords: []
action:
  related: []
  returnType: page.Pages
  signatures: [PAGES.ByWeight]
---

Assign a [weight] to a page using the `weight` field in front matter. The weight must be a non-zero integer. Lighter items float to the top, while heavier items sink to the bottom. Unweighted or zero-weighted pages are placed at the end of the collection.

[weight]: /getting-started/glossary/#weight

```go-html-template
{{ range .Pages.ByWeight }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range .Pages.ByWeight.Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
