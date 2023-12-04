---
title: Height
description: Applicable to images, returns the height of the given resource.
categories: []
keywords: []
action:
  related:
    - methods/resource/Width
  returnType: int
  signatures: [RESOURCE.Height]
---

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Height }} → 400
{{ end }}
```

Use the `Width` and `Height` methods together when rendering an `img` element:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
