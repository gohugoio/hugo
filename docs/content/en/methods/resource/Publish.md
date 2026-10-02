---
title: Publish
description: Publishes the given resource.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: nil
    signatures: [RESOURCE.Publish]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Publish` method on a `Resource` object writes the given resource to the [`publishDir`][].

This example uses [`resources.FromString`][] to create a resource from a string, then publishes it:

```go-html-template {file="layouts/baseof.html"}
{{ $content := printf "Contact: mailto:%s\n" site.Params.email }}
{{ with resources.FromString ".well-known/security.txt" $content }}
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

To publish a resource within a pipeline, use the [`resources.Publish`][] function instead.

[`publishDir`]: /configuration/all/#publishdir
[`resources.FromString`]: /functions/resources/fromstring/
[`resources.Publish`]: /functions/resources/publish/
