---
title: images.Config
description: Returns an image.Config structure from the image at the specified path, relative to the working directory.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: image.Config
    signatures: [images.Config PATH]
aliases: [/functions/imageconfig]
---

> [!note]
> This is a legacy function, superseded by the [`Width`][] and [`Height`][] methods for [global resources](g), [page resources](g), and [remote resources](g). See the [image processing][] section for details.

```go-html-template
{{ $ic := images.Config "/static/images/a.jpg" }}

{{ $ic.Width }} → 600 (int)
{{ $ic.Height }} → 400 (int)
```

Supported image formats include AVIF, BMP, GIF, HEIC, HEIF, JPEG, PNG, TIFF, and WebP.

[`Height`]: /methods/resource/height/
[`Width`]: /methods/resource/width/
[image processing]: /content-management/image-processing/
