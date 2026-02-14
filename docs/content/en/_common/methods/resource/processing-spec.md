---
_comment: Do not remove front matter.
---

## Processing specification

The processing specification is a space-delimited, case-insensitive list containing one or more of the following options in any sequence:

action
: Specify one of `crop`, `fill`, `fit`, or `resize`. This is applicable to the [`Process`][] method and the [`images.Process`][] filter. If you specify an action, you must also provide dimensions.

anchor
: The focal point used when cropping or filling an image. Valid options include `TopLeft`, `Top`, `TopRight`, `Left`, `Center`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`, or `Smart`. The `Smart` option utilizes the [`smartcrop.js`][] library to identify the most interesting area of the image. This defaults to the [`anchor`][] parameter in your site configuration.

background color
: The background color used when converting transparent images to formats that do not support transparency, such as PNG to JPEG. This color also fills the empty space created when rotating an image by a non-orthogonal angle if the space is not transparent and a background color is not specified in the  processing specification. The value must be an RGB [hexadecimal color][]. This defaults to the [`bgColor`][] parameter in your site configuration.

compression
: {{< new-in 0.153.5 />}}
: The encoding strategy used for the image. Options are `lossy` or `lossless`. Note that `lossless` is only supported by the WebP format. This defaults to the [`compression`][] parameter in your site configuration.

dimensions
: The dimensions of the resulting image, in pixels. The format is `WIDTHxHEIGHT` where `WIDTH` and `HEIGHT` are whole numbers. When resizing an image, you may specify only the width (such as `600x`) or only the height (such as `x400`) for proportional scaling. Specifying both width and height when resizing an image may result in non-proportional scaling. When cropping, fitting, or filling, you must provide both width and height such as `600x400`.

format
: The format of the resulting image. Valid options include `bmp`, `gif`, `jpeg`, `png`, `tiff`, or `webp`. This defaults to the format of the source image.

hint
: The encoding preset used when processing WebP images, equivalent to the `-preset` flag for the [`cwebp`][] CLI. Valid options include `drawing`, `icon`, `photo`, `picture`, or `text`. This defaults to the [`hint`][] parameter in your site configuration.

  Value|Example
  :--|:--
  `drawing`|Hand or line drawing with high-contrast details
  `icon`|Small colorful image
  `photo`|Outdoor photograph with natural lighting
  `picture`|Indoor photograph such as a portrait
  `text`|Image that is primarily text

quality
: The visual fidelity of the image, applicable to JPEG and WebP formats when using `lossy` compression. The format is `qQUALITY` where `QUALITY` is a whole number between `1` and `100`, inclusive. Lower numbers prioritize smaller file size, while higher numbers prioritize visual clarity. This defaults to the [`quality`][] parameter in your site configuration.

resampling filter
: The algorithm used to calculate new pixels when resizing, fitting, or filling an image. Common options include `box`, `lanczos`, `catmullRom`, `mitchellNetravali`, `linear`, or `nearestNeighbor`. This defaults to the [`resampleFilter`][] parameter in your site configuration.

  Filter|Description
  :--|:--
  `box`|Simple and fast averaging filter appropriate for downscaling
  `lanczos`|High-quality resampling filter for photographic images yielding sharp results
  `catmullRom`|Sharp cubic filter that is faster than the Lanczos filter while providing similar results
  `mitchellNetravali`|Cubic filter that produces smoother results with less ringing artifacts than CatmullRom
  `linear`|Bilinear resampling filter, produces smooth output, faster than cubic filters
  `nearestNeighbor`|Fastest resampling filter, no antialiasing

  Refer to the [source documentation][] for a complete list of available resampling filters. If you wish to improve image quality at the expense of performance, you may wish to experiment with the alternative filters.

rotation
: The number of whole degrees to rotate an image counter-clockwise. The format is `rDEGREES` where `DEGREES` is a whole number. Hugo performs rotation before any other transformations, so your [target dimensions](#dimensions) and any [anchor](#anchor) should refer to the image orientation after rotation. Use `r90`, `r180`, or `r270` for orthogonal rotations, or arbitrary angles such as `r45`. To rotate clockwise, use a negative number such as `r-45`. To automatically rotate an image based on its Exif orientation tag, use the [`images.AutoOrient`][] filter instead of manual rotation.

  Rotating by non-orthogonal values increases the image extents to fit the rotated corners. For formats supporting alpha channels such as PNG or WebP, this resulting empty space is transparent by default. If the target format does not support transparency such as JPEG, or if you explicitly specify a [background color](#background-color) in the processing specification, the space is filled. If a color is required but not specified in the processing string, it defaults to the [`bgColor`][] parameter in your site configuration.

[`anchor`]: /configuration/imaging/#anchor
[`bgcolor`]: /configuration/imaging/#bgcolor
[`compression`]: /configuration/imaging/#compression
[`cwebp`]: https://developers.google.com/speed/webp/docs/cwebp
[`hint`]: /configuration/imaging/#hint
[`images.AutoOrient`]: /functions/images/autoorient/
[`images.Process`]: /functions/images/process/
[`Process`]: /methods/resource/process
[`quality`]: /configuration/imaging/#quality
[`resampleFilter`]: /configuration/imaging/#resamplefilter
[`smartcrop.js`]: https://github.com/jwagner/smartcrop.js
[hexadecimal color]: https://developer.mozilla.org/en-US/docs/Web/CSS/hex-color
[source documentation]: https://github.com/disintegration/imaging#image-resizing
