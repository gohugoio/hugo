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

{{< new-in 0.128.0 >}}

<!-- TODO remove this admonition when feature is stable. -->

{{% note %}}
This is an experimental feature pending the release of TailwindCSS v4.0.

The functionality, configuration requirements, and documentation are subject to change at any time and may be not compatible with prior releases.
{{% /note %}}

## Prerequisites

To use this function you must install the Tailwind CSS CLI v4.0 or later. You may install the CLI as an npm package or as a standalone executable. See the [Tailwind CSS documentation] for details.

[Tailwind CSS documentation]: https://tailwindcss.com/docs/installation

{{% note %}}
Use npm to install the CLI prior to the v4.0 release of Tailwind CSS.

`npm install --save-dev tailwindcss@next @tailwindcss/cli@next`
{{% /note %}}

## Options

minify
: (`bool`) Whether to optimize and minify the output. Default is `false`.

optimize
: (`bool`) Whether to optimize the output without minifying. Default is `false`.

inlineImports
: (`bool`) Whether to enable inlining of `@import` statements. Inlining is performed recursively, but currently once only per file. It is not possible to import the same file in different scopes (root, media query, etc.). Note that this import routine does not care about the CSS specification, so you can have `@import` statements anywhere in the file. Default is `false`.

skipInlineImportsNotFound
: (`bool`) When `inlineImports` is enabled, we fail the build if an import cannot be resolved. Enable this option to allow the build to continue and leave the import statement in place. Note that the inline importer does not process URL location or imports with media queries, so those will be left as-is even without enabling this option. Default is `false`.

## Example

Define a [cache buster] in your site configuration:

[cache buster]: /getting-started/configuration/#configure-cache-busters

{{< code-toggle file=hugo >}}
[[build.cachebusters]]
source = 'layouts/.*'
target = 'css'
{{< /code-toggle >}}

Process the resource:

```go-html-template
{{ with resources.Get "css/main.css" }}
  {{ $opts := dict "minify" true }}
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

The example above publishes the minified CSS file to public/css/main.css.

See [this repository] for more information about the integration with Tailwind CSS v4.0.

[this repository]: https://github.com/bep/hugo-testing-tailwindcss-v4
