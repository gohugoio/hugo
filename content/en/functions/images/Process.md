---
title: images.Process
description: Returns an image filter that processes an image according to the given processing specification.
categories: []
keywords: [process]
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.Process SPECIFICATION]
---

Returns an image filter that processes an image according to the given [processing specification][]. This versatile filter supports the full range of image transformations, including resizing, cropping, rotation, and format conversion, all within a single specification string. Use this as an argument to the [`Filter`][] method or the [`images.Filter`][] function.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ $filter := images.Process "crop 200x200 TopRight webp q50" }}
  {{ with .Filter $filter }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"crop 200x200 TopRight webp q50"` is the _processing specification_.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Usage

Create a filter:

```go-html-template
{{ $filter := images.Process "crop 200x200 TopRight webp q50" }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

## Example

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="crop 200x200 TopRight webp q50"
  example=true
>}}

[`Filter`]: /methods/resource/filter/
[`images.Filter`]: /functions/images/filter
[processing specification]: #processing-specification
