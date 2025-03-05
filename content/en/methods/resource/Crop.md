---
title: Crop
description: Applicable to images, returns an image resource cropped to the given dimensions without resizing.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Crop SPEC]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Crop an image to match the given dimensions without resizing. You must provide both width and height.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Crop "200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Crop "200x200 topright webp q85 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="crop 200x200 topright webp q85 lanczos"
  example=true
>}}
