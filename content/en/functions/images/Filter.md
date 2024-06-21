---
title: images.Filter
description: Applies one or more image filters to the given image resource.
categories: []
keywords: []
action:
  aliases: []
  related:
    - methods/resource/Filter
  returnType: images.ImageResource
  signatures: [images.Filter FILTERS... IMAGE]
toc: true
---

Apply one or more [image filters](#image-filters) to the given image.

To apply a single filter:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with images.Filter images.Grayscale . }}
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
  {{ with images.Filter $filters . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also apply image filters using the [`Filter`] method on a `Resource` object.

[`Filter`]: /methods/resource/filter/

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with images.Filter images.Grayscale . }}
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

Use any of these filters with the `images.Filter` function, or with the `Filter` method on a `Resource` object.

{{< list-pages-in-section path=/functions/images filter=functions_images_no_filters filterType=exclude >}}
