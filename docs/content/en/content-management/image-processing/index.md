---
title: Image processing
description: Transform images to change their size, shape, and appearance.
categories: []
keywords: []
---

Hugo provides methods to transform and analyze images during the build process. While Hugo can manage any image format as a resource, only [processable images](g) can be transformed using the methods below. The results are cached to ensure subsequent builds remain fast.

> [!note]
> Use the [`reflect.IsImageResourceProcessable`][] function to verify that an image can be processed.

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

{{% include "/_common/functions/reflect/image-reflection-functions.md" %}}

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
> Metadata is not preserved during image transformation. Use the [`Meta`][] method with the original image resource to extract metadata from supported formats.

Select a method from the table below for syntax and usage examples, depending on your specific transformation or metadata requirements:

{{% render-table-of-pages-in-section
  path=/methods/resource
  filter=methods_resource_image_processing
  filterType=include
  headingColumn1=Method
  headingColumn2=Description
%}}{class="!mt-0"}

## Performance

### Caching

Hugo processes images on demand and returns a new resource object. To ensure subsequent builds remain fast, Hugo caches the results in the directory specified in the [file cache][] section of your project configuration.

If you host your site with Netlify, include the following in your project configuration to persist the image cache between builds:

```toml
[caches]
  [caches.images]
    dir = ':cacheDir/images'
```

### Garbage collection

If you change image processing methods, or rename/remove images, the cache will eventually contain unused files. To remove them and reclaim disk space, run Hugo's garbage collection:

```text
hugo build --gc
```

### Resource usage

The time and memory required to process an image increase with the image's dimensions. For example, a `4032x2268` image requires significantly more memory and processing time than a `1920x1080` image.

If your source images are much larger than the maximum size you intend to publish, consider scaling them down before the build to optimize performance.

## Configuration

See [configure imaging](/configuration/imaging).

[`Height`]: /methods/resource/height/
[`Meta`]: /methods/resource/meta/
[`Permalink`]: /methods/resource/permalink/
[`RelPermalink`]: /methods/resource/relpermalink/
[`Width`]: /methods/resource/width/
[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
[file cache]: /configuration/caches/
