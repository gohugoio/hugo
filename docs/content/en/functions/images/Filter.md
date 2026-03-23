---
title: images.Filter
description: Applies one or more image filters to the given image resource.
categories: []
keywords: [filter]
params:
  functions_and_methods:
    aliases: []
    returnType: images.ImageResource
    signatures: [images.Filter FILTER... RESOURCE]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `images.Filter` function returns a new resource from a [processable image](g) after applying one or more [image filters](#image-filters).

> [!note]
> Use the [`reflect.IsImageResourceProcessable`][] function to verify that an image can be processed.

## Usage

Use the `images.Filter` function to apply effects such as blurring, sharpening, or grayscale conversion. You can pass a single filter or a slice of filters. When providing a slice, Hugo applies the filters from left to right.

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

You can also apply image filters using the [`Filter`][] method on a `Resource` object.

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

{{% render-list-of-pages-in-section path=/functions/images filter=functions_images_no_filters filterType=exclude %}}

[`Filter`]: /methods/resource/filter/
[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
