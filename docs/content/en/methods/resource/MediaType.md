---
title: MediaType
description: Returns a media type object for the given resource.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: media.Type
    signatures: [RESOURCE.MediaType]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .MediaType.Type }} → image/jpeg
  {{ .MediaType.MainType }} → image
  {{ .MediaType.SubType }} → jpeg
  {{ .MediaType.Suffixes }} → [jpg jpeg jpe jif jfif]
  {{ .MediaType.FirstSuffix.Suffix }} → jpg
{{ end }}
```

## Methods

Use these methods on the `MediaType` object.

{{% include "/_common/methods/media-type/core-methods.md" %}}
