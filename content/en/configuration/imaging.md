---
title: Configure imaging
linkTitle: Imaging
description: Configure imaging.
categories: []
keywords: []
---

## Processing options

These are the default settings for processing images:

{{< code-toggle file=hugo >}}
[imaging]
anchor = 'Smart'
bgColor = '#ffffff'
compression = 'lossy'
quality = 75
resampleFilter = 'box'
{{< /code-toggle >}}

anchor
: (`string`) The focal point used when cropping or filling an image. Valid options include `TopLeft`, `Top`, `TopRight`, `Left`, `Center`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`, or `Smart`. The `Smart` option utilizes the [`smartcrop.js`][] library to identify the most interesting area of the image. Default is `Smart`.

bgColor
: (string) The background color used when converting transparent images to formats that do not support transparency, such as PNG to JPEG. This color also fills the empty space created when rotating an image by a non-orthogonal angle if the space is not transparent and a background color is not specified in the  processing specification. The value must be an RGB [hexadecimal color][]. Default is `#ffffff`.

compression
: {{< new-in 0.153.5 />}}
: (`string`) The encoding strategy used for the image. Options are `lossy` or `lossless`. Note that `lossless` is only supported by the WebP format. Default is `lossy`.

quality
: (`int`) The visual fidelity of the image, applicable to JPEG and WebP formats when using `lossy` compression. Expressed as a whole number from `1` to `100`, inclusive. Lower numbers prioritize smaller file size, while higher numbers prioritize visual clarity. Default is `75`.

resampleFilter
: (`string`) The algorithm used to calculate new pixels when resizing, fitting, or filling an image. Common options include `box`, `lanczos`, `catmullRom`, `mitchellNetravali`, `linear`, or `nearestNeighbor`. Default is `box`.

  Filter|Description
  :--|:--
  `box`|Simple and fast averaging filter appropriate for downscaling
  `lanczos`|High-quality resampling filter for photographic images yielding sharp results
  `catmullRom`|Sharp cubic filter that is faster than the Lanczos filter while providing similar results
  `mitchellNetravali`|Cubic filter that produces smoother results with less ringing artifacts than CatmullRom
  `linear`|Bilinear resampling filter, produces smooth output, faster than cubic filters
  `nearestNeighbor`|Fastest resampling filter, no antialiasing

  Refer to the [source documentation][] for a complete list of available resampling filters. If you wish to improve image quality at the expense of performance, you may wish to experiment with the alternative filters.

## WebP images

{{< new-in 0.155.0 />}}

These are the default settings specific to processing WebP images:

{{< code-toggle file=hugo >}}
[imaging.webp]
hint = 'photo'
method = 4
useSharpYuv = true
{{< /code-toggle >}}

hint
: (`string`) The encoding preset used when processing WebP images, equivalent to the `-preset` flag for the [`cwebp`][] CLI. Valid options include `drawing`, `icon`, `photo`, `picture`, or `text`. Default is `photo`.

  Value|Example
  :--|:--
  `drawing`|Hand or line drawing with high-contrast details
  `icon`|Small colorful image
  `photo`|Outdoor photograph with natural lighting
  `picture`|Indoor photograph such as a portrait
  `text`|Image that is primarily text

method
: (`int`) The effort level of the compression algorithm. Expressed as a whole number from `0` to `6`, inclusive, equivalent to the `-m` flag for the [`cwebp`][] CLI. Lower numbers prioritize processing speed, while higher numbers prioritize compression efficiency. Default is `4`.

useSharpYuv
: (`bool`) The conversion method used for RGB-to-YUV encoding, equivalent to the `-sharp_yuv` flag for the [`cwebp`][] CLI. Enabling this prioritizes image sharpness at the expense of processing speed. Default is `true`.

## Exif method

These are the default settings for the [`Exif`] method on an image `Resource` object:

{{< code-toggle file=hugo >}}
[imaging.exif]
disableDate = false
disableLatLong = false
excludeFields = ""
includeFields = ""
{{< /code-toggle >}}

disableDate
: (`bool`) Whether to disable the [`Date`][] method by returning its zero value. Default is `false`.

disableLatLong
: (`bool`) Whether to disable the [`Lat`][] and [`Long`][] methods by returning their zero values. Default is `false`.

excludeFields
: (`string`) A [regular expression](g) matching the fields to exclude when extracting metadata.

  > [!note]
  > By default, to improve performance and decrease cache size, Hugo excludes the following fields: `ColorSpace`, `Contrast`, `Exif`, `ExposureBias`, `ExposureMode`, `ExposureProgram`, `Flash`, `GPS`, `JPEG`, `Metering`, `Resolution`, `Saturation`, `Sensing`, `Sharp`, and `WhiteBalance`.

includeFields
: (`string`) A [regular expression](g) matching the fields to include when extracting metadata. If empty, a default set excluding technical metadata is used. Set&nbsp;to&nbsp;`'.*'`&nbsp;to include all fields.

## Meta method

{{< new-in 0.155.0 />}}

These are the default settings for the [`Meta`] method on an image `Resource` object:

{{< code-toggle file=hugo >}}
[imaging.meta]
fields = []
sources = ['exif', 'iptc']
{{< /code-toggle >}}

fields
: (`[]string`) A [glob slice](g) matching the fields to include when extracting metadata. If empty, a default set excluding technical metadata is used. Set&nbsp;to&nbsp;`['**']`&nbsp;to include all fields.

  > [!note]
  > By default, to improve performance and decrease cache size, Hugo excludes the following fields: `ColorSpace`, `Contrast`, `Exif`, `ExposureBias`, `ExposureMode`, `ExposureProgram`, `Flash`, `GPS`, `JPEG`, `Metering`, `Resolution`, `Saturation`, `Sensing`, `Sharp`, and `WhiteBalance`.

sources
: (`[]string`) The metadata sources to include, one or more of `exif`, `iptc`, or `xmp`. Default is `['exif', 'iptc']`. The XMP metadata is excluded by default to improve performance.

[`cwebp`]: https://developers.google.com/speed/webp/docs/cwebp
[`Exif`]: /methods/resource/exif/
[`Meta`]: /methods/resource/meta/
[`smartcrop.js`]: https://github.com/jwagner/smartcrop.js
[hexadecimal color]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color
[source documentation]: https://github.com/disintegration/imaging#image-resizing
