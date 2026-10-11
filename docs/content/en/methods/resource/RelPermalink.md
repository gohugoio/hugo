---
title: RelPermalink
description: Returns the relative permalink of the given resource, publishing it in the process.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [RESOURCE.RelPermalink]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `RelPermalink` method on a `Resource` object writes the resource to the `public` directory and returns its [relative permalink](g).

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .RelPermalink }} → /images/a.jpg
{{ end }}
```
