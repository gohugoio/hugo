---
title: resources.Copy
description: Returns a copy of the given resource at the target path.
categories: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.Copy TARGETPATH RESOURCE]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with resources.Copy "img/new-image-name.jpg" . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

The `TARGETPATH` is relative to the [site root](g).

> [!NOTE]
> Use the `resources.Copy` function with global, page, and remote resources.
