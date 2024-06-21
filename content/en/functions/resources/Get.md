---
title: resources.Get
description: Returns a global resource from the given path, or nil if none found.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/ByType
    - functions/resources/GetMatch
    - functions/resources/GetRemote
    - functions/resources/Match
    - methods/page/Resources
  returnType: resource.Resource
  signatures: [resources.Get PATH]
---

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{% note %}}
This function operates on global resources. A global resource is a file within the assets directory, or within any directory mounted to the assets directory.

For page resources, use the [`Resources.Get`] method on the Page object.

[`Resources.Get`]: /methods/page/resources/
{{% /note %}}
