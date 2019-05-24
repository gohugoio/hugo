---
title: Hugo 0.32 HOWTO
description: About page bundles, image processing and more.
date: 2017-12-28
keywords: [ssg,static,performance,security]
menu:
  docs:
    parent: "about"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: []
toc: true
images:
- images/blog/sunset.jpg
---


{{% note %}}
This documentation belongs in other places in this documentation site, but is put here first ... to get something up and running fast.
{{% /note %}}


Also see this demo project from [bep](https://github.com/bep/), the clever Norwegian behind these new features:

* https://temp.bep.is/hugotest/
* https://github.com/bep/hugotest (source)

## Page Resources

### Organize Your Content

{{< figure src="/images/hugo-content-bundles.png" title="Pages with image resources" >}}

The content folder above shows a mix of content pages (`md` (i.e. markdown) files) and image resources.

{{% note %}}
You can use any file type as a content resource as long as it is a MIME type recognized by Hugo (`json` files will, as one example, work fine). If you want to get exotic, you can define your [own media type](/templates/output-formats/#media-types).
{{% /note %}}

The 3 page bundles marked in red explained from top to bottom:

1. The home page with one image resource (`1-logo.png`)
2. The blog section with two images resources and two pages resources (`content1.md`, `content2.md`). Note that the `_index.md` represents the URL for this section.
3. An article (`hugo-is-cool`) with a folder with some images and one content resource (`cats-info.md`). Note that the `index.md` represents the URL for this article.

The content files below `blog/posts` are just regular standalone pages.

{{% note %}}
Note that changes to any resource inside the `content` folder will trigger a reload when running in watch (aka server or live reload mode), it will even work with `--navigateToChanged`.
{{% /note %}}

#### Sort Order

* Pages are sorted according to standard Hugo page sorting rules.
* Images and other resources are sorted in lexicographical order.

### Handle Page Resources in Templates


#### List all Resources

```go-html-template
{{ range .Resources }}
<li><a href="{{ .RelPermalink }}">{{ .ResourceType | title }}</a></li>
{{ end }}
```

For an absolute URL, use `.Permalink`.

**Note:** The permalink will be relative to the content page, respecting permalink settings. Also, included page resources will not have a value for `RelPermalink`.

#### List All Resources by Type

```go-html-template
{{ with .Resources.ByType "image" }}
{{ end }}

```

Type here is `page` for pages, else the main type in the MIME type, so `image`, `json` etc.

#### Get a Specific Resource

```go-html-template
{{ $logo := .Resources.GetByPrefix "logo" }}
{{ with $logo }}
{{ end }}
```

#### Include Page Resource Content

```go-html-template
{{ with .Resources.ByType "page" }}
{{ range . }}
<h3>{{ .Title }}</h3>
{{ .Content }}
{{ end }}
{{ end }}

```


## Image Processing

The `image` resource implements the methods `Resize`, `Fit` and `Fill`:

Resize
: Resize to the given dimension, `{{ $logo.Resize "200x" }}` will resize to 200 pixels wide and preserve the aspect ratio. Use `{{ $logo.Resize "200x100" }}` to control both height and width.

Fit
: Scale down the image to fit the given dimensions, e.g. `{{ $logo.Fit "200x100" }}` will fit the image inside a box that is 200 pixels wide and 100 pixels high.

Fill
: Resize and crop the image given dimensions, e.g. `{{ $logo.Fill "200x100" }}` will resize and crop to width 200 and height 100


{{% note %}}
Image operations in Hugo currently **do not preserve EXIF data** as this is not supported by Go's [image package](https://github.com/golang/go/search?q=exif&type=Issues&utf8=%E2%9C%93). This will be improved on in the future.
{{% /note %}}


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

```go-html-template
{{</* imgproc sunset Resize "300x" */>}}
```

### Image Processing Options

In addition to the dimensions (e.g. `200x100`) where either height or width can be omitted, Hugo supports a set of additional image options:

Anchor
: Only relevant for `Fill`. This is useful for thumbnail generation where the main motive is located in, say, the left corner. Valid are `Center`, `TopLeft`, `Top`, `TopRight`, `Left`, `Right`, `BottomLeft`, `Bottom`, `BottomRight`. Example: `{{ $logo.Fill "200x100 BottomLeft" }}`

JPEG Quality
: Only relevant for JPEG images, values 1 to 100 inclusive, higher is better. Default is 75. `{{ $logo.Resize "200x q50" }}`

Rotate
: Rotates an image by the given angle counter-clockwise. The rotation will be performed first to get the dimensions correct. `{{ $logo.Resize "200x r90" }}`. The main use of this is to be able to manually correct for [EXIF orientation](https://github.com/golang/go/issues/4341) of JPEG images.

Resample Filter
: Filter used in resizing. Default is `Box`, a simple and fast resampling filter appropriate for downscaling. See https://github.com/disintegration/imaging for more. If you want to trade quality for faster processing, this may be a option to test. 



### Performance

Processed images are stored below `<project-dir>/resources` (can be set with `resourceDir` config setting). This folder is deliberately placed in the project, as it is recommended to check these into source control as part of the project. These images are not "Hugo fast" to generate, but once generated they can be reused.

If you change your image settings (e.g. size), remove or rename images etc., you will end up with unused images taking up space and cluttering your project. 

To clean up, run:

```bash
hugo --gc
```


{{% note %}}
**GC** is short for **Garbage Collection**.
{{% /note %}}


## Configuration

### Default Image Processing Config

You can configure an `imaging` section in `config.toml` with default image processing options:

```toml
[imaging]
# Default resample filter used for resizing. Default is Box,
# a simple and fast averaging filter appropriate for downscaling.
# See https://github.com/disintegration/imaging
resampleFilter = "box"

# Default JPEG quality setting. Default is 75.
quality = 68
```





