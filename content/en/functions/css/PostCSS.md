---
title: css.PostCSS
description: Processes the given resource with PostCSS using any PostCSS plugin.
categories: []
keywords: []
action:
  aliases: [postCSS]
  related:
    - functions/css/Sass
    - functions/css/TailwindCSS
  returnType: resource.Resource
  signatures: ['css.PostCSS [OPTIONS] RESOURCE']
toc: true
---

{{< new-in 0.128.0 >}}

```go-html-template
{{ with resources.Get "css/main.css" | postCSS }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

## Setup

Follow the steps below to transform CSS using any of the available [PostCSS plugins].

Step 1
: Install [Node.js].

Step 2
: Install the required Node.js packages in the root of your project. For example, to add vendor prefixes to your CSS rules:

```sh
npm i -D postcss postcss-cli autoprefixer
```

Step 3
: Create a PostCSS configuration file in the root of your project. You must name this file `postcss.config.js` or another [supported file name]. For example:

```js
module.exports = {
  plugins: [
    require('autoprefixer')
  ]
};
```

{{% note %}}
{{% include "functions/resources/_common/postcss-windows-warning.md" %}}
{{% /note %}}

Step 4
: Place your CSS file within the `assets/css` directory.

Step 5
: Process the resource with PostCSS:

```go-html-template
{{ with resources.Get "css/main.css" | postCSS }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

## Options

The `css.PostCSS` method takes an optional map of options.

config
: (`string`) The directory that contains the PostCSS configuration file. Default is the root of the project directory.

noMap
: (`bool`) Default is `false`. If `true`, disables inline sourcemaps.

inlineImports
: (`bool`) Default is `false`. Enable inlining of @import statements. It does so recursively, but will only import a file once. URL imports (e.g. `@import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');`) and imports with media queries will be ignored. Note that this import routine does not care about the CSS spec, so you can have @import anywhere in the file. Hugo will look for imports relative to the module mount and will respect theme overrides.

skipInlineImportsNotFound
: (`bool`) Default is `false`. Before Hugo 0.99.0 when `inlineImports` was enabled and we failed to resolve an import, we logged it as a warning. We now fail the build. If you have regular CSS imports in your CSS that you want to preserve, you can either use imports with URL or media queries (Hugo does not try to resolve those) or set `skipInlineImportsNotFound` to true.

```go-html-template
{{ $opts := dict "config" "config-directory" "noMap" true }}
{{ with resources.Get "css/main.css" | postCSS $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

## No configuration file

To avoid using a PostCSS configuration file, you can specify a minimal configuration using the options map.

use
: (`string`) A space-delimited list of PostCSS plugins to use.

parser
: (`string`) A custom PostCSS parser.

stringifier
: (`string`) A custom PostCSS stringifier.

syntax
: (`string`) Custom postcss syntax.

```go-html-template
{{ $opts := dict "use" "autoprefixer postcss-color-alpha" }}
{{ with resources.Get "css/main.css" | postCSS $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

## Check environment

The current Hugo environment name (set by `--environment` or in configuration or OS environment) is available in the Node context, which allows constructs like this:

```js
const autoprefixer = require('autoprefixer');
const purgecss = require('@fullhuman/postcss-purgecss');
module.exports = {
  plugins: [
    autoprefixer,
    process.env.HUGO_ENVIRONMENT !== 'development' ? purgecss : null
  ]
}
```

[node.js]: https://nodejs.org/en/download
[postcss plugins]: https://www.postcss.parts/
[supported file name]: https://github.com/postcss/postcss-load-config#usage
[transpile to CSS]: /functions/css/sass/
