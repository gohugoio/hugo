---
title: resources.GetMatch
description: Returns the first global resource from paths matching the given glob pattern, or nil if none found.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/ByType
    - functions/resources/Get
    - functions/resources/GetRemote
    - functions/resources/Match
    - methods/page/Resources
  returnType: resource.Resource
  signatures: [resources.GetMatch PATTERN]
---

```go-html-template
{{ with resources.GetMatch "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{% note %}}
This function operates on global resources. A global resource is a file within the assets directory, or within any directory mounted to the assets directory.

For page resources, use the [`Resources.GetMatch`] method on the Page object.

[`Resources.GetMatch`]: /methods/page/resources/
{{% /note %}}

Hugo determines a match using a case-insensitive [glob pattern].

{{% include "functions/_common/glob-patterns.md" %}}

[glob pattern]: https://github.com/gobwas/glob#example
