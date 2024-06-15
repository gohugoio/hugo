---
title: images.UnsharpMask
description: Returns an image filter that sharpens an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.UnsharpMask SIGMA AMOUNT THRESHOLD]
toc: true
---

The sigma argument is used in a gaussian function and affects the radius of effect. Sigma must be positive. The sharpen radius is approximately 3 times the sigma value.

The amount argument controls how much darker and how much lighter the edge borders become. Typically between 0.5 and 1.5.

The threshold argument controls the minimum brightness change that will be sharpened. Typically between 0 and 0.05.

## Usage

Create the filter:

```go-html-template
{{ $filter := images.UnsharpMask 10 0.4 0.03 }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="UnsharpMask"
  filterArgs="10,0.4,0.03"
  example=true
>}}
