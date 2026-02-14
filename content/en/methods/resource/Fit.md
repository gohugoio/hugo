---
title: Fit
description: Applicable to images, returns a new image resource downscaled to fit according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Fit SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Downscale an image to fit according to the given [processing specification][] while maintaining the aspect ratio. You must provide both width and height (such as `600x400`) within the specification. Unlike [`Fill`][] or [`Resize`][], this method will never upscale an image; if the source image is smaller than the target dimensions, it remains its original size. The operation uses the [resampling filter](#resampling-filter) provided, if any.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "300x175 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"300x175 lanczos"` is the _processing specification_.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "300x175 lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fit 300x175 lanczos"
  example=true
>}}

[`Resize`]: /methods/resource/resize/
[`Fill`]: /methods/resource/fill/
[processing specification]: #processing-specification
