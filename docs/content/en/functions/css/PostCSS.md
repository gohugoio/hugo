---
title: css.PostCSS
description: Returns a resource created by processing the given CSS resource with PostCSS.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [postCSS]
    returnType: resource.Resource
    signatures: ['css.PostCSS [OPTIONS] RESOURCE']
aliases: [/functions/resources/postcss/]
---

The `css.PostCSS` function transforms CSS using [PostCSS][] and any of its [plugins][].

## Setup

Step 1
: Install [Node.js][].

Step 2
: Install the required Node packages in the root of your project. For example, to install PostCSS, its command-line interface, and the plugin to automatically add vendor prefixes to your CSS:

  ```sh
  npm install --save-dev postcss postcss-cli autoprefixer
  ```

Step 3
: Create a PostCSS configuration file in the root of your project. Hugo exposes several [environment variables](#environment-variables) to the PostCSS process, including the current Hugo [environment](g) name. For example, in this configuration, running `hugo server` disables vendor prefixes but enables an inline sourcemap. Conversely, when building for production, it applies vendor prefixes and disables the sourcemap:

  ```js {file="postcss.config.mjs" copy=true}
  import autoprefixer from 'autoprefixer';

  const isDev = process.env.HUGO_ENVIRONMENT === 'development';

  export default {
    plugins: [
      !isDev ? autoprefixer : null
    ],
    map: isDev ? { inline: true } : false
  };
  ```

Step 4
: Place your CSS file within the `assets/css` directory.

Step 5
: Create a _partial_ template to process the CSS:

  ```go-html-template {file="layouts/_partials/css.html" copy=true}
  {{ with resources.Get "css/main.css" }}
    {{ $opts := dict
      "inlineImports" true
    }}
    {{ with . | css.PostCSS $opts }}
      {{ if hugo.IsDevelopment }}
        <link rel="stylesheet" href="{{ .RelPermalink }}">
      {{ else }}
        {{ with . | minify | fingerprint }}
          <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
        {{ end }}
      {{ end }}
    {{ end }}
  {{ end }}
  ```

Step 6
: Call the _partial_ template from your _base_ template:

  ```go-html-template {file="layouts/baseof.html" copy=true}
  <head>
    {{ partial "css.html" . }}
  </head>
  ```

## Options

The `css.PostCSS` function accepts an options map.

The following options apply when using a PostCSS configuration file:

`config`
: (`string`) The path to the directory that contains the PostCSS configuration file. By default, Hugo searches the root of the project directory and any modules for `postcss.config.js`, `postcss.config.mjs`, and `postcss.config.cjs` in that order. Use this option only if your configuration file is located in a custom subdirectory.

`importContext`
: {{< new-in 0.165.0 />}}
: (`resource.ResourceGetter`) A [resource getter](g) to use when resolving `@import` statements. Hugo searches this context first, matching by the path as written in the statement, before falling back to the file system. Applicable when [`inlineImports`](#inlineimports) is `true`.

`inlineImports`
: (`bool`) Whether to enable inlining of import statements. It does so recursively, but will only import a file once. Hugo looks for imports relative to the module mount and respects theme overrides. Default is `false`.

  Note that Hugo's internal import routine does not adhere to the CSS specification; you can place `@import` statements anywhere in the file. However, external URL imports and imports with media queries are ignored during the inlining process.

  The following snippet illustrates an external URL import that Hugo will ignore:

  ```css
  @import url('https://fonts.googleapis.com/css?family=Open+Sans&display=swap');
  ```

`skipInlineImportsNotFound`
: (`bool`) Whether to allow the build process to continue despite unresolved import statements, preserving the original import declarations. Set this option to `true` if you want to retain standard CSS imports unparsed. Default is `false`.

The following options apply when configuring PostCSS without a configuration file:

`noMap`
: (`bool`) Whether to disable the default inline source maps. Default is `false`.

`parser`
: (`string`) A custom PostCSS parser.

`stringifier`
: (`string`) A custom PostCSS stringifier.

`syntax`
: (`string`) Custom PostCSS syntax.

`use`
: (`string`) A space-delimited list of PostCSS [plugins][] to use.

For example, to pass your plugins and disable source maps directly through the options map instead of a configuration file:

```go-html-template {file="layouts/_partials/css.html" copy=true}
{{ with resources.Get "css/main.css" }}
  {{ $opts := dict
    "noMap" true
    "use" "autoprefixer postcss-color-alpha"
  }}
  {{ with . | css.PostCSS $opts }}
    {{ if hugo.IsDevelopment }}
      <link rel="stylesheet" href="{{ .RelPermalink }}">
    {{ else }}
      {{ with . | fingerprint }}
        <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
```

## Environment variables

Hugo passes the following environment variables to the PostCSS process, making them available within your PostCSS configuration file.

`PWD`
: The absolute path to the project working directory.

`HUGO_ENVIRONMENT`
: The current Hugo environment, set with the `--environment` command line flag.
Default is `production` for `hugo build` and `development` for `hugo server`.

`HUGO_PUBLISHDIR`
: The absolute path to the publish directory, typically `public`. This value points to a directory on disk, even when rendering to memory with the `--renderToMemory` command line flag.

`HUGO_FILE_FILENAME`
: Hugo automatically mounts the following files from your project's root directory under `assets/_jsconfig`:

- `babel.config.js`, `babel.config.mjs`, `babel.config.cjs`
- `postcss.config.js`, `postcss.config.mjs`, `postcss.config.cjs`
- `tailwind.config.js`, `tailwind.config.mjs`, `tailwind.config.cjs`

For each file, Hugo creates a corresponding environment variable named `HUGO_FILE_FILENAME`, where `FILENAME` is the uppercase version of the filename with periods replaced by underscores. This allows you to access these files within your JavaScript, for example:

```js
let tailwindConfig = process.env.HUGO_FILE_TAILWIND_CONFIG_JS || './tailwind.config.js';
```

[Node.js]: https://nodejs.org/en
[PostCSS]: https://postcss.org/
[plugins]: https://postcss.org/docs/postcss-plugins
