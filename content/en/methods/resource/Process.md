---
title: Process
description: Applicable to images, returns a new image resource processed according to the given processing specification.
categories: []
keywords: [process]
params:
  alt_title: RESOURCE.Process
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Process SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Process an image according to the given [processing specification][]. This versatile method supports the full range of image transformations, including resizing, cropping, rotation, and format conversion, all within a single specification string.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Process "crop 200x200 TopRight webp q50" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"crop 200x200 TopRight webp q50"` is the _processing specification_.

You can also use this method to apply simple transformations such as rotation and conversion:

```go-html-template
{{/* Rotate 90 degrees counter-clockwise. */}}
{{ $image := $image.Process "r90" }}

{{/* Convert to WebP. */}}
{{ $image := $image.Process "webp" }}
```

The `Process` method is also available as a filter, which is more effective if you need to apply multiple filters to an image. See [`images.Process`].

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Process "crop 200x200 TopRight webp q50" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="crop 200x200 TopRight webp q50"
  example=true
>}}

[`images.Process`]: /functions/images/process/
[processing specification]: #processing-specification
