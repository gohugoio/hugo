---
title: images.Brightness
description: Returns an image filter that changes the brightness of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Brightness PERCENTAGE]
toc: true
---

The percentage must be in the range [-100, 100] where 0 has no effect. A value of `-100` produces a solid black image, and a value of `100` produces a solid white image.

## Usage

Create the image filter:

```go-html-template
{{ $filter := images.Brightness 12 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Brightness"
  filterArgs="12"
  example=true
>}}
