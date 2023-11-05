---
title: images.Colorize
description: Returns an image filter that produces a colorized version of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Colorize HUE SATURATION PERCENTAGE]
toc: true
---

The hue is the angle on the color wheel, typically in the range [0, 360].

The saturation must be in the range [0, 100].

The percentage specifies the strength of the effect, and must be in the range [0, 100].

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Colorize 180 50 20 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Colorize"
  filterArgs="180,50,20"
  example=true
>}}
