---
title: Hugo Pipes
linkTitle: Introduction
description: Hugo Pipes is Hugo's asset processing set of functions.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    weight: 20
weight: 20
toc: true
aliases: [/assets/]
---

## Find resources in assets

This is about global and remote resources.

global resource
: A file within the assets directory, or within any directory [mounted] to the assets directory.

remote resource
: A file on a remote server, accessible via HTTP or HTTPS.

For `.Page` scoped resources, see the [page resources] section.

[mounted]: /hugo-modules/configuration/#module-configuration-mounts
[page resources]: /content-management/page-resources/

## Get a resource

In order to process an asset with Hugo Pipes, it must be retrieved as a resource.

For global resources, use:

- [`resources.ByType`](/functions/resources/bytype/)
- [`resources.Get`](/functions/resources/get/)
- [`resources.GetMatch`](/functions/resources/getmatch/)
- [`resources.Match`](/functions/resources/match/)

For remote resources, use:

- [`resources.GetRemote`](/functions/resources/getremote/)

See the [GoDoc Page](https://pkg.go.dev/github.com/gohugoio/hugo/tpl/resources) for the `resources` package for an up to date overview of all template functions in this namespace.

## Copy a resource

See the [`resources.Copy`](/functions/resources/copy/) function.

## Asset directory

Asset files must be stored in the asset directory. This is `/assets` by default, but can be configured via the configuration file's `assetDir` key.

## Asset publishing

Hugo publishes assets to the `publishDir` (typically `public`) when you invoke `.Permalink`, `.RelPermalink`, or `.Publish`. You can use `.Content` to inline the asset.

## Go Pipes

For improved readability, the Hugo Pipes examples of this documentation will be written using [Go Pipes](/templates/introduction/#pipes):

```go-html-template
{{ $style := resources.Get "sass/main.scss" | css.Sass | resources.Minify | resources.Fingerprint }}
<link rel="stylesheet" href="{{ $style.Permalink }}">
```

## Caching

Hugo Pipes invocations are cached based on the entire *pipe chain*.

An example of a pipe chain is:

```go-html-template
{{ $mainJs := resources.Get "js/main.js" | js.Build "main.js" | minify | fingerprint }}
```

The pipe chain is only invoked the first time it is encountered in a site build, and results are otherwise loaded from cache. As such, Hugo Pipes can be used in templates which are executed thousands or millions of times without negatively impacting the build performance.
