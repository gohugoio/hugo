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

{{< new-in 0.128.0 />}}

Use the `css.TailwindCSS` function to process your Tailwind CSS files. This function uses the Tailwind CSS CLI to:

1. Scan your templates for Tailwind CSS utility class usage.
1. Compile those utility classes into standard CSS.
1. Generate an optimized CSS output file.

> [!note]
> Use this function with Tailwind CSS v4.0 and later, which require a relatively [modern browser] to render correctly.

[modern browser]: https://tailwindcss.com/docs/compatibility#browser-support

## Setup

### Step 1

Install the Tailwind CSS CLI v4.0 or later:

```sh {copy=true}
npm install --save-dev tailwindcss @tailwindcss/cli
```

The Tailwind CSS CLI is also available as a [standalone executable]. You must install it outside of your project directory and ensure its path is included in your system's `PATH` environment variable.


[standalone executable]: https://github.com/tailwindlabs/tailwindcss/releases/latest

### Step 2

Add this to your site configuration:

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

### Step 3

Create a CSS entry file:

```css {file="assets/css/main.css" copy=true}
@import "tailwindcss";
@source "hugo_stats.json";
```

Tailwind CSS respects `.gitignore` files. This means that if `hugo_stats.json` is listed in your `.gitignore` file, Tailwind CSS will ignore it. To make `hugo_stats.json` available to Tailwind CSS you must explicitly source it as shown in the example above.

### Step 4

Create a partial template to process the CSS with the Tailwind CSS CLI:

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

### Step 5

Call the partial template from your base template, deferring template execution until after all sites and output formats have been rendered:

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
