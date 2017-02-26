---
title: imageconfig
linktitle: imageConfig
description: Parses the image and returns the height, width, and color model.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [images]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

{{% warning %}}
`imageConfig` does not currently work in Hugo. See [the related `imageConfig` issue](https://github.com/spf13/hugo/issues/2806).
{{% /warning %}}

`imageConfig` parses the image and returns the height, width, and color model.

```golang
{{ with (imageConfig "favicon.ico") }}
favicon.ico: {{.Width}} x {{.Height}}
{{ end }}
```
