---
title: images.Invert
description: Returns an image filter that negates the colors of an image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Invert]
---

## Usage

Create the filter:

```go-html-template
{{ $filter := images.Invert }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Invert"
  filterArgs=""
  example=true
>}}
