---
title: Filter
description: Applicable to images, applies one or more image filters to the given image resource.
categories: []
keywords: [filter]
params:
  alt_title: RESOURCE.Filter
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Filter FILTER...]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Filter` method returns a new resource from a [processable image](g) after applying one or more [image filters](#image-filters).

> [!note]
> Use the [`reflect.IsImageResourceProcessable`][] function to verify that an image can be processed.

## Usage

Use the `Filter` method to apply effects such as blurring, sharpening, or grayscale conversion. You can pass a single filter or a slice of filters. When providing a slice, Hugo applies the filters from left to right.

To apply a single filter:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter images.Grayscale }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

To apply multiple filters:

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

You can also apply image filters using the [`images.Filter`][] function.

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

{{% render-list-of-pages-in-section path=/functions/images filter=functions_images_no_filters filterType=exclude %}}

[`images.Filter`]: /functions/images/filter/
[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
