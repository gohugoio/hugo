---
title: resources.Publish
description: Returns the given resource after publishing it.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.Publish RESOURCE]
---

{{< new-in 0.166.0 />}}

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `resources.Publish` function writes the given resource to the [`publishDir`][] and returns the resource, making it useful within a template pipeline.

```go-html-template
{{ resources.Get "main.js" | js.Build | resources.Publish }}
```

This is equivalent to:

```go-html-template
{{ (resources.Get "main.js" | js.Build).Publish }}
```

See the [`Publish`][] method for the non-pipeline form.

[`Publish`]: /methods/resource/publish/
[`publishDir`]: /configuration/all/#publishdir
