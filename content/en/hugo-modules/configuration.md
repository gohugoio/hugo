---
title: Configure Hugo modules
description: This page describes the configuration options for a module.
categories: [hugo modules]
keywords: [modules,themes]
menu:
  docs:
    parent: modules
    weight: 20
weight: 20
toc: true
---

## Module configuration: top level

{{< code-toggle file=hugo >}}
[module]
noProxy = 'none'
noVendor = ''
private = '*.*'
proxy = 'direct'
replacements = ''
vendorClosest = false
workspace = 'off'
{{< /code-toggle >}}

noProxy
: (`string`) Comma separated glob list matching paths that should not use the proxy configured above.

noVendor
: (`string`) A optional Glob pattern matching module paths to skip when vendoring, e.g. "github.com/**"

private
: (`string`) Comma separated glob list matching paths that should be treated as private.

proxy
: (`string`) Defines the proxy server to use to download remote modules. Default is `direct`, which means "git clone" and similar.

vendorClosest
: (`bool`) When enabled, we will pick the vendored module closest to the module using it. The default behavior is to pick the first. Note that there can still be only one dependency of a given module path, so once it is in use it cannot be redefined. Default is `false`.

workspace
: (`string`) The workspace file to use. This enables Go workspace mode. Note that this can also be set via OS env, e.g. `export HUGO_MODULE_WORKSPACE=/my/hugo.work` This only works with Go 1.18+. In Hugo `v0.109.0` we changed the default to `off` and we now resolve any relative work file names relative to the working directory.

replacements
: (`string`) A comma-separated list of mappings from module paths to directories, e.g. `github.com/bep/my-theme -> ../..,github.com/bep/shortcodes -> /some/path`. This is mostly useful for temporary local development of a module, in which case you might want to save it as an environment variable, e.g: `env HUGO_MODULE_REPLACEMENTS="github.com/bep/my-theme -> ../.."`. Relative paths are relative to [themesDir](/getting-started/configuration/#all-configuration-settings). Absolute paths are allowed.

Note that the above terms maps directly to their counterparts in Go Modules. Some of these setting may be natural to set as OS environment variables. To set the proxy server to use, as an example:

```txt
env HUGO_MODULE_PROXY=https://proxy.example.org hugo
```

{{< gomodules-info >}}

## Module configuration: hugoVersion

If your module requires a particular version of Hugo to work, you can indicate that in the `module` section and the user will be warned if using a too old/new version.

{{< code-toggle file=hugo >}}
[module]
[module.hugoVersion]
  min = ""
  max = ""
  extended = false

{{< /code-toggle >}}

Any of the above can be omitted.

min
: (`string`) The minimum Hugo version supported, e.g. `0.55.0`

max
: (`string`) The maximum Hugo version supported, e.g. `0.55.0`

extended
: (`bool`) Whether the extended version of Hugo is required.

## Module configuration: imports

{{< code-toggle file=hugo >}}
[module]
[[module.imports]]
  path = "github.com/gohugoio/hugoTestModules1_linux/modh1_2_1v"
  ignoreConfig = false
  ignoreImports = false
  disable = false
[[module.imports]]
  path = "my-shortcodes"
{{< /code-toggle >}}

path
: Can be either a valid Go Module module path, e.g. `github.com/gohugoio/myShortcodes`, or the directory name for the module as stored in your themes folder.

ignoreConfig
: If enabled, any module configuration file, e.g. `hugo.toml`, will not be loaded. Note that this will also stop the loading of any transitive module dependencies.

ignoreImports
: If enabled, module imports will not be followed.

disable
: Set to `true` to disable the module while keeping any version info in the `go.*` files.

noMounts
:  Do not mount any folder in this import.

noVendor
:  Never vendor this import (only allowed in main project).

{{< gomodules-info >}}

## Module configuration: mounts

{{% note %}}
When the `mounts` configuration was introduced in Hugo 0.56.0, we were careful to preserve the existing `contentDir`, `staticDir`, and similar configuration to make sure all existing sites just continued to work. But you should not have both: if you add a `mounts` section you should remove the old `contentDir`, `staticDir`, etc. settings.
{{% /note %}}

{{% note %}}
When you add a mount, the default mount for the concerned target root is ignored: be sure to explicitly add it.
{{% /note %}}

### Default mounts

{{< code-toggle file=hugo >}}
[module]
[[module.mounts]]
    source="content"
    target="content"
[[module.mounts]]
    source="static"
    target="static"
[[module.mounts]]
    source="layouts"
    target="layouts"
[[module.mounts]]
    source="data"
    target="data"
[[module.mounts]]
    source="assets"
    target="assets"
[[module.mounts]]
    source="i18n"
    target="i18n"
[[module.mounts]]
    source="archetypes"
    target="archetypes"
{{< /code-toggle >}}

source
: (`string`) The source directory of the mount. For the main project, this can be either project-relative or absolute. For other modules it must be project-relative.

target
: (`string`) Where it should be mounted into Hugo's virtual filesystem. It must start with one of Hugo's component folders: `static`, `content`, `layouts`, `data`, `assets`, `i18n`, or `archetypes`. E.g. `content/blog`.

disableWatch
{{< new-in 0.128.0 >}}
: (`bool`) Whether to disable watching in watch mode for this mount. Default is `false`.

lang
: (`string`) The language code, e.g. "en". Only relevant for `content` mounts, and `static` mounts when in multihost mode.

includeFiles
: (`string` or `string slice`) One or more [glob](https://github.com/gobwas/glob) patterns matching files or directories to include. If `excludeFiles` is not set, the files matching `includeFiles` will be the files mounted.

The glob patterns are matched to the file names starting from the `source` root, they should have Unix styled slashes even on Windows, `/` matches the mount root and `**` can be used as a  super-asterisk to match recursively down all directories, e.g `/posts/**.jpg`.

The search is case-insensitive.

excludeFiles
: (`string` or `string slice`) One or more glob patterns matching files to exclude.

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
