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

Use the [`reflect.IsImageResourceWithMeta`][] function to verify that Hugo can determine the dimensions before calling the `Width` method.

```go-html-template
{{ with resources.GetMatch "images/featured.*" }}
  {{ if reflect.IsImageResourceWithMeta . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ else }}
    <img src="{{ .RelPermalink }}" alt="">
  {{ end }}
{{ end }}
```

[`reflect.IsImageResourceWithMeta`]: /functions/reflect/isimageresourcewithmeta/
