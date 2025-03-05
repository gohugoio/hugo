---
title: Width
description: Applicable to images, returns the width of the given resource.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [RESOURCE.Width]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

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
