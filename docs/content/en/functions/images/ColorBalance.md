---
title: images.ColorBalance
description: Returns an image filter that changes the color balance of an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.ColorBalance PCTRED PCTGREEN PCTBLUE]
toc: true
---

The percentage for each channel (red, green, blue) must be in the range [-100, 500].

## Usage

Create the filter:

```go-html-template
{{ $filter := images.ColorBalance -10 10 50 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="ColorBalance"
  filterArgs="-10,10,50"
  example=true
>}}
