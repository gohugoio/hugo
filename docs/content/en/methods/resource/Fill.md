---
title: Fill
description: Applicable to images, returns a new image resource cropped and resized according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Fill SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Fill` method returns a new resource from a [processable image](g) according to the given [processing specification][].

> [!note]
> Use the [`reflect.IsImageResourceProcessable`][] function to verify that an image can be processed.

## Usage

When filling, you must provide both width and height (such as `500x200`) within the specification. `Fill` maintains the original aspect ratio by resizing the image to cover the target area and cropping any overflowing pixels based on the [anchor](#anchor) provided.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "500x200 TopRight" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"500x200 TopRight"` is the _processing specification.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fill "500x200 TopRight" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fill 500x200 TopRight"
  example=true
>}}

[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
[processing specification]: #processing-specification
