---
title: Configure build
linkTitle: Build
description: Configure global build options.
categories: []
keywords: []
---

This is the default configuration:

{{< code-toggle config=build />}}

`buildStats`
: See the [build stats](#build-stats) section below.

`cachebusters`
: See the [cache busters](#cache-busters) section below.

`cleanDestinationDir`
: See the [clean destination directory](#clean-destination-directory) section below.

`noJSConfigInAssets`
: (`bool`) Whether to disable writing a `jsconfig.json` in your `assets` directory with mapping of imports from running [`js.Build`][]. This file is intended to help with intellisense/navigation inside code editors such as [VS Code][]. Note that if you do not use `js.Build`, no file will be written.

`useResourceCacheWhen`
: (`string`) When to use the resource file cache, one of `never`, `fallback`, or `always`. Applicable when transpiling Sass to CSS. Default is `fallback`.

## Build stats

{{< code-toggle config=build.buildStats />}}

`enable`
: (`bool`) Whether to create a `hugo_stats.json` file in the root of your project. This file contains arrays of the `class` attributes, `id` attributes, and tags of every HTML element within your published site. Use this file as data source when [removing unused CSS][] from your site. This process is also known as pruning, purging, or tree shaking. Default is `false`.

`disableIDs`
: (`bool`) Whether to exclude `id` attributes. Default is `false`.

`disableTags`
: (`bool`) Whether to exclude element tags. Default is `false`.

`disableClasses`
: (`bool`) Whether to exclude `class` attributes. Default is `false`.

> [!NOTE]
> Given that CSS purging is typically limited to production builds, place the `buildStats` object below [`config/production`][].
>
> Built for speed, there may be "false positive" detections, such as HTML elements that are not HTML elements, while parsing the published site. These "false positives" are infrequent and inconsequential.

Due to the nature of partial server builds, new HTML entities are added while the server is running, but old values will not be removed until you restart the server or run `hugo build`.

## Cache busters

Use `build.cachebusters` to expire specific keys in the resource cache when a watched source file changes, triggering a rebuild of dependent resources such as CSS. For example, use this configuration when using the [`css.TailwindCSS`][] function:

<!-- markdownlint-disable MD049 -->
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
[security]
  [security.exec]
    allow = ['^(dart-)?sass$', '^go$', '^git$', '^node$', '^postcss$', '^tailwindcss$']
{{< /code-toggle >}}
<!-- markdownlint-enable MD049 -->

When [`buildStats`](#buildstats) is enabled, Hugo writes a `hugo_stats.json` file on each build, containing the classes, IDs, and tags used in the rendered output. Changes to this file trigger a rebuild of the CSS. See the [`css.TailwindCSS`][] function for a running example.

`source`
: (`string`) A [regular expression](g) matching file(s) relative to one of the virtual component directories in Hugo, typically `assets/...`.

`target`
: (`string`) A [regular expression](g) matching the keys in the resource cache that should be expired when `source` changes. You can use the matching regexp groups from `source` in the expression, such as `$1`.

## Clean destination directory

Hugo does not clear the `public` directory before building your project. Existing files are overwritten, but not deleted. This behavior is intentional, preventing the inadvertent removal of files that you may have added to the `public` directory after the build.

As a result, the `public` directory can accumulate stale files over time. For example, a rendered page may remain after you delete or rename its content file, or draft, expired, and future content may remain after it no longer meets the criteria for publication. Enable `cleanDestinationDir` to have Hugo remove these stale files automatically on every build.

This is the default configuration:

{{< code-toggle config=build.cleanDestinationDir />}}

`enable`
: {{< new-in 0.167.0 />}}
: (`bool`) Whether to clean the `public` directory before rendering the site. Hugo removes every file and directory in the `public` directory that does not have a corresponding static file, whether from the project's `static` directory, a module mount, or a theme. This removes stale files, such as old rendered pages and deleted static assets, as well as files you added to the `public` directory yourself, such as a `CNAME` or `_redirects` file. Use the [`keepDirs`](#keepdirs) and [`keepFiles`](#keepfiles) settings to preserve specific directories and files. This cleanup runs even if the project has no static files. Default is `false`. Override this setting for a single build with the `--cleanDestinationDir` command line flag.

`keepDirs`
: {{< new-in 0.167.0 />}}
: (`[]string`) A [glob slice](g) matching directories, relative to the `public` directory, to preserve when cleaning the destination directory. In a multilingual multihost project, patterns are relative to each language's subdirectory of the `public` directory. A matching directory is kept along with everything beneath it, including subdirectories and their contents. The default value, shown above, matches directories whose names begin with a dot, wherever they occur in the directory tree. A value you set replaces the default rather than adding to it, so include the default pattern to continue preserving these directories.

`keepFiles`
: {{< new-in 0.167.0 />}}
: (`[]string`) A [glob slice](g) matching files, relative to the `public` directory, to preserve when cleaning the destination directory. In a multilingual multihost project, patterns are relative to each language's subdirectory of the `public` directory. The default value, shown above, matches `.git`, `.gitignore`, and `.gitattributes` files, wherever they occur in the directory tree. A value you set replaces the default rather than adding to it, so include the default pattern to continue preserving these files.

[VS Code]: https://code.visualstudio.com/
[`config/production`]: /configuration/introduction/#configuration-directory
[`css.TailwindCSS`]: /functions/css/tailwindcss/
[`js.Build`]: /functions/js/build/
[removing unused CSS]: /functions/templates/defer/#remove-unused-css
