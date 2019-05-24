---
title: imageConfig
linktitle: imageConfig
description: Parses the image and returns the height, width, and color model.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [images]
signature: ["imageConfig PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

```
{{ with (imageConfig "favicon.ico") }}
favicon.ico: {{.Width}} x {{.Height}}
{{ end }}
```
