---
title: images.Contrast
description: Returns an image filter that changes the contrast of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Contrast PERCENTAGE]
toc: true
---

The percentage must be in the range [-100, 100] where 0 has no effect. A value of `-100` produces a solid grey image, and a value of `100` produces an over-contrasted image.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Contrast -20 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Contrast"
  filterArgs="-20"
  example=true
>}}
