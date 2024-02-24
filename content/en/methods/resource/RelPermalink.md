---
title: RelPermalink
description: Publishes the given resource and returns its relative permalink.
categories: []
keywords: []
action:
  related:
    - methods/resource/Permalink
    - methods/resource/Publish
  returnType: string
  signatures: [RESOURCE.RelPermalink]
---

The `Permalink` method on a `Resource` object writes the resource to the publish directory, typically `public`, and returns its [relative permalink].

[relative permalink]: /getting-started/glossary/#relative-permalink

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .RelPermalink }} → /images/a.jpg
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
