---
title: reflect.IsImageResourceProcessable
description: Reports whether the given value is a Resource object representing an image from which Hugo can extract dimensions and perform processing such as converting, resizing, cropping, or filtering.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsImageResourceProcessable INPUT]
---

{{< new-in 0.157.0 />}}

{{% glossary-term "processable image" %}}

## Usage

This example iterates through all project resources and uses `reflect.IsImageResourceProcessable` to ensure the image pipeline can perform transformations like resizing before processing begins.

```go-html-template
{{ range resources.Match "**" }}
  {{ if reflect.IsImageResourceProcessable . }}
    {{ with .Process "resize 300x webp" }}
      <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="Processed Image">
    {{ end }}
  {{ end }}
{{ end }}
```

{{% include "/_common/functions/reflect/image-reflection-functions.md" %}}
