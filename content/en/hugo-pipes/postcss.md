---
title: PostCSS
description: Process CSS files with PostCSS, using any of the available plugins.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    weight: 40
weight: 40
toc: true
action:
  aliases: [postCSS]
  returnType: resource.Resource
  signatures: ['css.PostCSS [OPTIONS] RESOURCE']
---

## Setup

Follow the steps below to transform CSS using any of the [available PostCSS plugins](https://www.postcss.parts/).

Step 1
: Install [Node.js](https://nodejs.org/en/download).

Step 2
: Install the required Node.js packages in the root of your project. For example, to add vendor prefixes to CSS rules:

```sh
npm i -D postcss postcss-cli autoprefixer
```

Step 3
: Create a PostCSS configuration file in the root of your project. You must name this file `postcss.config.js` or one of the other [supported file names]. For example:

[supported file names]: https://github.com/postcss/postcss-load-config#usage

{{< code file=postcss.config.js >}}
module.exports = {
  plugins: [
    require('autoprefixer')
  ]
};
{{< /code >}}

{{% note %}}
If you are a Windows user, and the path to your project contains a space, you must place the PostCSS configuration within the package.json file. See [this example](https://github.com/postcss/postcss-load-config#packagejson) and issue [#7333](https://github.com/gohugoio/hugo/issues/7333).
{{% /note %}}

Step 4
: Place your CSS file within the `assets` directory.

Step 5
: Capture the CSS file as a resource and pipe it through `css.PostCSS` (alias `postCSS`):

{{< code file=layouts/partials/css.html >}}
{{ with resources.Get "css/main.css" | postCSS }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
{{< /code >}}

If starting with a Sass file within the `assets` directory:

{{< code file=layouts/partials/css.html >}}
{{ with resources.Get "sass/main.scss" | toCSS | postCSS }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
{{< /code >}}

## Options

The `css.PostCSS` method takes an optional map of options.

config
: (`string`) The directory that contains the PostCSS configuration file. Default is the root of the project directory.

noMap
: (`bool`) Default is `false`. If `true`, disables inline sourcemaps.

inlineImports
: (`bool`) Default is `false`. Enable inlining of @import statements. It does so recursively, but will only import a file once.
URL imports (e.g. `@import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');`) and imports with media queries will be ignored.
Note that this import routine does not care about the CSS spec, so you can have @import anywhere in the file.
Hugo will look for imports relative to the module mount and will respect theme overrides.

skipInlineImportsNotFound
: (`bool`) Default is `false`. If you have regular CSS imports in your CSS that you want to preserve, you can either use imports with URL or media queries (Hugo does not try to resolve those) or set `skipInlineImportsNotFound` to true.

{{< code file=layouts/partials/css.html >}}
{{ $opts := dict "config" "config-directory" "noMap" true }}
{{ with resources.Get "css/main.css" | postCSS $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
{{< /code >}}

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

{{< code file=layouts/partials/css.html >}}
{{ $opts := dict "use" "autoprefixer postcss-color-alpha" }}
{{ with resources.Get "css/main.css" | postCSS $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
{{< /code >}}

## Check Hugo environment

The current Hugo environment name (set by `--environment` or in configuration or OS environment) is available in the Node context, which allows constructs like this:

{{< code file=postcss.config.js >}}
module.exports = {
  plugins: [
    require('autoprefixer'),
    ...process.env.HUGO_ENVIRONMENT === 'production'
      ? [purgecss]
      : []
  ]
}
{{< /code >}}
