---
title: Filter
description: Applicable to images, applies one or more image filters to the given image resource.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Filter FILTER...]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Apply one or more [image filters](#image-filters) to the given image.

To apply a single filter:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter images.Grayscale }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

To apply two or more filters, executing from left to right:

```go-html-template
{{ $filters := slice
  images.Grayscale
  (images.GaussianBlur 8)
}}
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter $filters }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also apply image filters using the [`images.Filter`] function.

[`images.Filter`]: /functions/images/filter/

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter images.Grayscale }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Grayscale"
  filterArgs=""
  example=true
>}}

## Image filters

Use any of these filters with the `Filter` method.

{{% list-pages-in-section path=/functions/images filter=functions_images_no_filters filterType=exclude %}}
