---
title: Fill
description: Applicable to images, returns a new image resource cropped and resized according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Fill SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Crop and resize an image according to the given [processing specification][]. You must provide both width and height (such as `500x200`) within the specification. Unlike [`Resize`][], which may stretch the image, `Fill` maintains the original aspect ratio by cropping the image to the target ratio before resizing. The operation uses the [anchor](#anchor) and [resampling filter](#resampling-filter) provided, if any.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "500x200 TopRight lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"500x200 TopRight lanczos"` is the _processing specification_.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "500x200 TopRight lanczos webp q85" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fill 500x200 TopRight lanczos webp q85"
  example=true
>}}

[`Resize`]: /methods/resource/resize/
[processing specification]: #processing-specification
