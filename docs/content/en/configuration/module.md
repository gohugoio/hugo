---
title: Configure modules
linkTitle: Modules
description: Configure modules.
categories: []
keywords: []
aliases: [/hugo-modules/configuration/]
---

{{% include "/_common/gomodules-info.md" %}}

## Top-level options

This is the default configuration:

<!-- markdownlint-disable MD049 -->
{{< code-toggle file=hugo >}}
[module]
noProxy = 'none'
noVendor = ''
private = '*.*'
proxy = 'direct'
vendorClosest = false
workspace = 'off'
{{< /code-toggle >}}
<!-- markdownlint-enable MD049 -->

auth
: {{< new-in 0.144.0 />}}
: (`string`) Configures `GOAUTH` when running the Go command for module operations. This is a semicolon-separated list of authentication commands for go-import and HTTPS module mirror interactions. This is useful for private repositories. See `go help goauth` for more information.

noProxy
: (`string`) A comma-separated list of [glob patterns](g),s matching paths that should not use the [configured proxy server](#proxy).

noVendor
: (`string`) A [glob pattern](g) matching module paths to skip when vendoring.

private
: (`string`) A comma-separated list of [glob patterns](g),s matching paths that should be treated as private.

proxy
: (`string`) The proxy server to use to download remote modules. Default is `direct`, which means `git clone` and similar.

replacements
: (`string`) Primarily useful for local module development, a comma-separated list of mappings from module paths to directories. Paths may be absolute or relative to the [`themesDir`][].

  {{< code-toggle file=hugo >}}
  [module]
  replacements = 'github.com/bep/my-theme -> ../..,github.com/bep/shortcodes -> /some/path'
  {{< /code-toggle >}}

vendorClosest
: (`bool`) Whether to pick the vendored module closest to the module using it. The default behavior is to pick the first. Note that there can still be only one dependency of a given module path, so once it is in use it cannot be redefined. Default is `false`.

workspace
: (`string`) The Go workspace file to use, either as an absolute path or a path relative to the current working directory. Enabling this activates Go workspace mode and requires Go 1.18 or later. The default is `off`.

You may also use environment variables to set any of the above. For example:

```sh
export HUGO_MODULE_PROXY="https://proxy.example.org"
export HUGO_MODULE_REPLACEMENTS="github.com/bep/my-theme -> ../.."
export HUGO_MODULE_WORKSPACE="/my/hugo.work"
```

## Hugo version

You can specify a required Hugo version for your module in the `module` section. Users will then receive a warning if their Hugo version is incompatible.

This is the default configuration:

{{< code-toggle config=module.hugoVersion />}}

You can omit any of the settings above.

extended
: (`bool`) Whether the extended edition of Hugo is required, satisfied by installing either the extended or extended/deploy edition.

  > [!note]
  > The extended version check is disabled in v0.153.2 and later.
  >
  > Historically, certain features—specifically WebP encoding and LibSass—required the Hugo Extended binary. However, as of v0.153.0:
  >
  > - WebP encoding is now supported in all Hugo editions.
  > - LibSass has been deprecated in favor of [Dart Sass][], which is compatible with any Hugo edition.
  >
  > Because these dependencies no longer require a specialized binary, the internal enforcement check for the extended version has been removed. Site and theme authors are encouraged to use Dart Sass to ensure cross-edition compatibility.

max
: (`string`) The maximum Hugo version supported, for example `0.153.0`.

min
: (`string`) The minimum Hugo version supported, for example `0.102.0`.

## Imports

{{< code-toggle file=hugo >}}
[[module.imports]]
disable = false
ignoreConfig = false
ignoreImports = false
path = "github.com/gohugoio/hugoTestModules1_linux/modh1_2_1v"
[[module.imports]]
path = "my-shortcodes"
{{< /code-toggle >}}

disable
: (`bool`) Whether to disable the module but keep version information in the `go.*` files. Default is `false`.

ignoreConfig
: (`bool`) Whether to ignore module configuration files, for example, `hugo.toml`. This will also prevent loading of any transitive module dependencies. Default is `false`.

ignoreImports
: (`bool`) Whether to ignore module imports. Default is `false`.

noMounts
: (`bool`) Whether to disable directory mounting for this import. Default is `false`.

noVendor
: (`bool`) Whether to disable vendoring for this import. This setting is restricted to the main project. Default is `false`.

path
: (`string`) The module path, either a valid Go module path (e.g., `github.com/gohugoio/myShortcodes`) or the directory name if stored in the [`themesDir`][].

version
: {{< new-in 0.150.0 />}}
: If set to a [version query](https://go.dev/ref/mod#version-queries), this import becomes a direct dependency, in contrast to dependencies managed by Go Modules. See [this issue](https://github.com/gohugoio/hugo/pull/13966) for more information.

## Mounts

{{% glossary-term mount %}}

> [!important]
> If you define one or more mounts to map a file system path to a component path, do not use these legacy configuration settings: [`archetypeDir`][], [`assetDir`][], [`contentDir`][], [`dataDir`][], [`i18nDir`][], [`layoutDir`][], or [`staticDir`][].

### Default mounts

Within a project, if you define a mount to map a file system path to a component path, the corresponding default mount for that component will be removed. This action essentially overwrites the standard, automatic mapping for that specific component with your custom one.

Within a module, if you define a mount to map a file system path to a component path, all of the default mounts will be removed. Defining a mount at the module level is a more sweeping change, causing all default mappings within that module to be discarded.

In either case, if you still need one of the default mounts, you must explicitly add it along with the new mount. Because custom mounts override defaults, any necessary default mappings must be re-added manually after you introduce your custom configuration.

These are the default mounts:

{{< code-toggle config=module.mounts />}}

source
: (`string`) The source directory of the mount. For the main project, this can be either project-relative or absolute. For other modules it must be project-relative.

target
: (`string`) Where the mount will reside within Hugo's [unified file system](g). It must begin with one of Hugo's [component](g) directories: archetypes, assets, content, data, i18n, layouts, or static. For example, content/blog.

disableWatch
: {{< new-in 0.128.0 />}}
: (`bool`) Whether to disable watching in watch mode for this mount. Default is `false`.

files
: {{< new-in 0.153.0 />}}
: (`[]string`) A [glob slice](g) defining the files to include or exclude.

sites
: {{< new-in 0.153.0 />}}
: (`map`) A map to define [sites matrix](g) and [sites complements](g) for the mount. Relevant for `content` and `layouts` mounts, and `static` mounts when in multihost mode. For `static` and `layouts`, only the `matrix` keyword is supported.

### Example

{{< code-toggle file=hugo >}}
[module]
[[module.mounts]]
    source="content"
    target="content"
    files=["! docs/*"]
[[module.mounts]]
    source="node_modules"
    target="assets"
[[module.mounts]]
    source="assets"
    target="assets"
{{< /code-toggle >}}

[`archetypeDir`]: /configuration/all/#archetypedir
[`assetDir`]: /configuration/all/#assetdir
[`contentDir`]: /configuration/all/#contentdir
[`dataDir`]: /configuration/all/#datadir
[`i18nDir`]: /configuration/all/#i18ndir
[`layoutDir`]: /configuration/all/#layoutdir
[`staticDir`]: /configuration/all/#staticdir
[`themesDir`]: /configuration/all/#themesdir
[Dart Sass]: /functions/css/sass/#dart-sass
