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

The `resources.PostProcess` function delays resource transformation steps until the build is complete, primarily for tasks like removing unused CSS rules.

## Example

In this example, after the build is complete, Hugo will:

1. Purge unused CSS using the [PurgeCSS] plugin for [PostCSS]
2. Add vendor prefixes to CSS rules using the [Autoprefixer] plugin for PostCSS
3. [Minify] the CSS
4. [Fingerprint] the CSS 

[autoprefixer]: https://github.com/postcss/autoprefixer
[fingerprint]: /functions/resources/fingerprint/
[minify]: /functions/resources/minify/
[postcss]: /functions/css/postcss/
[purgecss]: https://purgecss.com/plugins/postcss.html

Step 1
: Install [Node.js].

[node.js]: https://nodejs.org/en/download

Step 2
: Install the required Node.js packages in the root of your project:

```sh
npm i -D postcss postcss-cli autoprefixer @fullhuman/postcss-purgecss
```

Step 3
: Enable creation of the `hugo_stats.json` file when building the site. If you are only using this for the production build, consider placing it below [`config/production`].

[`config/production`]: /getting-started/configuration/#configuration-directory

{{< code-toggle file=hugo >}}
[build.buildStats]
enable = true
{{< /code-toggle >}}

See the [configure build] documentation for details and options.

[configure build]: /getting-started/configuration/#configure-build

Step 4
: Create a PostCSS configuration file in the root of your project.

{{< code file="postcss.config.js" copy=true >}}
const autoprefixer = require('autoprefixer');
const purgeCSSPlugin = require('@fullhuman/postcss-purgecss').default;

const purgecss = purgeCSSPlugin({
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
    process.env.HUGO_ENVIRONMENT !== 'development' ? purgecss : null,
    autoprefixer,
  ]
};
{{< /code >}}

{{% note %}}
{{% include "functions/resources/_common/postcss-windows-warning.md" %}}
{{% /note %}}

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

Hugo passes the environment variables below to PostCSS, allowing you to do something like:

```js
process.env.HUGO_ENVIRONMENT !== 'development' ? purgecss : null,
```

PWD
: The absolute path to the project working directory.

HUGO_ENVIRONMENT
: The current Hugo environment, set with the `--environment` command line flag.
Default is `production` for `hugo` and `development` for `hugo server`.

HUGO_PUBLISHDIR
: The absolute path to the publish directory, typically `public`. This value points to a directory on disk, even when rendering to memory with the `--renderToMemory` command line flag.

HUGO_FILE_X
: Hugo automatically mounts the following files from your project's root directory under `assets/_jsconfig`:

- `babel.config.js`
- `postcss.config.js`
- `tailwind.config.js`

For each file, Hugo creates a corresponding environment variable named `HUGO_FILE_:filename:`, where `:filename:` is the uppercase version of the filename with periods replaced by underscores.  This allows you to access these files within your JavaScript, for example:

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
