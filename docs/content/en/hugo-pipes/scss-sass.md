---
title: SASS / SCSS
description: Hugo Pipes allows the processing of SASS and SCSS files.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 30
weight: 02
sections_weight: 02
draft: false
---


Any SASS or SCSS file can be transformed into a CSS file using `resources.ToCSS` which takes two arguments, the resource object and a map of options listed below.

```go-html-template
{{ $sass := resources.Get "sass/main.scss" }}
{{ $style := $sass | resources.ToCSS }}
```

### Options
targetPath [string]
: If not set, the resource's target path will be the asset file original path with its extension replaced by `.css`.

outputStyle [string]
: Default is `nested`. Other available output styles are `expanded`, `compact` and `compressed`.

precision [int]
: Precision of floating point math.

enableSourceMap [bool]
: When enabled, a source map will be generated.

includePaths [string slice]
: Additional SCSS/SASS include paths. Paths must be relative to the project directory.

```go-html-template
{{ $options := (dict "targetPath" "style.css" "outputStyle" "compressed" "enableSourceMap" true "includePaths" (slice "node_modules/myscss")) }}
{{ $style := resources.Get "sass/main.scss" | resources.ToCSS $options }}
```

{{% note %}}
Setting `outputStyle` to `compressed` will handle SASS/SCSS files minification better than the more generic [`resources.Minify`]({{< ref "minification">}}).
{{% /note %}}
