---
title: images.Saturation
description: Returns an image filter that changes the saturation of an image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Saturation PERCENTAGE]
---

The percentage must be in the range [-100, 500] where 0 has no effect.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Saturation 65 }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Saturation"
  filterArgs="65"
  example=true
>}}
