---
title: css.TailwindCSS
description: Processes the given resource with the Tailwind CSS CLI.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/Fingerprint
    - functions/resources/Minify
    - functions/css/PostCSS
  returnType: resource.Resource
  signatures: ['css.TailwindCSS [OPTIONS] RESOURCE']
toc: true
---

{{< new-in 0.128.0 />}}

Use the `css.TailwindCSS` function to process your Tailwind CSS files.  This function uses the Tailwind CSS CLI to:

1. Scan your templates for Tailwind CSS utility class usage.
1. Compile those utility classes into standard CSS.
1. Generate an optimized CSS output file.

## Setup

Step 1
: Install the Tailwind CSS CLI v4.0 or later:

```sh
npm install --save-dev tailwindcss @tailwindcss/cli
```

The TailwindCSS CLI is also available as a [standalone executable] if you want to use it without installing Node.js.

[standalone executable]: https://github.com/tailwindlabs/tailwindcss/releases/latest

Step 2
: Add this to your site configuration:

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


Step 3
: Create a CSS entry file:

{{< code file=assets/css/main.css copy=true >}}
@import "tailwindcss";
@source "hugo_stats.json";
{{< /code >}}

Tailwind CSS respects `.gitignore` files. This means that if `hugo_stats.json` is listed in your `.gitignore` file, Tailwind CSS will ignore it. To make `hugo_stats.json` available to Tailwind CSS you must explicitly source it as shown in the example above.

Step 4
: Create a partial template to process the CSS with the Tailwind CSS CLI:

{{< code file=layouts/partials/css.html copy=true >}}
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
{{< /code >}}

Step 5
: Call the partial template from your base template:

{{< code file=layouts/default/baseof.html >}}
<head>
  ...
  {{ partialCached "css.html" . }}
  ...
<head>
{{< /code >}}

Step 6
: Optionally create a `tailwind.config.js` file in the root of your project as shown below. This is necessary if you use the [Tailwind CSS IntelliSense
extension] for Visual Studio Code.

[Tailwind CSS IntelliSense
extension]: https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss

{{< code file=tailwind.config.js copy=true >}}
/*
This file is present to satisfy a requirement of the Tailwind CSS IntelliSense
extension for Visual Studio Code.

https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss

The rest of this file is intentionally empty.
*/
{{< /code >}}

## Options

minify
: (`bool`) Whether to optimize and minify the output. Default is `false`.

optimize
: (`bool`) Whether to optimize the output without minifying. Default is `false`.

inlineImports
: (`bool`) Whether to enable inlining of `@import` statements. Inlining is performed recursively, but currently once only per file. It is not possible to import the same file in different scopes (root, media query, etc.). Note that this import routine does not care about the CSS specification, so you can have `@import` statements anywhere in the file. Default is `false`.

skipInlineImportsNotFound
: (`bool`) When `inlineImports` is enabled, we fail the build if an import cannot be resolved. Enable this option to allow the build to continue and leave the import statement in place. Note that the inline importer does not process URL location or imports with media queries, so those will be left as-is even without enabling this option. Default is `false`.
