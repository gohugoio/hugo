---
title: "Image Processing"
description: "Image Page resources can be resized and cropped."
date: 2018-01-24T13:10:00-05:00
lastmod: 2018-01-26T15:59:07-05:00
linktitle: "Image Processing"
categories: ["content management"]
keywords: [bundle,content,resources,images]
weight: 4004
draft: false
toc: true
menu:
  docs:
    parent: "content-management"
    weight: 32
---

## The Image Page Resource

The `image` is a [Page Resource]({{< relref "content-management/page-resources" >}}), and the processing methods listed below does not work on images inside your `/static` folder.


To get all images in a [Page Bundle]({{< relref "content-management/organization#page-bundles" >}}):


```html
{{ with .Resources.ByType "image" }}
{{ end }}

```

## Image Processing Methods


The `image` resource implements the methods `Resize`, `Fit` and `Fill`, each returning the transformed image using the specified dimensions and processing options.

Resize
: Resizes the image to the specified width and height.

```go
// Resize to a width of 600px and preserve ratio
{{ $image := $resource.Resize "600x" }} 

// Resize to a height of 400px and preserve ratio
{{ $image := $resource.Resize "x400" }} 

// Resize to a width 600px and a height of 400px
{{ $image := $resource.Resize "600x400" }}
```

Fit
: Scale down the image to fit the given dimensions while maintaining aspect ratio. Both height and width are required.

```go
{{ $image := $resource.Fit "600x400" }} 
```

Fill
: Resize and crop the image to match the given dimensions. Both height and width are required.

```go
{{ $image := $resource.Fill "600x400" }} 
```


{{% note %}}
Image operations in Hugo currently **do not preserve EXIF data** as this is not supported by Go's [image package](https://github.com/golang/go/search?q=exif&type=Issues&utf8=%E2%9C%93). This will be improved on in the future.
{{% /note %}}


## Image Processing Options

In addition to the dimensions (e.g. `600x400`), Hugo supports a set of additional image options.


JPEG Quality
: Only relevant for JPEG images, values 1 to 100 inclusive, higher is better. Default is 75.

```go
{{ $image.Resize "600x q50" }}
```

Rotate
: Rotates an image by the given angle counter-clockwise. The rotation will be performed first to get the dimensions correct. The main use of this is to be able to manually correct for [EXIF orientation](https://github.com/golang/go/issues/4341) of JPEG images.

```go
{{ $image.Resize "600x r90" }}
```

Anchor
: Only relevant for the `Fill` method. This is useful for thumbnail generation where the main motive is located in, say, the left corner. 
Valid are `Center`, `TopLeft`, `Top`, `TopRight`, `Left`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`.

```go
{{ $image.Fill "300x200 BottomLeft" }}
```

Resample Filter
: Filter used in resizing. Default is `Box`, a simple and fast resampling filter appropriate for downscaling. 

Examples are: `Box`, `NearestNeighbor`, `Linear`, `Gaussian`.

See https://github.com/disintegration/imaging for more. If you want to trade quality for faster processing, this may be a option to test. 

```go
{{ $image.Resize "600x400 Gaussian" }}
```

### Image Processing Examples

_The photo of the sunset used in the examples below is Copyright [Bjørn Erik Pedersen](https://commons.wikimedia.org/wiki/User:Bep) (Creative Commons Attribution-Share Alike 4.0 International license)_


{{< imgproc sunset Resize "300x" />}}

{{< imgproc sunset Fill "90x120 left" />}}

{{< imgproc sunset Fill "90x120 right" />}}

{{< imgproc sunset Fit "90x90" />}}

{{< imgproc sunset Resize "300x q10" />}}


This is the shortcode used in the examples above:


{{< code file="layouts/shortcodes/imgproc.html" >}}
{{< readfile file="layouts/shortcodes/imgproc.html" >}}   
{{< /code >}}

And it is used like this:

```html
{{</* imgproc sunset Resize "300x" /*/>}}
```


{{% note %}}
**Tip:** Note the self-closing shortcode syntax above. The `imgproc` shortcode can be called both with and without **inner content**.
{{% /note %}}