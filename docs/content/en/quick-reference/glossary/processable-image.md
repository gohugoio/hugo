---
title: processable image
---

A _processable image_ is an image file characterized by one of the following [_media types_](g):

  - `image/bmp`
  - `image/gif`
  - `image/jpeg`
  - `image/png`
  - `image/tiff`
  - `image/webp`

  Hugo can decode and encode these image formats, allowing you to use any of the [resource methods][] applicable to images such as `Width`, `Height`, `Crop`, `Fill`, `Fit`, `Filter`, `Process`, `Resize`, etc.

  Use the [`reflect.IsImageResourceProcessable`][] function to determine if an image can be processed.

  [`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
  [resource methods]: /methods/resource
