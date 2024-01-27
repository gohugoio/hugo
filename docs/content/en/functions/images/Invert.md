---
title: images.Invert
description: Returns an image filter that negates the colors of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Invert]
toc: true
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Invert }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Invert"
  filterArgs=""
  example=true
>}}
