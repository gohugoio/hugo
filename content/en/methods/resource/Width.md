---
title: Width
description: Applicable to images, returns the width of the given resource.
categories: []
keywords: []
action:
  related:
    - methods/resource/Height
  returnType: int
  signatures: [RESOURCE.Width]
---

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Width }} → 600
{{ end }}
```

Use the `Width` and `Height` methods together when rendering an `img` element:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
