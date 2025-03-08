---
title: resources.GetMatch
description: Returns the first global resource from paths matching the given glob pattern, or nil if none found.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.GetMatch PATTERN]
---

```go-html-template
{{ with resources.GetMatch "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

> [!note]
> This function operates on global resources. A global resource is a file within the `assets` directory, or within any directory mounted to the `assets` directory.
>
> For page resources, use the [`Resources.GetMatch`] method on a `Page` object.

Hugo determines a match using a case-insensitive [glob](g) pattern.

{{% include "/_common/glob-patterns.md" %}}

[`Resources.GetMatch`]: /methods/page/resources/
