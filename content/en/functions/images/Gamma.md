---
title: images.Gamma
description: Returns an image filter that performs gamma correction on an image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Gamma GAMMA]
---

The gamma value must be positive. A value greater than 1 lightens the image, while a value less than 1 darkens the image. The filter has no effect when the gamma value is&nbsp;1.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Gamma 1.667 }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Gamma"
  filterArgs="1.667"
  example=true
>}}
