---
title: resources.Match
description: Returns a collection of global resources from paths matching the given glob pattern, or nil if none found.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resources
    signatures: [resources.Match PATTERN]
---

```go-html-template
{{ range resources.Match "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

> [!note]
> This function operates on global resources. A global resource is a file within the `assets` directory, or within any directory mounted to the `assets` directory.
>
> For page resources, use the [`Resources.Match`] method on a `Page` object.

Hugo determines a match using a case-insensitive [glob pattern].

{{% include "/_common/glob-patterns.md" %}}

[`Resources.Match`]: /methods/page/resources/
[glob pattern]: https://github.com/gobwas/glob#example
