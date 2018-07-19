---
title: Hugo Pipes Introduction
description: Hugo Pipes is Hugo's asset processing set of functions.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 20
weight: 01
sections_weight: 01
draft: false
aliases: [/assets/]
--- 

### Asset directory

Asset files must be stored in the asset directory. This is `/assets` by default, but can be configured via the configuration file's `assetDir` key.

### From file to resource

In order to process an asset with Hugo Pipes, it must be retrieved as a resource using `resources.Get`, which takes one argument: the filepath of the file relative to the asset directory.

```go-html-template
{{ $style := resources.Get "sass/main.scss" }}
```

### Asset publishing

Assets will only be published (to `/public`) if `.Permalink` or `.RelPermalink` is used.

### Go Pipes

For improved readability, the Hugo Pipes examples of this documentation will be written using [Go Pipes](/templates/introduction/#pipes):
```go-html-template
{{ $style := resources.Get "sass/main.scss" | resources.ToCSS | resources.Minify | resources.Fingerprint }}
<link rel="stylesheet" href="{{ $style.Permalink }}">
```

### Method aliases

Each Hugo Pipes `resources` transformation method uses a __camelCased__ alias (`toCSS` for `resources.ToCSS`).
Non-transformation methods deprived of such aliases are `resources.Get`, `resources.FromString`, `resources.ExecuteAsTemplate` and `resources.Concat`.

The example above can therefore also be written as follows:
```go-html-template
{{ $style := resources.Get "sass/main.scss" | toCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $style.Permalink }}">
```
