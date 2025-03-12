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
hint = 'photo'
quality = 75
resampleFilter = 'box'
{{< /code-toggle >}}

anchor
: (`string`) When using the [`Crop`] or [`Fill`] method, the anchor determines the placement of the crop box. One of `TopLeft`, `Top`, `TopRight`, `Left`, `Center`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`, or `Smart`. Default is `Smart`.

bgColor
: (`string`) The background color of the resulting image. Applicable when converting from a format that supports transparency to a format that does not support transparency, for example, when converting from PNG to JPEG. Expressed as an RGB [hexadecimal] value. Default is `#ffffff`.

[hexadecimal]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color

hint
: (`string`) Applicable to WebP images, this option corresponds to a set of predefined encoding parameters. One of `drawing`, `icon`, `photo`, `picture`, or `text`. Default is `photo`. See&nbsp;[details](/content-management/image-processing/#hint).

quality
: (`int`) Applicable to JPEG and WebP images, this value determines the quality of the converted image. Higher values produce better quality images, while lower values produce smaller files. Set this value to a whole number between `1` and `100`, inclusive. Default is `75`.

resampleFilter
: (`string`) The resampling filter used when resizing an image. Default is `box`. See&nbsp;[details](/content-management/image-processing/#resampling-filter)

## EXIF data

These are the default settings for extracting EXIF data from images:

{{< code-toggle file=hugo >}}
[imaging.exif]
includeFields = ""
excludeFields = ""
disableDate = false
disableLatLong = false
{{< /code-toggle >}}

disableDate
: (`bool`) Whether to disable extraction of the image creation date/time. Default is `false`.

disableLatLong
: (`bool`) Whether to disable extraction of the GPS latitude and longitude. Default is `false`.

excludeFields
: (`string`) A [regular expression](g) matching the tags to exclude when extracting EXIF data.

includeFields
: (`string`) A [regular expression](g) matching the tags to include when extracting EXIF data. To include all available tags, set this value to&nbsp;`".*"`.

> [!note]
> To improve performance and decrease cache size, Hugo excludes the following tags: `ColorSpace`, `Contrast`, `Exif`, `Exposure[M|P|B]`, `Flash`, `GPS`, `JPEG`, `Metering`, `Resolution`, `Saturation`, `Sensing`, `Sharp`, and `WhiteBalance`.
>
> To control tag availability, change the `excludeFields` or `includeFields` settings as described above.

[`Crop`]: /methods/resource/crop/
[`Fill`]: /methods/resource/fill/
