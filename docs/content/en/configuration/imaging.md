---
title: Configure imaging
linkTitle: Imaging
description: Configure imaging.
categories: []
keywords: []
---

These are the default settings for processing images:

{{< code-toggle config=imaging />}}

## Top-level settings

These settings apply to all image formats.

`anchor`
: (`string`) The focal point used when cropping or filling an image. Valid case-insensitive options include `TopLeft`, `Top`, `TopRight`, `Left`, `Center`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`, or `Smart`. The `Smart` option utilizes the [`muesli/smartcrop`][] package to identify the most interesting area of the image. Default is `smart`.

`bgColor`
: (`string`) The background color used when converting transparent images to formats that do not support transparency, such as PNG to JPEG. This color also fills the empty space created when rotating an image by a non-orthogonal angle if the space is not transparent and a background color is not specified in the processing specification. The value must be an RGB [hexadecimal color][]. Default is `#ffffff`.

`compression`
: {{< deprecated-in 0.163.0 />}}
: Use the format-specific `compression` setting instead, applicable to [AVIF](#avif) and [WebP](#webp) images.

`hint`
: {{< deprecated-in 0.163.0 />}}
: Use the format-specific `hint` setting instead, applicable to [AVIF](#avif) and [WebP](#webp) images.

`quality`
: {{< deprecated-in 0.163.0 />}}
: Use the format-specific `quality` setting instead, applicable to [AVIF](#avif), [JPEG](#jpeg), and [WebP](#webp) images.

`resampleFilter`
: (`string`) The algorithm used to calculate new pixels when resizing, fitting, or filling an image. Common case-insensitive options include `box`, `lanczos`, `catmullRom`, `mitchellNetravali`, `linear`, or `nearestNeighbor`. Default is `box`.

  Filter|Description
  :--|:--
  `box`|Simple and fast averaging filter appropriate for downscaling
  `lanczos`|High-quality resampling filter for photographic images yielding sharp results
  `catmullRom`|Sharp cubic filter that is faster than the Lanczos filter while providing similar results
  `mitchellNetravali`|Cubic filter that produces smoother results with less ringing artifacts than CatmullRom
  `linear`|Bilinear resampling filter, produces smooth output, faster than cubic filters
  `nearestNeighbor`|Fastest resampling filter, no antialiasing

  Refer to the [source documentation][] for a complete list of available resampling filters. If you wish to improve image quality at the expense of performance, you may wish to experiment with the alternative filters.

## AVIF

{{< new-in 0.162.0 />}}

These settings apply when encoding AVIF images.

> [!NOTE]
> When exporting HDR AVIF images from Lightroom, in the Export dialog under File Settings, uncheck Maximize Compatibility to improve Hugo's AVIF decoding speed.

> [!NOTE]
> Encoding animated images to AVIF produces a single-frame (static) image. Converting an animated AVIF to another format such as GIF works as expected.

{{< code-toggle config=imaging.avif />}}

`compression`
: {{< new-in 0.163.0 />}}
: (`string`) The encoding strategy. Options are `lossy` or `lossless`. Default is `lossy`.

`encoderSpeed`
: (`int`) The encoder speed. Expressed as a whole number from `1` to `10`, inclusive, equivalent to the `-s` flag for the [`avifenc`][] CLI. Lower numbers reduce file size at the cost of build time. At typical web image sizes, quality is indistinguishable across settings. Values below `5` may cause significantly longer build times. Default is `10`.

`hint`
: {{< new-in 0.163.0 />}}
: (`string`) The content hint. Valid options include `drawing`, `icon`, `photo`, `picture`, or `text`. Hugo uses the `4:2:0` chroma subsampling format with `photo` and `picture`, and `4:4:4` with the remaining options. Default is `photo`.

  Value|Example
  :--|:--
  `drawing`|Hand or line drawing with high-contrast details
  `icon`|Small colorful image
  `photo`|Outdoor photograph with natural lighting
  `picture`|Indoor photograph such as a portrait
  `text`|Image that is primarily text

`quality`
: {{< new-in 0.163.0 />}}
: (`int`) The visual fidelity when using `lossy` compression. Expressed as a whole number from `1` to `100`, inclusive. Lower numbers prioritize smaller file size, while higher numbers prioritize visual clarity. Default is `60`. Quality values are encoder-specific and not directly comparable across formats; a value of `60` for AVIF is perceptually similar to `75` for JPEG.

## JPEG

{{< new-in 0.163.0 />}}

These settings apply when encoding JPEG images.

{{< code-toggle config=imaging.jpeg />}}

`quality`
: (`int`) The visual fidelity. Expressed as a whole number from `1` to `100`, inclusive. Lower numbers prioritize smaller file size, while higher numbers prioritize visual clarity. Default is `75`.

## WebP

{{< new-in 0.155.0 />}}

These settings apply when encoding WebP images.

{{< code-toggle config=imaging.webp />}}

`compression`
: {{< new-in 0.163.0 />}}
: (`string`) The encoding strategy. Options are `lossy` or `lossless`. Default is `lossy`.

`hint`
: (`string`) The content hint, equivalent to the `-preset` flag for the [`cwebp`][] CLI. Valid options include `drawing`, `icon`, `photo`, `picture`, or `text`. Default is `photo`.

  Value|Example
  :--|:--
  `drawing`|Hand or line drawing with high-contrast details
  `icon`|Small colorful image
  `photo`|Outdoor photograph with natural lighting
  `picture`|Indoor photograph such as a portrait
  `text`|Image that is primarily text

`method`
: (`int`) The effort level of the compression algorithm. Expressed as a whole number from `0` to `6`, inclusive, equivalent to the `-m` flag for the [`cwebp`][] CLI. Lower numbers prioritize processing speed, while higher numbers prioritize compression efficiency and image quality. Default is `2`.

`quality`
: {{< new-in 0.163.0 />}}
: (`int`) The visual fidelity when using `lossy` compression. Expressed as a whole number from `1` to `100`, inclusive. Lower numbers prioritize smaller file size, while higher numbers prioritize visual clarity. Default is `75`.

`useSharpYuv`
: (`bool`) The conversion method used for RGB-to-YUV encoding, equivalent to the `-sharp_yuv` flag for the [`cwebp`][] CLI. Enabling this prioritizes image sharpness at the expense of processing speed. Default is `false`.

## Exif method

{{< deprecated-in 0.155.0 >}}
Use the [`Meta`](#meta-method) method instead.
{{< /deprecated-in >}}

## Meta method

{{< new-in 0.155.0 />}}

The following parameters allow you to control how Hugo extracts and filters metadata when using the [`Meta`][] method, helping you balance data granularity with build performance.

`fields`
: (`[]string`) A [glob slice](g) matching the fields to include when extracting metadata. If empty, a default set excluding technical metadata is used. Set&nbsp;to&nbsp;`['**']`&nbsp;to include all fields.

  > [!NOTE]
  > By default, to improve performance and decrease cache size, Hugo excludes the following fields: `ColorSpace`, `Contrast`, `Exif`, `ExposureBias`, `ExposureMode`, `ExposureProgram`, `Flash`, `GPS`, `JPEG`, `Metering`, `Resolution`, `Saturation`, `Sensing`, `Sharp`, and `WhiteBalance`.

`sources`
: (`[]string`) The metadata sources to include, one or more of `exif`, `iptc`, or `xmp`. Default is `['exif', 'iptc']`. The XMP metadata is excluded by default to improve performance.

[`avifenc`]: https://github.com/aomediacodec/libavif
[`cwebp`]: https://developers.google.com/speed/webp/docs/cwebp
[`muesli/smartcrop`]: https://github.com/muesli/smartcrop
[hexadecimal color]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color
[source documentation]: https://github.com/disintegration/imaging#image-resizing
