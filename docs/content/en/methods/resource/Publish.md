---
title: Publish
description: Publishes the given resource.
categories: []
keywords: []
action:
  related:
    - methods/resource/Permalink
    - methods/resource/RelPermalink
  returnType: nil
  signatures: [RESOURCE.Publish]
---

The `Publish` method on a `Resource` object writes the resource to the publish directory, typically `public`.

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Publish }}
{{ end }}
```

The `Permalink` and `RelPermalink` methods also publish a resource. `Publish` is a convenience method for publishing without a return value. For example, this:

```go-html-template
{{ $resource.Publish }}
```

Instead of this:

```go-html-template
{{ $noop := $resource.Permalink }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
