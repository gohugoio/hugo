---
title: Image processing
description: Process, transform, and analyze images.
categories: []
keywords: []
---

Hugo provides methods to transform and analyze images during the build process. The results are cached to ensure subsequent builds remain fast.

## Resources

To process an image you must capture the file as a page resource, a global resource, or a remote resource.

### Page

{{% glossary-term "page resource" %}}

```text
content/
└── posts/
    └── post-1/           <-- page bundle
        ├── index.md
        └── sunset.jpg    <-- page resource
```

To capture an image as a page resource:

```go-html-template
{{ $image := .Resources.Get "sunset.jpg" }}
```

### Global

{{% glossary-term "global resource" %}}

```text
assets/
└── images/
    └── sunset.jpg    <-- global resource
```

To capture an image as a global resource:

```go-html-template
{{ $image := resources.Get "images/sunset.jpg" }}
```

### Remote

{{% glossary-term "remote resource" %}}

To capture an image as a remote resource:

```go-html-template
{{ $image := resources.GetRemote "https://gohugo.io/img/hugo-logo.png" }}
```

## Rendering

Once you have captured an image as a resource, render it in your templates using the [`Permalink`][], [`RelPermalink`][], [`Width`][], and [`Height`][] methods.

Example 1: Throw an error if the resource is not found.

```go-html-template
{{ $image := .Resources.GetMatch "sunset.jpg" }}
<img src="{{ $image.RelPermalink }}" width="{{ $image.Width }}" height="{{ $image.Height }}">
```

Example 2: Skip image rendering if the resource is not found.

```go-html-template
{{ $image := .Resources.GetMatch "sunset.jpg" }}
{{ with $image }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
{{ end }}
```

Example 3: A more concise way to skip image rendering if the resource is not found.

```go-html-template
{{ with .Resources.GetMatch "sunset.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
{{ end }}
```

Example 4: Skip rendering if there's problem accessing a remote resource.

```go-html-template
{{ $url := "https://gohugo.io/img/hugo-logo.png" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

## Processing

To transform an image, apply a processing method to the image resource. Hugo generates the processed image on demand, caches the result, and returns a new resource object.

```go-html-template
{{ with .Resources.Get "sunset.jpg" }}
  {{ with .Resize "400x" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
  {{ end }}
{{ end }}
```

> [!note]
> Metadata is not preserved during image transformation. Use the `Exif` or `Meta` methods with the _original_ image resource to extract metadata from JPEG, PNG, TIFF, and WebP images.

Each method serves a specific transformation or metadata requirement:

Method|Description
:--|:--
[`Colors`]|Returns a slice of the most dominant colors using a simple histogram method.
[`Crop`]|Returns a new image resource cropped according to the given processing specification.
[`Exif`]|Applicable to JPEG, PNG, TIFF, and WebP images, returns an object containing Exif metadata.
[`Fill`]|Returns a new image resource cropped and resized according to the given processing specification.
[`Filter`]|Applies one or more image filters to the given image resource.
[`Fit`]|Returns a new image resource downscaled to fit according to the given processing specification.
[`Meta`]|Applicable to JPEG, PNG, TIFF, and WebP images, returns an object containing Exif, IPTC, and XMP metadata.
[`Process`]|Returns a new image resource processed according to the given processing specification.
[`Resize`]|Returns a new image resource resized according to the given processing specification.
{class="!mt-0"}

Select a method from the table above for syntax and usage examples.

## Performance

### Caching

Hugo processes images on demand and returns a new resource object. To ensure subsequent builds remain fast, Hugo caches the results in the directory specified in the [file cache] section of your site configuration.

If you host your site with Netlify, include the following in your site configuration to persist the image cache between builds:

```toml
[caches]
  [caches.images]
    dir = ':cacheDir/images'
```

### Garbage collection

If you change image processing methods, or rename/remove images, the cache will eventually contain unused files. To remove them and reclaim disk space, run Hugo's garbage collection:

```text
hugo --gc
```

### Resource usage

The time and memory required to process an image increase with the image's dimensions. For example, a `4032x2268` image requires significantly more memory and processing time than a `1920x1080` image.

If your source images are much larger than the maximum size you intend to publish, consider scaling them down before the build to optimize performance.

## Configuration

See [configure imaging](/configuration/imaging).

[`Colors`]: /methods/resource/colors/
[`Crop`]: /methods/resource/crop/
[`Exif`]: /methods/resource/exif/
[`Fill`]: /methods/resource/fill/
[`Filter`]: /methods/resource/filter/
[`Fit`]: /methods/resource/fit/
[`Height`]: /methods/resource/height/
[`Meta`]: /methods/resource/meta/
[`Permalink`]: /methods/resource/permalink/
[`Process`]: /methods/resource/process/
[`RelPermalink`]: /methods/resource/relpermalink/
[`Resize`]: /methods/resource/resize/
[`Width`]: /methods/resource/width/
