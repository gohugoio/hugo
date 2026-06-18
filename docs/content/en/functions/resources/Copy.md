---
title: resources.Copy
description: Copies the given resource to the target path.
categories: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.Copy TARGETPATH RESOURCE]
---

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with resources.Copy "img/new-image-name.jpg" . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

The `TARGETPATH` is relative to the server root. A leading slash is optional and has no effect.

> [!NOTE]
> Use the `resources.Copy` function with global, page, and remote resources.
