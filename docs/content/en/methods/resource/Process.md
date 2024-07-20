---
title: Process
description: Applicable to images, returns an image resource processed with the given specification.
categories: []
keywords: []
action:
  related:
    - methods/resource/Crop
    - methods/resource/Fit
    - methods/resource/Fill
    - methods/resource/Resize
    - functions/images/Process
  returnType: images.ImageResource
  signatures: [RESOURCE.Process SPEC]
toc: true
---

Process an image with the given specification. The specification can contain an optional action, one of `crop`, `fill`, `fit`, or `resize`. This means that you can use this method instead of [`Crop`], [`Fill`], [`Fit`], or [`Resize`].

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Process "crop 200x200" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also use this method to apply simple transformations such as rotation and conversion:

```go-html-template
{{/* Rotate 90 degrees counter-clockwise. */}}
{{ $image := $image.Process "r90" }}

{{/* Convert to WebP. */}}
{{ $image := $image.Process "webp" }}
```

The `Process` method is also available as a filter, which is more effective if you need to apply multiple filters to an image. See [`images.Process`].

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

{{% include "/methods/resource/_common/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Process "crop 200x200 topright webp q85 lanczos" }}
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

[`Crop`]: /methods/resource/crop/
[`Fill`]: /methods/resource/fill/
[`Fit`]: /methods/resource/fit/
[`Resize`]: /methods/resource/resize/
[`images.Process`]: /functions/images/process/
