---
title: Fit
description: Applicable to images, returns an image resource downscaled to fit the given dimensions while maintaining aspect ratio. 
categories: []
keywords: []
action:
  related:
    - methods/resource/Crop
    - methods/resource/Fill
    - methods/resource/Resize
    - methods/resource/Process
    - functions/images/Process
  returnType: images.ImageResource
  signatures: [RESOURCE.Fit SPEC]
toc: true
---

Downscale an image to fit the given dimensions while maintaining aspect ratio. You must provide both width and height.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

{{% include "/methods/resource/_common/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "300x175 webp q85 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fit 300x175 webp q85 lanczos"
  example=true
>}}
