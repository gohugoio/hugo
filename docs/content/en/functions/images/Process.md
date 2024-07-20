---
title: images.Process
description: Returns an image filter that processes the given image using the given specification.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
    - methods/resource/Process
  returnType: images.filter
  signatures: [images.Process SPEC]
toc: true
---

{{< new-in 0.119.0 >}}

This filter has the same options as the [`Process`] method on a `Resource` object, but using it as a filter may be more effective if you need to apply multiple filters to an image.

[`Process`]: /methods/resource/process/

The process specification is a space-delimited, case-insensitive list of one or more of the following in any sequence:

action
: Specify zero or one of `crop`, `fill`, `fit`, or `resize`. If you specify an action you must also provide dimensions. See&nbsp;[details](content-management/image-processing/#image-processing-methods).

```go-html-template
{{ $filter := images.Process "resize 300x" }}
```

dimensions
: Required if you specify an action. Provide width _or_ height when using `resize`, else provide both width _and_ height. See&nbsp;[details](/content-management/image-processing/#dimensions).

```go-html-template
{{ $filter := images.Process "crop 200x200" }}
```

anchor
: Use with the `crop` or `fill` action. Specify zero or one of `TopLeft`, `Top`, `TopRight`, `Left`, `Center`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`, or `Smart`. Default is `Smart`. See&nbsp;[details](/content-management/image-processing/#anchor).

```go-html-template
{{ $filter := images.Process "crop 200x200 center" }}
```

rotation
: Typically specify zero or one of `r90`, `r180`, or `r270`. Also supports arbitrary rotation angles. See&nbsp;[details](/content-management/image-processing/#rotation).

```go-html-template
{{ $filter := images.Process "r90" }}
{{ $filter := images.Process "crop 200x200 center r90" }}
```

target format
: Specify zero or one of `gif`, `jpeg`, `png`, `tiff`, or `webp`. See&nbsp;[details](/content-management/image-processing/#target-format).

```go-html-template
{{ $filter := images.Process "webp" }}
{{ $filter := images.Process "crop 200x200 center r90 webp" }}
```

quality
: Applicable to JPEG and WebP images. Optionally specify `qN` where `N` is an integer in the range [0, 100]. Default is `75`. See&nbsp;[details](/content-management/image-processing/#quality).

```go-html-template
{{ $filter := images.Process "q50" }}
{{ $filter := images.Process "crop 200x200 center r90 webp q50" }}
```

hint
: Applicable to WebP images and equivalent to the `-preset` flag for the [`cwebp`] encoder. Specify zero or one of `drawing`, `icon`, `photo`, `picture`, or `text`. Default is `photo`. See&nbsp;[details](/content-management/image-processing/#hint).

[`cwebp`]: https://developers.google.com/speed/webp/docs/cwebp


```go-html-template
{{ $filter := images.Process "webp" "icon" }}
{{ $filter := images.Process "crop 200x200 center r90 webp q50 icon" }}
```

background color
: When converting a PNG or WebP with transparency to a format that does not support transparency, optionally specify a background color using a 3-digit or a 6-digit hexadecimal color code. Default is `#ffffff` (white). See&nbsp;[details](/content-management/image-processing/#background-color).

```go-html-template
{{ $filter := images.Process "jpeg #000" }}
{{ $filter := images.Process "crop 200x200 center r90 q50 jpeg #000" }}
```

resampling filter
: Typically specify zero or one of `Box`, `Lanczos`, `CatmullRom`, `MitchellNetravali`, `Linear`, or `NearestNeighbor`. Other resampling filters are available. See&nbsp;[details](/content-management/image-processing/#resampling-filter).

```go-html-template
{{ $filter := images.Process "resize 300x lanczos" }}
{{ $filter := images.Process "resize 300x r90 q50 jpeg #000 lanczos" }}
```

## Usage

Create a filter:

```go-html-template
{{ $filter := images.Process "resize 256x q40 webp" }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="resize 256x q40 webp"
  example=true
>}}
