---
title: images.Pixelate
description: Returns an image filter that applies a pixelation effect to an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Pixelate SIZE]
toc: true
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Pixelate 4 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Pixelate"
  filterArgs="4"
  example=true
>}}
