---
title: Configure modules
linkTitle: Modules
description: Configure modules.
categories: []
keywords: []
aliases: [/hugo-modules/configuration/]
---

## Top-level options

This is the default configuration:

{{< code-toggle file=hugo >}}
[module]
noProxy = 'none'
noVendor = ''
private = '*.*'
proxy = 'direct'
vendorClosest = false
workspace = 'off'
{{< /code-toggle >}}

auth
: {{< new-in 0.144.0 />}}
: (`string`) Configures `GOAUTH` when running the Go command for module operations. This is a semicolon-separated list of authentication commands for go-import and HTTPS module mirror interactions. This is useful for private repositories. See `go help goauth` for more information.

noProxy
: (`string`) A comma-separated list of [glob](g) patterns matching paths that should not use the [configured proxy server](#proxy).

noVendor
: (`string`) A [glob](g) pattern matching module paths to skip when vendoring.

private
: (`string`) A comma-separated list of [glob](g) patterns matching paths that should be treated as private.

proxy
: (`string`) The proxy server to use to download remote modules. Default is `direct`, which means `git clone` and similar.

replacements
: (`string`) Primarily useful for local module development, a comma-separated list of mappings from module paths to directories. Paths may be absolute or relative to the [`themesDir`].

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

{{% include "/_common/gomodules-info.md" %}}

## Hugo version

You can specify a required Hugo version for your module in the `module` section. Users will then receive a warning if their Hugo version is incompatible.

This is the default configuration:

{{< code-toggle config=module.hugoVersion />}}

You can omit any of the settings above.

extended
: (`bool`) Whether the extended edition of Hugo is required, satisfied by installing either the extended or extended/deploy edition.

max
: (`string`) The maximum Hugo version supported, for example `0.143.0`.

min
: (`string`) The minimum Hugo version supported, for example `0.123.0`.

[`themesDir`]: /configuration/all/#themesdir

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
: (`string`) The module path, either a valid Go module path (e.g., `github.com/gohugoio/myShortcodes`) or the directory name if stored in the [`themesDir`].

[`themesDir`]: /configuration/all#themesDir

{{% include "/_common/gomodules-info.md" %}}

## Mounts

Before Hugo v0.56.0, custom component paths could only be configured by setting [`archetypeDir`], [`assetDir`], [`contentDir`], [`dataDir`], [`i18nDir`], [`layoutDi`], or [`staticDir`] in the site configuration. Module mounts offer greater flexibility than these legacy settings, but
you cannot use both.

[`archetypeDir`]: /configuration/all/
[`assetDir`]: /configuration/all/
[`contentDir`]: /configuration/all/
[`dataDir`]: /configuration/all/
[`i18nDir`]: /configuration/all/
[`layoutDi`]: /configuration/all/
[`staticDir`]: /configuration/all/

> [!note]
> If you use module mounts do not use the legacy settings.

### Default mounts

> [!note]
> Adding a new mount to a target root will cause the existing default mount for that root to be ignored. If you still need the default mount, you must explicitly add it along with the new mount.

The are the default mounts:

{{< code-toggle config=module.mounts />}}

source
: (`string`) The source directory of the mount. For the main project, this can be either project-relative or absolute. For other modules it must be project-relative.

target
: (`string`) Where the mount will reside within Hugo's virtual file system. It must begin with one of Hugo's component directories: `archetypes`, `assets`, `content`, `data`, `i18n`, `layouts`, or `static`. For example, `content/blog`.

disableWatch
: {{< new-in 0.128.0 />}}
: (`bool`) Whether to disable watching in watch mode for this mount. Default is `false`.

lang
: (`string`) The language code, e.g. "en". Relevant for `content` mounts, and `static` mounts when in multihost mode.

includeFiles
: (`string` or `[]string`) One or more [glob](g) patterns matching files or directories to include. If `excludeFiles` is not set, the files matching `includeFiles` will be the files mounted.

  The glob patterns are matched against file names relative to the source root. Use Unix-style forward slashes (`/`), even on Windows. A single forward slash (`/`) matches the mount root, and double asterisks (`**`) act as a recursive wildcard, matching all directories and files beneath a given point (e.g., `/posts/**.jpg`). The search is case-insensitive.

excludeFiles
: (`string` or `[]string`) One or more [glob](g) patterns matching files to exclude.

### Example

{{< code-toggle file=hugo >}}
[module]
[[module.mounts]]
    source="content"
    target="content"
    excludeFiles="docs/*"
[[module.mounts]]
    source="node_modules"
    target="assets"
[[module.mounts]]
    source="assets"
    target="assets"
{{< /code-toggle >}}
