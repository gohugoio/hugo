---
title: resources.Concat
description: Concatenates a slice of resources.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: resource.Resource
  signatures: ['resources.Concat TARGETPATH [RESOURCE...]']
---

```go-html-template
{{ $plugins := resources.Get "js/plugins.js" }}
{{ $global := resources.Get "js/global.js" }}
{{ $js := slice $plugins $global | resources.Concat "js/bundle.js" }}
```

Asset files of the same [media type] can be bundled into one resource using the `resources.Concat` function which takes two arguments, the target path for the created resource bundle and a slice of resource objects to be concatenated.

[media type]: https://en.wikipedia.org/wiki/Media_type
