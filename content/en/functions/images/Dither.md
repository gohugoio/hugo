---
title: images.Dither
description: Returns an image filter that dithers an image.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - functions/images/Process
    - methods/resource/Colors
    - methods/resource/Filter
  returnType: images.filter
  signatures: ['images.Dither [OPTIONS]']
toc: true
---

{{< new-in 0.123.0 >}}

## Options

colors
: (`string array`) A slice of two or more colors that make up the dithering palette, each expressed as an RGB or RGBA [hexadecimal] value, with or without a leading hash mark. The default values are opaque black (`000000ff`) and opaque white (`ffffffff`).

[hexadecimal]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color

method
: (`string`) The dithering method. See the [dithering methods](#dithering-methods) section below for a list of the available methods. Default is `FloydSteinberg`.

serpentine
: (`bool`) Applicable to error diffusion dithering methods, serpentine controls whether the error diffusion matrix is applied in a serpentine manner, meaning that it goes right-to-left every other line. This greatly reduces line-type artifacts. Default is `true`.

strength
: (`float`) The strength at which to apply the dithering matrix, typically a value in the range [0, 1]. A value of `1.0` applies the dithering matrix at 100% strength (no modification of the dither matrix). The `strength` is inversely proportional to contrast; reducing the strength increases the contrast. Setting `strength` to a value such as `0.8` can be useful to reduce noise in the dithered image. Default is `1.0`.

## Usage

Create the options map:

```go-html-template
{{ $opts := dict
  "colors" (slice "222222" "808080" "dddddd")
  "method" "ClusteredDot4x4"
  "strength" 0.85
}}
```

Create the filter:

```go-html-template
{{ $filter := images.Dither $opts }}
```

Or create the filter using the default settings:

```go-html-template
{{ $filter := images.Dither }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Dithering methods

See the [Go documentation] for descriptions of each of the dithering methods below.

[Go documentation]: https://pkg.go.dev/github.com/makeworld-the-better-one/dither/v2#pkg-variables 

Error diffusion dithering methods:

- Atkinson
- Burkes
- FalseFloydSteinberg
- FloydSteinberg
- JarvisJudiceNinke
- Sierra
- Sierra2
- Sierra2_4A
- Sierra3
- SierraLite
- Simple2D
- StevenPigeon
- Stucki
- TwoRowSierra

Ordered dithering methods:

- ClusteredDot4x4
- ClusteredDot6x6
- ClusteredDot6x6_2
- ClusteredDot6x6_3
- ClusteredDot8x8
- ClusteredDotDiagonal16x16
- ClusteredDotDiagonal6x6
- ClusteredDotDiagonal8x8
- ClusteredDotDiagonal8x8_2
- ClusteredDotDiagonal8x8_3
- ClusteredDotHorizontalLine
- ClusteredDotSpiral5x5
- ClusteredDotVerticalLine
- Horizontal3x5
- Vertical5x3

## Example

This example uses the default dithering options.

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Dither"
  filterArgs=""
  example=true
>}}

## Recommendations

Regardless of dithering method, do both of the following to obtain the best results:

1. Scale the image _before_ dithering
2. Output the image to a lossless format such as GIF or PNG

The example below does both of these, and it sets the dithering palette to the three most dominant colors in the image.


```go-html-template
{{ with resources.Get "original.jpg" }}
  {{ $opts := dict
    "method" "ClusteredDotSpiral5x5"
    "colors" (first 3 .Colors)
  }}
  {{ $filters := slice
    (images.Process "resize 800x")
    (images.Dither $opts)
    (images.Process "png")
  }}
  {{ with . | images.Filter $filters }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

For best results, if the dithering palette is grayscale, convert the image to grayscale before dithering.

```go-html-template
{{ $opts := dict "colors" (slice "222" "808080" "ddd") }}
{{ $filters := slice
  (images.Process "resize 800x")
  (images.Grayscale)
  (images.Dither $opts)
  (images.Process "png")
}}
{{ with images.Filter $filters . }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

The example above:

1. Resizes the image to be 800 px wide
2. Converts the image to grayscale
3. Dithers the image using the default (`FloydSteinberg`) dithering method with a grayscale palette
4. Converts the image to the PNG format
