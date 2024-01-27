---
title: images.GaussianBlur
description: Returns an image filter that applies a gaussian blur to an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.GaussianBlur SIGMA]
toc: true
---

The sigma value must be positive, and indicates how much the image will be blurred. The blur-affected radius is approximately 3 times the sigma value.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.GaussianBlur 5 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="GaussianBlur"
  filterArgs="5"
  example=true
>}}
