---
title: resources.Match
description: Returns a collection of global resources from paths matching the given glob pattern, or nil if none found.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/ByType
    - functions/resources/Get
    - functions/resources/GetMatch
    - functions/resources/GetRemote
    - methods/page/Resources
  returnType: resource.Resources
  signatures: [resources.Match PATTERN]
---

```go-html-template
{{ range resources.Match "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

{{% note %}}
This function operates on global resources. A global resource is a file within the assets directory, or within any directory mounted to the assets directory.

For page resources, use the [`Resources.Match`] method on the Page object.

[`Resources.Match`]: /methods/page/resources
{{% /note %}}

Hugo determines a match using a case-insensitive [glob pattern].

{{% include "functions/_common/glob-patterns.md" %}}

[glob pattern]: https://github.com/gobwas/glob#example
