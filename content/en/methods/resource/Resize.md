---
title: Resize
description: Applicable to images, returns a new image resource resized according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Resize SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Resize an image according to the given [processing specification][]. You may specify only the width (such as `300x`) or only the height (`such as x150`) for proportional scaling. If you specify both width and height (such as `300x150`), the resulting image will be scaled to those exact dimensions; if the aspect ratio differs from the original, the image will be non-proportionally scaled (stretched or squashed). The operation uses the [resampling filter](#resampling-filter) provided, if any.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Resize "300x lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"300x lanczos"` is the _processing specification_.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Resize "300x lanczos" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="resize 300x lanczos"
  example=true
>}}

[processing specification]: #processing-specification
