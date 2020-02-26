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

inlineImports [bool] {{< new-in "0.66.0" >}}
: Default is `false`. Enable inlining of @import statements. It does so recursively, but will only import a file once.
URL imports (e.g. `@import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');`) and imports with media queries will be ignored.
Note that this import routine does not care about the CSS spec, so you can have @import anywhere in the file.
Hugo will look for imports relative to the module mount and will respect theme overrides.

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

## Check Hugo Environment from postcss.config.js

{{< new-in "0.66.0" >}}

The current Hugo environment name (set by `--environment` or in config or OS environment) is available in the Node context, which allows constructs like this:

```js
module.exports = {
  plugins: [
    require('autoprefixer'),
    ...process.env.HUGO_ENVIRONMENT === 'production'
      ? [purgecss]
      : []
  ]
}
```