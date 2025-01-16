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

{{< new-in 0.141.0 >}}

The `images.Mask` filter applies a mask to an image. Black pixels in the mask make the corresponding areas of the base image transparent, while white pixels keep them opaque. Color images are converted to grayscale for masking purposes. The mask is automatically resized to match the dimensions of the base image.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Mask "images/mask.png" }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

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
