---
title: images.Overlay
description: Returns an image filter that overlays the source image at the given coordinates, relative to the upper left corner.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Overlay RESOURCE X Y]
---

## Usage

Capture the overlay image as a resource:

```go-html-template
{{ $overlay := "" }}
{{ $path := "images/logo.png" }}
{{ with resources.Get $path }}
  {{ $overlay = . }}
{{ else }}
  {{ errorf "Unable to get resource %q" $path }}
{{ end }}
```

The overlay image can be a [global resource](g), a [page resource](g), or a [remote resource](g).

Create the filter:

```go-html-template
{{ $filter := images.Overlay $overlay 20 20 }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Overlay"
  filterArgs="images/logos/logo-64x64.png,20,20"
  example=true
>}}
