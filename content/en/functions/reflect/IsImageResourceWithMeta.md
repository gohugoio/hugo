---
title: reflect.IsImageResourceWithMeta
description: Reports whether the given value is a Resource object representing an image from which Hugo can extract dimensions and, if present, Exif, IPTC, and XMP data.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsImageResourceWithMeta INPUT]
---

{{< new-in 0.157.0 />}}

## Usage

This example iterates through all project resources and uses `reflect.IsImageResourceWithMeta` to safely display image dimensions and metadata only for supported formats.

```go-html-template
{{ range resources.Match "**" }}
  {{ if reflect.IsImageResourceWithMeta . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="Image with Meta">
    {{ with .Meta }}
      <p>Taken on: {{ .Date }}</p>
    {{ end }}
  {{ end }}
{{ end }}
```

{{% include "/_common/functions/reflect/image-reflection-functions.md" %}}
