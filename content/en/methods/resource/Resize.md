---
title: Resize
description: Applicable to images, returns an image resource resized to the given width and/or height.
categories: []
keywords: []
action:
  related:
    - methods/resource/Crop
    - methods/resource/Fit
    - methods/resource/Fill
    - methods/resource/Process
    - functions/images/Process
  returnType: images.ImageResource
  signatures: [RESOURCE.Resize SPEC]
---

Resize an image to the given width and/or height.

If you specify both width and height, the resulting image will be disproportionally scaled unless the original image has the same aspect ratio.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Resize "300x" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

{{% include "/methods/resource/_common/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Resize "300x webp q85 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="resize 300x webp q85 lanczos"
  example=true
>}}
