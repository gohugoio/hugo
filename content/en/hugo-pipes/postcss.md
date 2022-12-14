---
title: PostCSS
description: Hugo Pipes can process CSS files with PostCSS.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 40
weight: 40
---

Any asset file can be processed using `resources.PostCSS` which takes for argument the resource object and a slice of options listed below.

The resource will be processed using the project's or theme's own `postcss.config.js` or any file set with the `config` option.

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $style := $css | resources.PostCSS }}
```

You must install the required Node.js packages to use the PostCSS feature. For example, to use the `autoprefixer` package, run these commands from the root of your project:

```text
npm install -D postcss postcss-cli autoprefixer
```

### Options

config [string]
: Set a custom directory to look for a config file

noMap [bool]
: Default is `false`. Disable the default inline sourcemaps

inlineImports [bool]
: Default is `false`. Enable inlining of @import statements. It does so recursively, but will only import a file once.
URL imports (e.g. `@import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');`) and imports with media queries will be ignored.
Note that this import routine does not care about the CSS spec, so you can have @import anywhere in the file.
Hugo will look for imports relative to the module mount and will respect theme overrides.

skipInlineImportsNotFound [bool] {{< new-in "0.99.0" >}}
: Default is `false`. Before Hugo 0.99.0 when `inlineImports` was enabled and we failed to resolve an import, we logged it as a warning. We now fail the build. If you have regular CSS imports in your CSS that you want to preserve, you can either use imports with URL or media queries (Hugo does not try to resolve those) or set `skipInlineImportsNotFound` to true.

_If no configuration file is used:_

use [string]
: Space-delimited list of PostCSS plugins to use

parser [string]
: Custom PostCSS parser

stringifier [string]
: Custom PostCSS stringifier

syntax [string]
: Custom postcss syntax

```go-html-template
{{ $options := dict "config" "/path/to/custom-config-directory" "noMap" true }}
{{ $style := resources.Get "css/main.css" | resources.PostCSS $options }}

{{ $options := dict "use" "autoprefixer postcss-color-alpha" }}
{{ $style := resources.Get "css/main.css" | resources.PostCSS $options }}
```

## Check Hugo Environment from postcss.config.js

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
