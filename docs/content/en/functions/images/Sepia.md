---
title: images.Sepia
description: Returns an image filter that produces a sepia-toned version of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Sepia PERCENTAGE]
toc: true
---

The percentage must be in the range [0, 100] where 0 has no effect.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Sepia 75 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Sepia"
  filterArgs="75"
  example=true
>}}
