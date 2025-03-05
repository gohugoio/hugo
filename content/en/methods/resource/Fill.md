---
title: Fill
description: Applicable to images, returns an image resource cropped and resized to the given dimensions.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Fill SPEC]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Crop and resize an image to match the given dimensions. You must provide both width and height.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "200x200 top webp q85 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fill 200x200 top webp q85 lanczos"
  example=true
>}}
