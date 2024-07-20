---
title: Resources
description: Returns a collection of page resources.
categories: []
keywords: []
action:
  related:
    - functions/resources/ByType
    - functions/resources/Get
    - functions/resources/GetMatch
    - functions/resources/GetRemote
    - functions/resources/Match
  returnType: resource.Resources
  signatures: [PAGE.Resources]
toc: true
---

The `Resources` method on a `Page` object returns a collection of page resources. A page resource is a file within a [page bundle].

To work with global or remote resources, see the [`resources`] functions.

## Methods

###### ByType

(`resource.Resources`) Returns a collection of page resources of the given [media type], or nil if none found. The media type is typically one of `image`, `text`, `audio`, `video`, or `application`.

```go-html-template
{{ range .Resources.ByType "image" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

When working with global resources instead of page resources, use the [`resources.ByType`] function.

###### Get

(`resource.Resource`) Returns a page resource from the given path, or nil if none found.

```go-html-template
{{ with .Resources.Get "images/a.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

When working with global resources instead of page resources, use the [`resources.Get`] function.

###### GetMatch

(`resource.Resource`) Returns the first page resource from paths matching the given [glob pattern], or nil if none found.

```go-html-template
{{ with .Resources.GetMatch "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

When working with global resources instead of page resources, use the [`resources.GetMatch`] function.

###### Match

(`resource.Resources`) Returns a collection of page resources from paths matching the given [glob pattern], or nil if none found.

```go-html-template
{{ range .Resources.Match "images/*.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

When working with global resources instead of page resources, use the [`resources.Match`] function.

## Pattern matching

With the `GetMatch` and `Match` methods, Hugo determines a match using a case-insensitive [glob pattern].

{{% include "functions/_common/glob-patterns.md" %}}

[`resources.ByType`]: /functions/resources/ByType/
[`resources.GetMatch`]: /functions/resources/ByType/
[`resources.Get`]: /functions/resources/ByType/
[`resources.Match`]: /functions/resources/ByType/
[`resources`]: /functions/resources/
[glob pattern]: https://github.com/gobwas/glob#example
[media type]: https://en.wikipedia.org/wiki/Media_type
[page bundle]: /getting-started/glossary/#page-bundle
