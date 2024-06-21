---
title: Permalink
description:  Publishes the given resource and returns its permalink.
categories: []
keywords: []
action:
  related:
    - methods/resource/RelPermalink
    - methods/resource/Publish
  returnType: string
  signatures: [RESOURCE.Permalink]
---

The `Permalink` method on a `Resource` object writes the resource to the publish directory, typically `public`, and returns its [permalink].

[permalink]: /getting-started/glossary/#permalink

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Permalink }} → https://example.org/images/a.jpg
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
