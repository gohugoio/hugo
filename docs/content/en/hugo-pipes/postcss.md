---
title: PostCSS
description: Hugo Pipes can process CSS files with PostCSS.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 40
weight: 40
sections_weight: 40
draft: false
---


Any asset file can be processed using `resources.PostCSS` which takes for argument the resource object and a slice of options listed below. 

The resource will be processed using the project's or theme's own `postcss.config.js` or any file set with the `config` option.


```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $style := $css | resources.PostCSS }}
```

{{% note %}}
Hugo Pipe's PostCSS requires the `postcss-cli` JavaScript package to be installed in the environment (`npm install -g postcss-cli`) along with any PostCSS plugin(s) used (e.g., `npm install -g autoprefixer`).

If you are using the Hugo Snap package, PostCSS and plugin(s) need to be installed locally within your Hugo site directory, e.g., `npm install postcss-cli` without the `-g` flag.
{{% /note %}}
### Options

config [string]
: Path to the PostCSS configuration file

noMap [bool]
: Default is `true`. Disable the default inline sourcemaps

_If no configuration file is used:_

use [string]
: List of PostCSS plugins to use

parser [string]
: Custom PostCSS parser

stringifier [string]
: Custom PostCSS stringifier

syntax [string]
: Custom postcss syntax

```go-html-template
{{ $style := resources.Get "css/main.css" | resources.PostCSS (dict "config" "customPostCSS.js" "noMap" true) }}
```