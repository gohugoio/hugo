---
title: Configure build
linkTitle: Build
description: Configure global build options.
categories: []
keywords: []
aliases: [/getting-started/configuration-build/]
---

The `build` configuration section contains global build-related configuration options.

{{< code-toggle config=build />}}

buildStats
: See the [build stats](#build-stats) section below.

cachebusters
: See the [cache busters](#cache-busters) section below.

noJSConfigInAssets
: (`bool`) Whether to disable writing a `jsconfig.json` in your `assets` directory with mapping of imports from running [js.Build](/hugo-pipes/js). This file is intended to help with intellisense/navigation inside code editors such as [VS Code](https://code.visualstudio.com/). Note that if you do not use `js.Build`, no file will be written.

useResourceCacheWhen
: (`string`) When to use the resource file cache, one of `never`, `fallback`, or `always`. Applicable when transpiling Sass to CSS. Default is `fallback`.

## Cache busters

The `build.cachebusters` configuration option was added to support development using Tailwind 3.x's JIT compiler where a `build` configuration may look like this:

{{< code-toggle file=hugo >}}
[build]
  [build.buildStats]
    enable = true
  [[build.cachebusters]]
    source = "assets/watching/hugo_stats\\.json"
    target = "styles\\.css"
  [[build.cachebusters]]
    source = "(postcss|tailwind)\\.config\\.js"
    target = "css"
  [[build.cachebusters]]
    source = "assets/.*\\.(js|ts|jsx|tsx)"
    target = "js"
  [[build.cachebusters]]
    source = "assets/.*\\.(.*)$"
    target = "$1"
{{< /code-toggle >}}

When `buildStats` is enabled, Hugo writes a `hugo_stats.json` file on each build with HTML classes etc. that's used in the rendered output. Changes to this file will trigger a rebuild of the `styles.css` file. You also need to add `hugo_stats.json` to Hugo's server watcher. See [Hugo Starter Tailwind Basic](https://github.com/bep/hugo-starter-tailwind-basic) for a running example.

source
: (`string`) A [regular expression](g) matching file(s) relative to one of the virtual component directories in Hugo, typically `assets/...`.

target
: (`string`) A [regular expression](g) matching the keys in the resource cache that should be expired when `source` changes. You can use the matching regexp groups from `source` in the expression, e.g. `$1`.

## Build stats

{{< code-toggle config=build.buildStats />}}

enable
: (`bool`) Whether to create a `hugo_stats.json` file in the root of your project. This file contains arrays of the `class` attributes, `id` attributes, and tags of every HTML element within your published site. Use this file as data source when [removing unused CSS] from your site. This process is also known as pruning, purging, or tree shaking. Default is `false`.

[removing unused CSS]: /functions/resources/postprocess/

disableIDs
: (`bool`) Whether to exclude `id` attributes. Default is `false`.

disableTags
: (`bool`) Whether to exclude element tags. Default is `false`.

disableClasses
: (`bool`) Whether to exclude `class` attributes. Default is `false`.

> [!note]
> Given that CSS purging is typically limited to production builds, place the `buildStats` object below [`config/production`].
>
> Built for speed, there may be "false positive" detections (e.g., HTML elements that are not HTML elements) while parsing the published site. These "false positives" are infrequent and inconsequential.

Due to the nature of partial server builds, new HTML entities are added while the server is running, but old values will not be removed until you restart the server or run a regular `hugo` build.

[`config/production`]: /configuration/introduction/#configuration-directory
