---
title: css.PostCSS
description: Processes the given resource with PostCSS using any PostCSS plugin.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [postCSS]
    returnType: resource.Resource
    signatures: ['css.PostCSS [OPTIONS] RESOURCE']
---

{{< new-in 0.128.0 />}}

```go-html-template
{{ with resources.Get "css/main.css" | postCSS }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
```

## Setup

Follow the steps below to transform CSS using any of the available [PostCSS plugins].

### Step 1

Install [Node.js].

### Step 2

Install the required Node.js packages in the root of your project. For example, to add vendor prefixes to your CSS rules:

```sh
npm i -D postcss postcss-cli autoprefixer
```

### Step 3

Create a PostCSS configuration file in the root of your project.

```js {file="postcss.config.js"}
module.exports = {
  plugins: [
    require('autoprefixer')
  ]
};
```

> [!note]
> If you are a Windows user, and the path to your project contains a space, you must place the PostCSS configuration within the package.json file. See [this example] and issue [#7333].

### Step 4

Place your CSS file within the `assets/css` directory.

### Step 5

Process the resource with PostCSS:

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
: (`bool`) Whether to disable inline source maps. Default is `false`.

inlineImports
: (`bool`) Whether to enable inlining of import statements. It does so recursively, but will only import a file once. URL imports (e.g. `@import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');`) and imports with media queries will be ignored. Note that this import routine does not care about the CSS spec, so you can have @import anywhere in the file. Hugo will look for imports relative to the module mount and will respect theme overrides. Default is `false`.

skipInlineImportsNotFound
: (`bool`) Whether to allow the build process to continue despite unresolved import statements, preserving the original import declarations. If you have regular CSS imports in your CSS that you want to preserve, you can either use imports with URL or media queries (Hugo does not try to resolve those) or set this option to `true`. Default is `false`."

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
module.exports = {
  plugins: [
    process.env.HUGO_ENVIRONMENT !== 'development' ? autoprefixer : null
  ]
}
```

[#7333]: https://github.com/gohugoio/hugo/issues/7333
[Node.js]: https://nodejs.org/en
[PostCSS plugins]: https://postcss.org/docs/postcss-plugins
[this example]: https://github.com/postcss/postcss-load-config#packagejson
