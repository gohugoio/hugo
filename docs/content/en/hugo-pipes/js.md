---
title: JavaScript Building
description: Hugo Pipes can process JavaScript files with [ESBuild](https://github.com/evanw/esbuild).
date: 2020-07-20
publishdate: 2020-07-20
lastmod: 2020-07-20
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 45
weight: 45
sections_weight: 45
draft: false
---

Any JavaScript resource file can be transpiled and "tree shaken" using `js.Build` which takes for argument either a string for the filepath or a dict of options listed below.

### Options

targetPath [string]
: If not set, the source path will be used as the base target path. 
Note that the target path's extension may change if the target MIME type is different, e.g. when the source is TypeScript.

minify [bool]
: Let `js.Build` handle the minification.

target [string]
: The language target.
  One of: `es5`, `es2015`, `es2016`, `es2017`, `es2018`, `es2019`, `es2020` or `esnext`.
  Default is `esnext`.

externals [slice]
: External dependencies. If a dependency should not be included in the bundle (Ex. library loaded from a CDN.), it should be listed here.

```go-html-template
{{ $externals := slice "react" "react-dom" }}
```

defines [map]
: Allow to define a set of string replacement to be performed when building. Should be a map where each key is to be replaced by its value.

```go-html-template
{{ $defines := dict "process.env.NODE_ENV" `"development"` }}
```

format [string] {{< new-in "0.74.3" >}}
: The output format.
  One of: `iife`, `cjs`, `esm`.
  Default is `iife`, a self-executing function, suitable for inclusion as a <script> tag. 

### Examples

```go-html-template
{{ $built := resources.Get "js/index.js" | js.Build "main.js" }}
```

Or with options:

```go-html-template
{{ $externals := slice "react" "react-dom" }}
{{ $defines := dict "process.env.NODE_ENV" `"development"` }}

{{ $opts := dict "targetPath" "main.js" "externals" $externals "defines" $defines }}
{{ $built := resources.Get "scripts/main.js" | js.Build $opts }}
<script type="text/javascript" src="{{ $built.RelPermalink }}" defer></script>
```
