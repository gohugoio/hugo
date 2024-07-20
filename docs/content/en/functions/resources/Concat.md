---
title: resources.Concat
description: Returns a concatenated slice of resources.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: resource.Resource
  signatures: ['resources.Concat TARGETPATH [RESOURCE...]']
---

The `resources.Concat` function returns a concatenated slice of resources, caching the result using the target path as its cache key. Each resource must have the same [media type].

Hugo publishes the resource to the target path when you call its [`Publish`], [`Permalink`], or [`RelPermalink`] methods. 

[media type]: https://en.wikipedia.org/wiki/Media_type
[`publish`]: /methods/resource/publish/
[`permalink`]: /methods/resource/permalink/
[`relpermalink`]: /methods/resource/relpermalink/

```go-html-template
{{ $plugins := resources.Get "js/plugins.js" }}
{{ $global := resources.Get "js/global.js" }}
{{ $js := slice $plugins $global | resources.Concat "js/bundle.js" }}
```
