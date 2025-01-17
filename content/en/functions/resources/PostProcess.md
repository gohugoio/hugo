---
title: resources.PostProcess
description: Processes the given resource after the build.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/css/PostCSS
    - functions/css/Sass
  returnType: postpub.PostPublishedResource
  signatures: [resources.PostProcess RESOURCE]
toc: true
---

```go-html-template
{{ with resources.Get "css/main.css" }}
  {{ if hugo.IsDevelopment }}
    <link rel="stylesheet" href="{{ .RelPermalink }}">
  {{ else }}
    {{ with . | postCSS | minify | fingerprint | resources.PostProcess }}
      <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
    {{ end }}
  {{ end }}
{{ end }}
```

Marking a resource with `resources.PostProcess` postpones transformations until the build has finished.

Call `resources.PostProcess` when one or more of the steps in the transformation chain depends on the result of the build.

A prime use case for this is purging unused CSS rules using the [PurgeCSS] plugin for the PostCSS Node.js package.

## CSS Purging

{{% note %}}
There are several ways to set up CSS purging with PostCSS in Hugo. If you have a simple project, you should consider going the simpler route and drop the use of `resources.PostProcess` and just extract keywords from the templates. See the [Tailwind documentation](https://tailwindcss.com/docs/controlling-file-size/#app) for examples.
{{% /note %}}

Step 1
: Install [Node.js].

Step 2
: Install the required Node.js packages in the root of your project:

```sh
npm i -D postcss postcss-cli autoprefixer @fullhuman/postcss-purgecss
```

Step 3
: Create a PostCSS configuration file in the root of your project. You must name this file `postcss.config.js` or another [supported file name]. For example:

```js
const autoprefixer = require('autoprefixer');
const purgecss = require('@fullhuman/postcss-purgecss')({
  content: ['./hugo_stats.json'],
  defaultExtractor: content => {
    const els = JSON.parse(content).htmlElements;
    return [
      ...(els.tags || []),
      ...(els.classes || []),
      ...(els.ids || []),
    ];
  },
  // https://purgecss.com/safelisting.html
  safelist: []
});

module.exports = {
  plugins: [
    autoprefixer,
    process.env.HUGO_ENVIRONMENT !== 'development' ? purgecss : null
  ]
};
```

{{% note %}}
{{% include "functions/resources/_common/postcss-windows-warning.md" %}}
{{% /note %}}

Step 4
: Enable creation of the `hugo_stats.json` file when building the site. If you are only using this for the production build, consider placing it below [`config/production`].

{{< code-toggle file=hugo >}}
[build.buildStats]
enable = true
{{< /code-toggle >}}

See the [configure build] documentation for details and options.

Step 5
: Place your CSS file within the `assets/css` directory.

Step 6
: If the current environment is not `development`, process the resource with PostCSS:

```go-html-template
{{ with resources.Get "css/main.css" }}
  {{ if hugo.IsDevelopment }}
    <link rel="stylesheet" href="{{ .RelPermalink }}">
  {{ else }}
    {{ with . | postCSS | minify | fingerprint | resources.PostProcess }}
      <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
    {{ end }}
  {{ end }}
{{ end }}
```

## Environment variables

Hugo passes these environment variables to PostCSS, which allows you to do something like:

```js
process.env.HUGO_ENVIRONMENT === 'production' ? [autoprefixer] : []
```

PWD
: The absolute path to the project working directory.

HUGO_ENVIRONMENT
: The current Hugo environment, set with the `--environment` command line flag.
Default is `production` for `hugo` and `development` for `hugo server`.

HUGO_PUBLISHDIR
: The absolute path to the publish directory (the `public` directory). Note that the value will always point to a directory on disk even when running `hugo server` in memory mode. If you write to this directory from PostCSS when running the server, you could run the server with one of these flags:

```sh
hugo server --renderToDisk
hugo server --renderStaticToDisk
```

Also, Hugo will add environment variables for all files mounted below `assets/_jsconfig`. A default mount will be set up with files in the project root matching this regexp: `(babel|postcss|tailwind)\.config\.js`.

These will get environment variables named on the form `HUGO_FILE_:filename:` where `:filename:` is all upper case with periods replaced with underscore. This allows you to do something like:

```js
let tailwindConfig = process.env.HUGO_FILE_TAILWIND_CONFIG_JS || './tailwind.config.js';
```

## Limitations

Do not use `resources.PostProcess` when running Hugo's built-in development server. The examples above specifically prevent this by verifying that the current environment is not "development".

The `resources.PostProcess` function only works within templates that produce HTML files.

You cannot manipulate the values returned from the resource’s methods. For example, the `strings.ToUpper` function in this example will not work as expected:

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $css = $css | css.PostCSS | minify | fingerprint | resources.PostProcess }}
{{ $css.RelPermalink | strings.ToUpper }}
```

[node.js]: https://nodejs.org/en/download
[supported file name]: https://github.com/postcss/postcss-load-config#usage
[`config/production`]: /getting-started/configuration/#configuration-directory
[configure build]: /getting-started/configuration/#configure-build
[purgecss]: https://github.com/FullHuman/purgecss#readme
