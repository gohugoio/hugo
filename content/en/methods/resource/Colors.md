---
title: Colors
description: Applicable to images, returns a slice of the most dominant colors using a simple histogram method.
categories: []
keywords: []
action:
  related: []
  returnType: '[]string'
  signatures: [RESOURCE.Colors]
---

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Colors }} → [#bebebd #514947 #768a9a #647789 #90725e #a48974]
{{ end }}
```

This method is fast, but if you also scale down your images, it would be good for performance to extract the colors from the scaled image.

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
