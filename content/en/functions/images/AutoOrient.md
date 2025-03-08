---
title: images.AutoOrient
description: Returns an image filter that rotates and flips an image as needed per its EXIF orientation tag.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: images.filter
    signatures: [images.AutoOrient]
---

{{< new-in 0.121.2 />}}

## Usage

Create the filter:

```go-html-template
{{ $filter := images.AutoOrient }}
```

{{% include "/_common/functions/images/apply-image-filter.md" %}}

> [!note]
> When using with other filters, specify `images.AutoOrient` first.

```go-html-template
{{ $filters := slice
  images.AutoOrient
  (images.Process "resize 200x")
}}
{{ with resources.Get "images/original.jpg" }}
  {{ with images.Filter $filters . }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

## Example

{{< img
  src="images/examples/landscape-exif-orientation-5.jpg"
  alt="Zion National Park"
  filter="AutoOrient"
  filterArgs=""
  example=true
>}}
