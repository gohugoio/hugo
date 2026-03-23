---
title: reflect.IsImageResource
description: Reports whether the given value is a Resource object representing an image as defined by its media type.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsImageResource INPUT]
---

{{< new-in 0.154.0 />}}

## Usage

This example iterates through all project resources and uses `reflect.IsImageResource` to decide whether to render an image tag or provide a download link for non-image files.

```go-html-template
{{ range resources.Match "**" }}
  {{ if reflect.IsImageResource . }}
    <img src="{{ .RelPermalink }}" alt="Image">
  {{ else }}
    <a href="{{ .RelPermalink }}">Download</a>
  {{ end }}
{{ end }}
```

{{% include "/_common/functions/reflect/image-reflection-functions.md" %}}
