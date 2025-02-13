---
title: images.Mask
description: Returns an image filter that applies a mask to the source image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Mask RESOURCE]
toc: true
---

{{< new-in 0.141.0 />}}

The `images.Mask` filter applies a mask to an image. Black pixels in the mask make the corresponding areas of the base image transparent, while white pixels keep them opaque. Color images are converted to grayscale for masking purposes. The mask is automatically resized to match the dimensions of the base image.

{{% note %}}
Of the formats supported by Hugo's imaging pipelie, only PNG and WebP have an alpha channel to support transparency. If your source image has a different format and you require transparent masked areas, convert it to either PNG or WebP as shown in the example below.
{{% /note %}}

When applying a mask to a non-transparent image format such as JPEG, the masked areas will be filled with the color specified by the `bgColor` parameter in your [site configuration]. You can override that color with a `Process` image filter:

```go-html-template
{{ $filter := images.Process "#00ff00" }}
```

[site configuration]: /content-management/image-processing/#imaging-configuration

## Usage

Create a slice of filters, one for WebP conversion and the other for mask application:

```go-html-template
{{ $filter1 := images.Process "webp" }}
{{ $filter2 := images.Mask (resources.Get "images/mask.png") }}
{{ $filters := slice $filter1 $filter2 }}
```

Apply the filters using the [`images.Filter`] function:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with . | images.Filter $filters }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also apply the filter using the [`Filter`] method on a 'Resource' object:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter $filters }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

[`images.Filter`]: /functions/images/filter/
[`Filter`]: /methods/resource/filter/

## Example

Mask

{{< img
  src="images/examples/mask.png"
  example=false
>}}

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="mask"
  filterArgs="images/examples/mask.png"
  example=true
>}}
