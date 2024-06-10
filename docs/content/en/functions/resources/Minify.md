---
title: resources.Minify
description: Minifies the given resource.
categories: []
keywords: []
action:
  aliases: [minify]
  related:
    - functions/js/Build
    - functions/resources/Babel
    - functions/resources/Fingerprint
    - functions/resources/PostCSS
    - functions/resources/ToCSS
  returnType: resource.Resource
  signatures: [resources.Minify RESOURCE]
---

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $style := $css | minify }}
```

Any CSS, JS, JSON, HTML, SVG, or XML resource can be minified using resources.Minify which takes for argument the resource object.
