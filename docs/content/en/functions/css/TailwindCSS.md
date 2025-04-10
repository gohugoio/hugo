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

> [!caution]
> Tailwind CSS v4.0 and later requires a relatively [modern browser](https://tailwindcss.com/docs/compatibility#browser-support) to render correctly.

## Setup

### Step 1

Install the Tailwind CSS CLI v4.0 or later:

```sh
npm install --save-dev tailwindcss @tailwindcss/cli
```

The TailwindCSS CLI is also available as a [standalone executable] if you want to use it without installing Node.js.

[standalone executable]: https://github.com/tailwindlabs/tailwindcss/releases/latest

### Step 2

Add this to your site configuration:

{{< code-toggle file=hugo copy=true >}}
[[module.mounts]]
source = "assets"
target = "assets"
[[module.mounts]]
source = "hugo_stats.json"
target = "assets/notwatching/hugo_stats.json"
disableWatch = true
[build.buildStats]
enable = true
[[build.cachebusters]]
source = "assets/notwatching/hugo_stats\\.json"
target = "css"
[[build.cachebusters]]
source = "(postcss|tailwind)\\.config\\.js"
target = "css"
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

```go-html-template {file="layouts/partials/css.html" copy=true}
{{ with (templates.Defer (dict "key" "global")) }}
  {{ with resources.Get "css/main.css" }}
    {{ $opts := dict
      "minify" hugo.IsProduction
      "inlineImports" true
    }}
    {{ with . | css.TailwindCSS $opts }}
      {{ if hugo.IsProduction }}
        {{ with . | fingerprint }}
          <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
        {{ end }}
      {{ else }}
        <link rel="stylesheet" href="{{ .RelPermalink }}">
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
```

### Step 5

Call the partial template from your base template:

```go-html-template {file="layouts/_default/baseof.html"}
<head>
  ...
  {{ partialCached "css.html" . }}
  ...
<head>
```

### Step 6

Optionally create a `tailwind.config.js` file in the root of your project as shown below. This is necessary if you use the [Tailwind CSS IntelliSense
extension] for Visual Studio Code.

[Tailwind CSS IntelliSense
extension]: https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss

```js {file="tailwind.config.js" copy=true}
/*
This file is present to satisfy a requirement of the Tailwind CSS IntelliSense
extension for Visual Studio Code.

https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss

The rest of this file is intentionally empty.
*/
```

## Options

minify
: (`bool`) Whether to optimize and minify the output. Default is `false`.

optimize
: (`bool`) Whether to optimize the output without minifying. Default is `false`.

inlineImports
: (`bool`) Whether to enable inlining of `@import` statements. Inlining is performed recursively, but currently once only per file. It is not possible to import the same file in different scopes (root, media query, etc.). Note that this import routine does not care about the CSS specification, so you can have `@import` statements anywhere in the file. Default is `false`.

skipInlineImportsNotFound
: (`bool`) Whether to allow the build process to continue despite unresolved import statements, preserving the original import declarations. It is important to note that the inline importer does not process URL-based imports or those with media queries, and these will remain unaltered even when this option is disabled. Default is `false`.
