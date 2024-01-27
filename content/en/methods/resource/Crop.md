---
title: Crop
description: Applicable to images, returns an image resource cropped to the given dimensions without resizing.
categories: []
keywords: []
action:
  related:
    - methods/resource/Fit
    - methods/resource/Fill
    - methods/resource/Resize
    - methods/resource/Process
    - functions/images/Process
  returnType: images.ImageResource
  signatures: [RESOURCE.Crop SPEC]
toc: true
---

Crop an image to match the given dimensions without resizing. You must provide both width and height.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Crop "200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

{{% include "/methods/resource/_common/processing-spec.md" %}}

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
