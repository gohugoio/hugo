---
title: images.Grayscale
description: Returns an image filter that produces a grayscale version of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Grayscale]
toc: true
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Grayscale }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Grayscale"
  filterArgs=""
  example=true
>}}
