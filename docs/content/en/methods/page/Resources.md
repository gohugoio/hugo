---
title: Resources
description: Returns a collection of page resources.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: resource.Resources
    signatures: [PAGE.Resources]
---

The `Resources` method on a `Page` object returns a collection of page resources. A page resource is a file within a [page bundle](g).

To work with global or remote resources, see the [`resources`][] functions.

## Methods

Use these methods on the `Resources` object.

`ByType`
: (`resource.Resources`) Returns a collection of page resources of the given [media type](g), or `nil` if none found. The media type is typically one of `image`, `text`, `audio`, `video`, or `application`.

  ```go-html-template
  {{ range .Resources.ByType "image" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
  ```

  When working with global resources instead of page resources, use the [`resources.ByType`][] function.

`Get`
: (`resource.Resource`) Returns a page resource from the given path, or `nil` if none found.

  ```go-html-template
  {{ with .Resources.Get "images/a.jpg" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
  ```

  When working with global resources instead of page resources, use the [`resources.Get`][] function.

`GetMatch`
: (`resource.Resource`) Returns the first page resource from paths matching the given [glob pattern](g), or `nil` if none found.

  ```go-html-template
  {{ with .Resources.GetMatch "images/*.jpg" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
  ```

  When working with global resources instead of page resources, use the [`resources.GetMatch`][] function.

`Match`
: (`resource.Resources`) Returns a collection of page resources from paths matching the given [glob pattern](g), or `nil` if none found.

  ```go-html-template
  {{ range .Resources.Match "images/*.jpg" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
  ```

  When working with global resources instead of page resources, use the [`resources.Match`][] function.

`Mount`
: {{< new-in 0.140.0 />}}
: (`resource.ResourceGetter`) Mounts the given resources, remapping from the base path (first argument) to the target path (second argument), and returns a [resource getter](g). A leading slash in the target marks an absolute path. Relative target paths allow you to mount resources relative to another set, such as a [page bundle](g):

  ```go-html-template
  {{ $common := resources.Match "/js/headlessui/*.*" }}
  {{ $importContext := (slice $.Page ($common.Mount "/js/headlessui" ".")) }}
  ```

## Pattern matching

With the `GetMatch` and `Match` methods, Hugo determines a match using a case-insensitive [glob pattern](g). For syntax rules and examples, see the [glob patterns quick reference guide][].

[`resources.ByType`]: /functions/resources/bytype/
[`resources.GetMatch`]: /functions/resources/getmatch/
[`resources.Get`]: /functions/resources/get/
[`resources.Match`]: /functions/resources/match/
[`resources`]: /functions/resources/
[glob patterns quick reference guide]: /quick-reference/glob-patterns/
