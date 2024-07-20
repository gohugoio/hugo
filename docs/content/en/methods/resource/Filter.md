---
title: Filter
description: Applicable to images, applies one or more image filters to the given image resource.
categories: []
keywords: []
action:
  related:
    - functions/images/Filter
  returnType: resources.resourceAdapter
  signatures: [RESOURCE.Filter FILTER...]
toc: true
---

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

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

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

{{< list-pages-in-section path=/functions/images filter=functions_images_no_filters filterType=exclude >}}
