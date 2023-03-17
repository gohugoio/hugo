---
title: Sass / SCSS
description: Hugo Pipes allows the processing of Sass and SCSS files.
date: 2018-07-14
publishdate: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 30
weight: 02
sections_weight: 02
---

Any Sass or SCSS file can be transformed into a CSS file using `resources.ToCSS` which takes two arguments, the resource object and a map of options listed below.

```go-html-template
{{ $sass := resources.Get "sass/main.scss" }}
{{ $style := $sass | resources.ToCSS }}
```

### Options

transpiler [string]

: The `transpiler` to use, valid values are `libsass` (default) and `dartsass`. If you want to use Hugo with Dart Sass you need to download a release binary from [Embedded Dart Sass](https://github.com/sass/dart-sass-embedded/releases) and make sure it's in your PC's `$PATH` (or `%PATH%` on Windows).

targetPath [string]
: If not set, the transformed resource's target path will be the asset file original path with its extension replaced by `.css`.

vars [map]
: Map of key/value pairs that will be available in the `hugo:vars` namespace, e.g. with `@use "hugo:vars" as v;` or (globally) with `@import "hugo:vars";` {{< new-in "0.109.0" >}}

outputStyle [string]
: Default is `nested` (LibSass) and `expanded` (Dart Sass). Other available output styles for LibSass are `expanded`, `compact` and `compressed`. Dart Sass only supports `expanded` and `compressed`.

precision [int]
: Precision of floating point math. **Note:** This option is not supported by Dart Sass.

enableSourceMap [bool]
: When enabled, a source map will be generated.

sourceMapIncludeSources [bool]
: When enabled, sources will be embedded in the generated source map. (Dart Sass only). {{< new-in "0.108.0" >}}

includePaths [string slice]
: Additional SCSS/Sass include paths. Paths must be relative to the project directory.

```go-html-template
{{ $options := (dict "targetPath" "style.css" "outputStyle" "compressed" "enableSourceMap" (not hugo.IsProduction) "includePaths" (slice "node_modules/myscss")) }}
{{ $style := resources.Get "sass/main.scss" | resources.ToCSS $options }}
```

{{% note %}}
Setting `outputStyle` to `compressed` will handle Sass/SCSS files minification better than the more generic [`resources.Minify`]({{< ref "minification">}}).
{{% /note %}}
