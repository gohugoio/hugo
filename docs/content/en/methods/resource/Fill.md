---
title: Fill
description:  Applicable to images, returns an image resource cropped and resized to the given dimensions.
categories: []
keywords: []
action:
  related:
    - methods/resource/Crop
    - methods/resource/Fit
    - methods/resource/Resize
    - methods/resource/Process
    - functions/images/Process
  returnType: images.ImageResource
  signatures: [RESOURCE.Fill SPEC]
toc: true
---

Crop and resize an image to match the given dimensions. You must provide both width and height.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

{{% include "/methods/resource/_common/processing-spec.md" %}}

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
