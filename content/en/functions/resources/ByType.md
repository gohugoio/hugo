---
title: resources.ByType
description: Returns a collection of global resources of the given media type, or nil if none found.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/Get
    - functions/resources/GetMatch
    - functions/resources/GetRemote
    - functions/resources/Match
    - methods/page/Resources
  returnType: resource.Resources
  signatures: [resources.ByType MEDIATYPE]
---

The [media type] is typically one of `image`, `text`, `audio`, `video`, or `application`.

```go-html-template
{{ range resources.ByType "image" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{% note %}}
This function operates on global resources. A global resource is a file within the assets directory, or within any directory mounted to the assets directory.

For page resources, use the [`Resources.ByType`] method on a `Page` object.

[`Resources.ByType`]: /methods/page/resources/
{{% /note %}}

[media type]: https://en.wikipedia.org/wiki/Media_type
