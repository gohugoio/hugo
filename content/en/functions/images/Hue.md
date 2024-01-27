---
title: images.Hue
description: Returns an image filter that rotates the hue of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.Hue SHIFT]
toc: true
---

The hue angle shift is typically in the range [-180, 180] where 0 has no effect.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Hue -15 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Hue"
  filterArgs="-15"
  example=true
>}}
