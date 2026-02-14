---
title: Crop
description: Applicable to images, returns a new image resource cropped according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Crop SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Crop an image according to the given [processing specification][]. When cropping, you must provide both width and height (such as `200x200`) within the specification. This method does not perform any resizing; it simply extracts a region of the image based on the dimensions and the [anchor](#anchor) provided, if any.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Crop "200x200 TopRight" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"200x200 TopRight"` is the _processing specification_.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Crop "200x200 TopRight" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="crop 200x200 TopRight"
  example=true
>}}

[processing specification]: #processing-specification
