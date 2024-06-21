---
title: Page
description: Returns the taxonomy page or nil if the taxonomy has no terms.
categories: []
keywords: []
action:
  related: []
  returnType: page.Page
  signatures: [TAXONOMY.Page]
---

{{< new-in 0.125.0 >}}

This `TAXONOMY` method returns nil if the taxonomy has no terms, so you must code defensively:

```go-html-template
{{ with .Site.Taxonomies.tags.Page }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

This is rendered to:

```html
<a href="/tags/">Tags</a>
```
