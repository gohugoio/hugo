---
title: Permalink
description: Publishes the given resource and returns its permalink.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [RESOURCE.Permalink]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Permalink` method on a `Resource` object writes the resource to the publish directory, typically `public`, and returns its [permalink](g).

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Permalink }} → https://example.org/images/a.jpg
{{ end }}
```
