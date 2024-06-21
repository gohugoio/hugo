---
title: images.Config
description: Returns an image.Config structure from the image at the specified path, relative to the working directory.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: image.Config
  signatures: [images.Config PATH]
aliases: [/functions/imageconfig]
---

See [image processing] for an overview of Hugo's image pipeline.

[image processing]: /content-management/image-processing/

```go-html-template
{{ $ic := images.Config "/static/images/a.jpg" }}

{{ $ic.Width }} → 600 (int)
{{ $ic.Height }} → 400 (int)
```

Supported image formats include GIF, JPEG, PNG, TIFF, and WebP.

{{% note %}}
This is a legacy function, superseded by the [`Width`] and [`Height`] methods for [global], [page], and [remote] resources. See the [image processing] section for details.

[`Width`]: /methods/resource/width/
[`Height`]: /methods/resource/height/
[global]: /getting-started/glossary/#global-resource
[image processing]: /content-management/image-processing/
[page]: /getting-started/glossary/#page-resource
[remote]: /getting-started/glossary/#remote-resource
{{% /note %}}
