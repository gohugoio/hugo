---
title: css.TailwindCSS
description: Processes the given resource with the Tailwind CSS CLI.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: ['css.TailwindCSS [OPTIONS] RESOURCE']
---

Use the `css.TailwindCSS` function to process your Tailwind CSS files. This function uses the Tailwind CSS CLI to:

1. Scan your templates for Tailwind CSS utility class usage.
1. Compile those utility classes into standard CSS.
1. Generate an optimized CSS output file.

> [!note]
> Use this function with Tailwind CSS v4.0 and later, which require a relatively [modern browser] to render correctly.

[modern browser]: https://tailwindcss.com/docs/compatibility#browser-support

## Setup

Step 1
: Install Tailwind CSS v4.0 or later:

  ```sh {copy=true}
  npm install --save-dev tailwindcss @tailwindcss/cli @tailwindcss/typography
  ```

  <!-- TODO: remove the admonition below somewhere after v0.172.0 -->
  
  > [!note]
  > As of v0.161.0, Hugo no longer supports the Tailwind [standalone binary]. You must now install the Tailwind CSS CLI via `npm` as shown above.

  [standalone binary]: https://github.com/tailwindlabs/tailwindcss/releases/latest

Step 2
: Add this to your project configuration:

  {{< code-toggle file=hugo copy=true >}}
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

Step 3
: Create a CSS entry file:

  ```css {file="assets/css/main.css" copy=true}
  @import "tailwindcss";
  @plugin "@tailwindcss/typography";
  @source "hugo_stats.json";
  ```

  Tailwind CSS respects `.gitignore` files. This means that if `hugo_stats.json` is listed in your `.gitignore` file, Tailwind CSS will ignore it. To make `hugo_stats.json` available to Tailwind CSS you must explicitly source it as shown in the example above.

Step 4
: Create a _partial_ template to process the CSS with the Tailwind CSS CLI:

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

Step 5
: Call the _partial_ template from your base template, deferring template execution until after all sites and output formats have been rendered:

  ```go-html-template {file="layouts/baseof.html" copy=true}
  <head>
    ...
    {{ with (templates.Defer (dict "key" "global")) }}
      {{ partial "css.html" . }}
    {{ end }}
    ...
  </head>
  ```

## Options

minify
: (`bool`) Whether to optimize and minify the output. Default is `false`.

optimize
: (`bool`) Whether to optimize the output without minifying. Default is `false`.

disableInlineImports
: {{< new-in 0.147.4 />}}
: (`bool`) Whether to disable inlining of `@import` statements. Inlining is performed recursively, but currently once only per file. It is not possible to import the same file in different scopes (root, media query, etc.). Note that this import routine does not care about the CSS specification, so you can have `@import` statements anywhere in the file. Default is `false`.

skipInlineImportsNotFound
: (`bool`) Whether to allow the build process to continue despite unresolved import statements, preserving the original import declarations. It is important to note that the inline importer does not process URL-based imports or those with media queries, and these will remain unaltered even when this option is disabled. Default is `false`.

## Inject CSS variables with `vars`

The [css.Build](./Build) function has a [vars](./Build#vars) option that can be used to inject CSS variables into your stylesheets. This is particularly useful for dynamically setting values based on your site's configuration or other data. To use this with Tailwind CSS, you can use [css.Build](./Build) as a preprocessor step before passing the result to `css.TailwindCSS`. Here's how you can do it:

```go-html-template
{{ with resources.Get "css/styles.css" }}
  {{ $cssOpts := dict
    "vars" (dict "favourite-color" "#7f93c9")
    "externals" (slice "tailwindcss")
  }}
  {{ $tailwindOpts := dict "disableInlineImports" true }}
  {{ with . | css.Build $cssOpts | css.TailwindCSS $tailwindOpts }}
    <link rel="stylesheet" href="{{ .RelPermalink }}">
  {{ end }}
{{ end }}
```

Some notes to the above:

- Marking `tailwindcss` as an external in the `css.Build` options prevents it from being processed by the build step, allowing it to be correctly handled by the Tailwind CSS CLI in the subsequent step.
- The `disableInlineImports` option is set to `true` for the Tailwind CSS step as imports are handled by the `css.Build`.
