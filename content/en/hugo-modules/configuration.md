---
title: Configure Modules
linktitle: Configure Modules
description: This page describes the configuration options for a module.
date: 2019-07-24
categories: [hugo modules]
keywords: [themes, source, organization, directories]
menu:
  docs:
    parent: "modules"
    weight: 10
weight: 10
sections_weight: 10
toc: true
---

## Module Config: Top level

{{< code-toggle file="config">}}
[module]
proxy = "direct"
noProxy = "none"
private = "*.*"
{{< /code-toggle >}}

proxy
: Defines the proxy server to use to download remote modules. Default is `direct`, which means "git clone" and similar.

noProxy
: Comma separated glob list matching paths that should not use the proxy configured above.

private
: Comma separated glob list matching paths that should be treated as private.

Note that the above terms maps directly to their counterparts in Go Modules. Some of these setting may be natural to set as OS environment variables. To set the proxy server to use, as an example:

```
env HUGO_MODULE_PROXY=https://proxy.example.org hugo
```

{{< gomodules-info >}}

## Module Config: hugoVersion

If your module requires a particular version of Hugo to work, you can indicate that in the `module` section and the user will be warned if using a too old/new version.

{{< code-toggle file="config">}}
[module]
[module.hugoVersion]
  min = ""
  max = ""
  extended = false

{{< /code-toggle >}}

Any of the above can be omitted.

min
: The minimum Hugo version supported, e.g. `0.55.0`

max 
: The maximum Hugo version supported, e.g. `0.55.0`

extended
: Whether the extended version of Hugo is required.

## Module Config: imports

{{< code-toggle file="config">}}
[module]
[[module.imports]]
  path = "github.com/gohugoio/hugoTestModules1_linux/modh1_2_1v"
  ignoreConfig = false
  disable = false
[[module.imports]]
  path = "my-shortcodes"  
{{< /code-toggle >}}

path
: Can be either a valid Go Module module path, e.g. `github.com/gohugoio/myShortcodes`, or the directory name for the module as stored in your themes folder.

ignoreConfig
: If enabled, any module configuration file, e.g. `config.toml`, will not be loaded. Note that this will also stop the loading of any transitive module dependencies.

disable
: Set to `true` to disable the module off while keeping any version info in the `go.*` files.

{{< gomodules-info >}}


## Module Config: mounts

{{% note %}}
When the `mounts` config was introduced in Hugo 0.56.0, we were careful to preserve the existing `staticDir` and similar configuration to make sure all existing sites just continued to work.

But you should not have both. So if you add a `mounts` section you should make it complete and remove the old `staticDir` etc. settings.
{{% /note %}}

{{< code-toggle file="config">}}
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
: The source directory of the mount. For the main project, this can be either project-relative or absolute and even a symbolic link. For other modules it must be project-relative.

target
: Where it should be mounted into Hugo's virtual filesystem. It must start with one of Hugo's component folders: `static`, `content`, `layouts`, `data`, `assets`, `i18n`, or `archetypes`. E.g. `content/blog`.

lang
: The language code, e.g. "en". Only relevant for `content` mounts, and `static` mounts when in multihost mode.

