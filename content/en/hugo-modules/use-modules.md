---
title: Use Hugo Modules
description: Use modules to manage the content, layout, presentation, and behavior of your site.
categories: []
keywords: []
weight: 20
aliases: [/themes/usage/,/themes/installing/,/installing-and-using-themes/]
---

> [!note]
> To work with modules you must install [Git][] and [Go][] 1.18 or later.

## Introduction

{{% glossary-term module %}}

- Modules can be imported in any combination or sequence.
- Module imports are recursive; importing Module A can trigger the import of Module B, and so on.
- Modules can provide configuration files and directories, subject to the constraints described in the [merge configuration settings][] section of the documentation.
- External directories, including those from non-Hugo projects, can be mounted to create a [unified file system](g).

## Import

To import a module, first initialize the project itself as a module. For example:

```sh
hugo mod init github.com/user/project
```

This will generate a [`go.mod`][] file in the project root.

> [!note]
> The module name is a unique identifier rather than a hosting requirement. Using a name like `github.com/user/project` is a common convention but it does not mean you must use Git or host your code on GitHub. You can use any name you like if you do not plan to have others import your project as a module. For example, you could use a simple name such as `my-project` when you run the initialization command.

Then define one or more imports in your project configuration. This contrived example imports three modules, each containing custom shortcodes:

{{< code-toggle file=hugo >}}
[module]
  [[module.imports]]
    path = 'shortcodes-a'
  [[module.imports]]
    path = '/home/user/shortcodes-b'
  [[module.imports]]
    path = 'github.com/user/shortcodes-c'
{{< /code-toggle >}}

Import precedence is top-down. For example, if `shortcodes-a`, `shortcodes-b`, and `shortcodes-c` each define an `image` shortcode, the `image` shortcode from `shortcodes-a` will take effect.

> [!note]
> If multiple modules contain data files or [translation tables](g) with identical paths, the data is deeply merged, following top-down precedence.

When you build your project, Hugo will:

1. Download the modules
1. Cache them for future use
1. Generate a [`go.sum`][] file in the project root

See [configuring module imports][] for details and options.

## Update

When you import a module, Hugo creates `go.mod` and `go.sum` files in your project root, storing version and checksum data. Clearing the module cache and rebuilding will re-download the originally imported module version, as specified in the `go.mod` file, ensuring consistent builds. Modules can be updated to other versions as needed.

To update a module to the latest version:

```sh
hugo mod get -u github.com/user/shortcodes-c
```

To update a module to a specific version:

```sh
hugo mod get -u github.com/user/shortcodes-c@v0.42.0
```

To update all modules to the latest version:

```sh
hugo mod get -u
```

To recursively update all modules to the latest version:

```sh
hugo mod get -u ./...
```

## Tidy

To remove unused entries from the `go.mod` and `go.sum` files:

```sh
hugo mod tidy
```

## Cache

Hugo caches modules to avoid repeated downloads during site builds. By default, these are stored in the `modules` directory within the [`cacheDir`][].

To clean the module cache for the current project:

```sh
hugo mod clean
```

To clean the module cache for all projects:

```sh
hugo mod clean --all
```

For details on cache location and eviction, see [configuring file caches][].

## Vendor

{{% glossary-term vendor %}}

Vendoring a module provides the benefits described above and allows for local inspection of its [components](g).

```sh
hugo mod vendor
```

This command creates a `_vendor` directory containing copies of all imported modules, used in subsequent builds. Note that:

- The `hugo mod vendor` command can be run from any module tree level.
- Modules within the `themes` directory are not vendored.
- The `--ignoreVendorPaths` flag allows you to exclude vendored modules matching a [glob pattern](g) from specific commands.

> [!important]
> Instead of modifying files directly within the `_vendor` directory, override them by creating a corresponding file with the same relative path in your project's root.

To remove the vendored modules, delete the `_vendor` directory.

## Replace

For local module development, use a `replace` directive in `go.mod` pointing to your local directory:

```text
replace github.com/user/module => /home/user/projects/module
```

With `hugo serve`r running, this change will trigger a configuration reload and add the local directory to the watch list. Alternatively, configure replacements by setting the [`replacements`][] parameter in your project configuration.

## Workspace

{{% glossary-term "workspace" %}}

Workspaces simplify local development of sites with modules. Create a `.work` file to define a workspace, and activate it via the [`workspace`][] configuration parameter or the `HUGO_MODULE_WORKSPACE` environment variable.

A `.work` file example:

```text
go 1.24

use .
use ../my-hugo-module
```

Use the `use` directive to list module paths, including the main project (`.`). Start the Hugo server with the workspace enabled:

```sh
HUGO_MODULE_WORKSPACE=hugo.work hugo server --ignoreVendorPaths "**"
```

The `--ignoreVendorPaths` flag, used to ignore vendored dependencies (if applicable), enables live reloading of local edits within the workspace.

## Graph

To generate a [dependency graph](g), including vendoring, module replacement, and disabled module information, execute `hugo mod graph` within the target module directory. For example:

```sh
$ hugo mod graph

github.com/bep/my-modular-site github.com/bep/hugotestmods/mymounts@v1.2.0
github.com/bep/my-modular-site github.com/bep/hugotestmods/mypartials@v1.0.7
github.com/bep/hugotestmods/mypartials@v1.0.7 github.com/bep/hugotestmods/myassets@v1.0.4
github.com/bep/hugotestmods/mypartials@v1.0.7 github.com/bep/hugotestmods/myv2@v1.0.0
DISABLED github.com/bep/my-modular-site github.com/spf13/hyde@v0.0.0-20190427180251-e36f5799b396
github.com/bep/my-modular-site github.com/bep/hugo-fresh@v1.0.1
github.com/bep/my-modular-site in-themesdir
```

## Mounts

Imported modules automatically mount their component directories to Hugo's [unified file system](g). You can also manually mount any directory, including those from non-Hugo projects, to component directories.

See [configuring module mounts][] for details.

[`cacheDir`]: /configuration/all/#cachedir
[`go.mod`]: https://go.dev/ref/mod#go-mod-file
[`go.sum`]: https://go.dev/ref/mod#go-sum-files
[`replacements`]: /configuration/module/#replacements
[`workspace`]: /configuration/module/#workspace
[configuring file caches]: /configuration/caches/
[configuring module imports]: /configuration/module/#imports
[configuring module mounts]: /configuration/module/#mounts
[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[merge configuration settings]: /configuration/introduction/#merge-configuration-settings
