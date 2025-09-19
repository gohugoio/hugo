---
title: resources.PostProcess
description: Processes the given resource after the build.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: postpub.PostPublishedResource
    signatures: [resources.PostProcess RESOURCE]
---

The `resources.PostProcess` function delays resource transformation steps until the build is complete, primarily for tasks like removing unused CSS rules.

## Example

In this example, after the build is complete, Hugo will:

1. Purge unused CSS using the [PurgeCSS] plugin for [PostCSS]
1. Add vendor prefixes to CSS rules using the [Autoprefixer] plugin for PostCSS
1. [Minify] the CSS
1. [Fingerprint] the CSS

Step 1
: Install [Node.js].

Step 2
: Install the required Node.js packages in the root of your project:

  ```sh {copy=true}
  npm i -D postcss postcss-cli autoprefixer @fullhuman/postcss-purgecss
  ```

Step 3
: Enable creation of the `hugo_stats.json` file when building the site. If you are only using this for the production build, consider placing it below [`config/production`].

  {{< code-toggle file=hugo copy=true >}}
  [build.buildStats]
  enable = true
  {{< /code-toggle >}}

  See the [configure build] documentation for details and options.

Step 4
: Create a PostCSS configuration file in the root of your project.

  ```js {file="postcss.config.js" copy=true}
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
  ```

  > [!note]
  > If you are a Windows user, and the path to your project contains a space, you must place the PostCSS configuration within the package.json file. See [this example] and issue [#7333].

Step 5
: Place your CSS file within the `assets/css` directory.

Step 6
: If the current environment is not `development`, process the resource with PostCSS:

  ```go-html-template {copy=true}
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

For each file, Hugo creates a corresponding environment variable named `HUGO_FILE_:filename:`, where `:filename:` is the uppercase version of the filename with periods replaced by underscores. This allows you to access these files within your JavaScript, for example:

```js
let tailwindConfig = process.env.HUGO_FILE_TAILWIND_CONFIG_JS || './tailwind.config.js';
```

## Limitations

Do not use `resources.PostProcess` when running Hugo's built-in development server. The examples above specifically prevent this by verifying that the current environment is not `development`.

The `resources.PostProcess` function only works within templates that produce HTML files.

You cannot manipulate the values returned from the resource's methods. For example, the `strings.ToUpper` function in this example will not work as expected:

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $css = $css | css.PostCSS | minify | fingerprint | resources.PostProcess }}
{{ $css.RelPermalink | strings.ToUpper }}
```

[#7333]: https://github.com/gohugoio/hugo/issues/7333
[`config/production`]: /configuration/introduction/#configuration-directory
[Autoprefixer]: https://github.com/postcss/autoprefixer
[configure build]: /configuration/build/
[Fingerprint]: /functions/resources/fingerprint/
[Minify]: /functions/resources/minify/
[Node.js]: https://nodejs.org/en
[PostCSS]: https://postcss.org/
[PurgeCSS]: https://github.com/FullHuman/purgecss
[this example]: https://github.com/postcss/postcss-load-config#packagejson
