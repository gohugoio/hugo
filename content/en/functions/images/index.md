---
title: Image Functions
description: The images namespace provides a list of filters and other image related functions.
godocref:
date: 2017-02-01
categories: [functions]
aliases: [/functions/imageconfig/]
menu:
  docs:
    parent: "functions"
keywords: [images]
toc: true
---


## Image Filters

See [images.Filter](#filter) for how to apply these filters to an image.

### Brightness

{{% funcsig %}}
images.Brightness PERCENTAGE
{{% /funcsig %}}

Brightness creates a filter that changes the brightness of an image.
The percentage parameter must be in range (-100, 100).

### ColorBalance

{{% funcsig %}}
images.ColorBalance PERCENTAGERED PERCENTAGEGREEN PERCENTAGEBLUE
{{% /funcsig %}}

ColorBalance creates a filter that changes the color balance of an image.
The percentage parameters for each color channel (red, green, blue) must be in range (-100, 500).

### Colorize

{{% funcsig %}}
images.Colorize HUE SATURATION PERCENTAGE
{{% /funcsig %}}

Colorize creates a filter that produces a colorized version of an image.
The hue parameter is the angle on the color wheel, typically in range (0, 360).
The saturation parameter must be in range (0, 100).
The percentage parameter specifies the strength of the effect, it must be in range (0, 100).

### Contrast

{{% funcsig %}}
images.Contrast PERCENTAGE
{{% /funcsig %}}

Contrast creates a filter that changes the contrast of an image.
The percentage parameter must be in range (-100, 100).

### Gamma

{{% funcsig %}}
images.Gamma GAMMA
{{% /funcsig %}}

Gamma creates a filter that performs a gamma correction on an image.
The gamma parameter must be positive. Gamma = 1 gives the original image.
Gamma less than 1 darkens the image and gamma greater than 1 lightens it.

### GaussianBlur

{{% funcsig %}}
images.GaussianBlur SIGMA
{{% /funcsig %}}

GaussianBlur creates a filter that applies a gaussian blur to an image.

### Grayscale

{{% funcsig %}}
images.Grayscale
{{% /funcsig %}}

Grayscale creates a filter that produces a grayscale version of an image.

### Hue

{{% funcsig %}}
images.Hue SHIFT
{{% /funcsig %}}

Hue creates a filter that rotates the hue of an image.
The hue angle shift is typically in range -180 to 180.

### Invert

{{% funcsig %}}
images.Invert
{{% /funcsig %}}

Invert creates a filter that negates the colors of an image.

### Pixelate

{{% funcsig %}}
images.Pixelate SIZE
{{% /funcsig %}}

Pixelate creates a filter that applies a pixelation effect to an image.

### Saturation

{{% funcsig %}}
images.Saturation PERCENTAGE
{{% /funcsig %}}

Saturation creates a filter that changes the saturation of an image.

### Sepia

{{% funcsig %}}
images.Sepia PERCENTAGE
{{% /funcsig %}}

Sepia creates a filter that produces a sepia-toned version of an image.

### Sigmoid

{{% funcsig %}}
images.Sigmoid MIDPOINT FACTOR
{{% /funcsig %}}

Sigmoid creates a filter that changes the contrast of an image using a sigmoidal function and returns the adjusted image.
It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.

### UnsharpMask

{{% funcsig %}}
images.UnsharpMask SIGMA AMOUNT THRESHOLD
{{% /funcsig %}}

UnsharpMask creates a filter that sharpens an image.
The sigma parameter is used in a gaussian function and affects the radius of effect.
Sigma must be positive. Sharpen radius roughly equals 3 * sigma.
The amount parameter controls how much darker and how much lighter the edge borders become. Typically between 0.5 and 1.5.
The threshold parameter controls the minimum brightness change that will be sharpened. Typically between 0 and 0.05.

## Other Functions

### Filter

{{% funcsig %}}
IMAGE | images.Filter FILTERS...
{{% /funcsig %}}

Can be used to apply a set of filters to an image:

```go-html-template
{{ $img := $img | images.Filter (images.GaussianBlur 6) (images.Pixelate 8) }}
```

Also see the [Filter Method](/content-management/image-processing/#filter).

### ImageConfig

Parses the image and returns the height, width, and color model.

{{% funcsig %}}
images.ImageConfig PATH
{{% /funcsig %}}

```go-html-template
{{ with (imageConfig "favicon.ico") }}
favicon.ico: {{.Width}} x {{.Height}}
{{ end }}
```
