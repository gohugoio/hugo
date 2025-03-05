---
title: images.Grayscale
description: Returns an image filter that produces a grayscale version of an image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Grayscale]
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Grayscale }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Grayscale"
  filterArgs=""
  example=true
>}}
