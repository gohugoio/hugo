---
title: images.Pixelate
description: Returns an image filter that applies a pixelation effect to an image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Pixelate SIZE]
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Pixelate 4 }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Pixelate"
  filterArgs="4"
  example=true
>}}
