---
title: Height
description: Applicable to images, returns the height of the given resource.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [RESOURCE.Height]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

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
