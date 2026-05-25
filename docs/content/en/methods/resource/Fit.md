---
title: Fit
description: Applicable to images, returns a new image resource downscaled to fit according to the given processing specification.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: images.ImageResource
    signatures: [RESOURCE.Fit SPECIFICATION]
---

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Fit` method returns a new resource from a [processable image](g) according to the given [processing specification][].

> [!note]
> Use the [`reflect.IsImageResourceProcessable`][] function to verify that an image can be processed.

## Usage

When fitting, you must provide both width and height (such as `300x175`) within the specification. `Fit` maintains the original aspect ratio by downscaling the image until it fits within the specified dimensions. Unlike [`Fill`][] or [`Resize`][], this method will never upscale an image; if the source image is smaller than the target dimensions, the dimensions of the resulting image are the same as the original.

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "300x175" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

In the example above, `"300x175"` is the processing specification.

{{% include "/_common/methods/resource/processing-spec.md" %}}

## Example

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Fit "300x175" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

{{< img
  src="images/examples/zion-national-park.jpg"
  alt="Zion National Park"
  filter="Process"
  filterArgs="fit 300x175"
  example=true
>}}

[`Fill`]: /methods/resource/fill/
[`Resize`]: /methods/resource/resize/
[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
[processing specification]: #processing-specification
