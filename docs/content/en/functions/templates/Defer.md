---
title: templates.Defer
description: Defers execution of a template until all sites and output formats have been rendered.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [templates.Defer OPTIONS]
aliases: [/functions/templates.defer]
---

## Usage

The `templates.Defer` function defers the execution of a template until all sites and output formats have been rendered.

> [!NOTE]
> Do not call this function within a _partial_ template called with the `partialCached` function. This restriction applies transitively: if `partialCached` calls a _partial_ template that calls `templates.Defer`, Hugo returns an error. Using this function within shortcode or render hook templates may also lead to unpredictable results.

> [!NOTE]
> This function only works in combination with the `with` keyword.
>
> Variables defined on the outside are not visible on the inside and vice versa. To pass in data, use the `data` option.

The `templates.Defer` function requires a single argument, a map with the following optional keys:

`key`
: (`string`) The key to use for the deferred template. This will, combined with a hash of the template content, be used as a cache key. Without this key, Hugo executes the deferred template on every render, which is inefficient for shared resources like CSS and JavaScript.

`data`
: (`map`) Optional map to pass as data to the deferred template. This will be available in the deferred template as `.` or `$`.

```go-html-template
Language Outside: {{ site.Language.Name }}
Page Outside: {{ .RelPermalink }}
I18n Outside: {{ i18n "hello" }}
{{ $data := (dict "page" . )}}
{{ with (templates.Defer (dict "data" $data )) }}
     Language Inside: {{ site.Language.Name }}
     Page Inside: {{ .page.RelPermalink }}
     I18n Inside: {{ i18n "hello" }}
{{ end }}
```

The [output format][], [site][], and [language][] will be the same, even if the execution is deferred. In the example above, this means that the `site.Language.Name` and `.RelPermalink` will be the same on the inside and the outside of the deferred template.

## Examples

The following examples defer CSS processing until Hugo has finished rendering all pages.

### Process Tailwind CSS

The [`css.TailwindCSS`][] function uses [`hugo_stats.json`][] to determine which classes and other HTML identifiers are in use. Because Hugo must finish rendering all pages before writing this file, use the following templates to defer the CSS processing:

```go-html-template {file="layouts/baseof.html" copy=true}
<head>
  {{ with (templates.Defer (dict "key" "global")) }}
    {{ partial "css.html" . }}
  {{ end }}
</head>
```

```go-html-template {file="layouts/_partials/css.html" copy=true}
{{ with resources.Get "css/main.css" }}
  {{ $opts := dict "minify" (not hugo.IsDevelopment) }}
  {{ with . | css.TailwindCSS $opts }}
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

For the above to work well when running the server (or `hugo -w`), add the following configuration:

{{< code-toggle file=hugo >}}
[build]
  [build.buildStats]
    enable = true
  [[build.cachebusters]]
    source = 'assets/notwatching/hugo_stats\.json'
    target = 'css'
  [[build.cachebusters]]
    source = '(postcss|tailwind)\.config\.js'
    target = 'css'
[module]
  [[module.mounts]]
    source = 'assets'
    target = 'assets'
  [[module.mounts]]
    disableWatch = true
    source = 'hugo_stats.json'
    target = 'assets/notwatching/hugo_stats.json'
{{< /code-toggle >}}

### Remove unused CSS

The `postcss.config.mjs` configuration file created in [Step 4](#step-4) directs the [`PurgeCSS`][] plugin to read [`hugo_stats.json`][] when the [`css.PostCSS`][] function runs. Because Hugo only writes this file after all pages are rendered, the CSS processing must be deferred.

Step 1
: Install [Node.js][].

Step 2
: Install the required Node packages in the root of your project:

  ```sh {copy=true}
  npm i -D postcss postcss-cli autoprefixer @fullhuman/postcss-purgecss
  ```

Step 3
: Enable creation of the [`hugo_stats.json`][] file when building the site. If you are only using this for the production build, consider placing it below [`config/production`][].

  {{< code-toggle file=hugo copy=true >}}
  [build.buildStats]
  enable = true
  {{< /code-toggle >}}

  See the [configure build][] documentation for details and options.

Step 4
: Create a configuration file for the `PurgeCSS` and `autoprefixer` plugins in the root of your project.

  ```js {file="postcss.config.mjs" copy=true}
  import autoprefixer from 'autoprefixer';
  import purgeCSSPlugin from '@fullhuman/postcss-purgecss';

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

  export default {
    plugins: [
      process.env.HUGO_ENVIRONMENT !== 'development' ? purgecss : null,
      autoprefixer,
    ]
  };
  ```

Step 5
: Place your CSS file within the `assets/css` directory.

Step 6
: Defer CSS processing so it runs after the site has been rendered:

  ```go-html-template {file="layouts/baseof.html" copy=true}
  <head>
    {{ with (templates.Defer (dict "key" "global")) }}
      {{ partial "css.html" . }}
    {{ end }}
  </head>
  ```

  ```go-html-template {file="layouts/_partials/css.html" copy=true}
  {{ with resources.Get "css/main.css" }}
    {{ if hugo.IsDevelopment }}
      <link rel="stylesheet" href="{{ .RelPermalink }}">
    {{ else }}
      {{ with . | postCSS | minify | fingerprint }}
        <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
      {{ end }}
    {{ end }}
  {{ end }}
  ```

[Node.js]: https://nodejs.org/en
[`PurgeCSS`]: https://github.com/FullHuman/purgecss
[`config/production`]: /configuration/introduction/#configuration-directory
[`css.PostCSS`]: /functions/css/postcss/
[`css.TailwindCSS`]: /functions/css/tailwindcss/
[`hugo_stats.json`]: /configuration/build/
[configure build]: /configuration/build/
[language]: /methods/site/language/
[output format]: /configuration/output-formats/
[site]: /methods/page/site/
