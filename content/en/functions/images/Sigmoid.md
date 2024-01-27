---
title: images.Sigmoid
description: Returns an image filter that changes the contrast of an image using a sigmoidal function.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Sigmoid MIDPOINT FACTOR]
toc: true
---

This is a non-linear contrast change useful for photo adjustments; it preserves highlight and shadow detail.

The midpoint is the midpoint of contrast. It must be in the range [0, 1], typically 0.5.

The factor indicates how much to increase or decrease the contrast, typically in the range [-10, 10] where 0 has no effect. A positive value increases contrast, while a negative value decrease contrast.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Sigmoid 0.6 -4 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Sigmoid"
  filterArgs="0.6,-4"
  example=true
>}}
