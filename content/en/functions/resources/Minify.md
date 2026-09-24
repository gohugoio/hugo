---
title: resources.Minify
description: Returns a minified version of the given resource.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [minify]
    returnType: resource.Resource
    signatures: [resources.Minify RESOURCE]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $style := $css | minify }}
```

Any CSS, JS, JSON, HTML, SVG, or XML resource can be minified using resources.Minify which takes for argument the resource object.
