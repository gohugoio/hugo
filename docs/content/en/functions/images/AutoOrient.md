---
title: images.AutoOrient
description: Returns an image filter that rotates and flips an image as needed per its EXIF orientation tag.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/images/Filter
    - methods/resource/Filter
  returnType: images.filter
  signatures: [images.AutoOrient]
toc: true
---

{{< new-in 0.121.2 >}}

## Usage

Create the filter:

```go-html-template
{{ $filter := images.AutoOrient }}
```

{{% include "functions/images/_common/apply-image-filter.md" %}}

{{% note %}}
When using with other filters, specify `images.AutoOrient` first.
{{% /note %}}

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
