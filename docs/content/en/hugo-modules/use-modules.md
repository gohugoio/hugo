---
title: Use Hugo Modules
description: How to use Hugo Modules to build and manage your site.
categories: [hugo modules]
keywords: [modules,themes]
menu:
  docs:
    parent: modules
    weight: 30
weight: 30
aliases: [/themes/usage/,/themes/installing/,/installing-and-using-themes/]
toc: true
---

## Prerequisite

{{< gomodules-info >}}

## Initialize a new module

Use `hugo mod init` to initialize a new Hugo Module. If it fails to guess the module path, you must provide it as an argument, e.g.:

```sh
hugo mod init github.com/<your_user>/<your_project>
```

Also see the [CLI Doc](/commands/hugo_mod_init/).

## Use a module for a theme

The easiest way to use a Module for a theme is to import it in the configuration.

1. Initialize the hugo module system: `hugo mod init github.com/<your_user>/<your_project>`
2. Import the theme:

{{< code-toggle file=hugo >}}
[module]
  [[module.imports]]
    path = "github.com/spf13/hyde"
{{< /code-toggle >}}

## Update modules

Modules will be downloaded and added when you add them as imports to your configuration, see [Module Imports](/hugo-modules/configuration/#module-configuration-imports).

To update or manage versions, you can use `hugo mod get`.

Some examples:

### Update all modules

```sh
hugo mod get -u
```

### Update all modules recursively

```sh
hugo mod get -u ./...
```

### Update one module

```sh
hugo mod get -u github.com/gohugoio/myShortcodes
```

### Get a specific version

```sh
hugo mod get github.com/gohugoio/myShortcodes@v1.0.7
```

Also see the [CLI Doc](/commands/hugo_mod_get/).

## Make and test changes in a module

One way to do local development of a module imported in a project is to add a replace directive to a local directory with the source in `go.mod`:

```sh
replace github.com/bep/hugotestmods/mypartials => /Users/bep/hugotestmods/mypartials
```

If you have the `hugo server` running, the configuration will be reloaded and `/Users/bep/hugotestmods/mypartials` put on the watch list.

Instead of modifying the `go.mod` files, you can also use the modules configuration [`replacements`](/hugo-modules/configuration/#module-configuration-top-level) option.

## Print dependency graph

Use `hugo mod graph` from the relevant module directory and it will print the dependency graph, including vendoring, module replacement or disabled status.

E.g.:

```txt
hugo mod graph

github.com/bep/my-modular-site github.com/bep/hugotestmods/mymounts@v1.2.0
github.com/bep/my-modular-site github.com/bep/hugotestmods/mypartials@v1.0.7
github.com/bep/hugotestmods/mypartials@v1.0.7 github.com/bep/hugotestmods/myassets@v1.0.4
github.com/bep/hugotestmods/mypartials@v1.0.7 github.com/bep/hugotestmods/myv2@v1.0.0
DISABLED github.com/bep/my-modular-site github.com/spf13/hyde@v0.0.0-20190427180251-e36f5799b396
github.com/bep/my-modular-site github.com/bep/hugo-fresh@v1.0.1
github.com/bep/my-modular-site in-themesdir
```

Also see the [CLI Doc](/commands/hugo_mod_graph/).

## Vendor your modules

`hugo mod vendor` will write all the module dependencies to a `_vendor` folder, which will then be used for all subsequent builds.

Note that:

* You can run `hugo mod vendor` on any level in the module tree.
* Vendoring will not store modules stored in your `themes` folder.
* Most commands accept a `--ignoreVendorPaths` flag, which will then not use the vendored modules in `_vendor` for the module paths matching the [Glob](https://github.com/gobwas/glob) pattern given.

Also see the [CLI Doc](/commands/hugo_mod_vendor/).

## Tidy go.mod, go.sum

Run `hugo mod tidy` to remove unused entries in `go.mod` and `go.sum`.

Also see the [CLI Doc](/commands/hugo_mod_clean/).

## Clean module cache

Run `hugo mod clean` to delete the entire modules cache.

Note that you can also configure the `modules` cache with a `maxAge`, see [File Caches](/getting-started/configuration/#configure-file-caches).

Also see the [CLI Doc](/commands/hugo_mod_clean/).

## Module workspaces

{{< new-in 0.109.0 >}}

Workspace support was added in [Go 1.18](https://go.dev/blog/get-familiar-with-workspaces) and Hugo got solid support for it in the `v0.109.0` version.

A common use case for a workspace is to simplify local development of a site with its theme modules.

A workspace can be configured in a `*.work` file and activated with the [module.workspace](/hugo-modules/configuration/) setting, which for this use is commonly controlled via the `HUGO_MODULE_WORKSPACE` OS environment variable.

See the [hugo.work](https://github.com/gohugoio/hugo/blob/master/docs/hugo.work) file in the Hugo Docs repo for an example:

```text
go 1.20

use .
use ../gohugoioTheme
```

Using the `use` directive, list all the modules you want to work on, pointing to its relative location. As in the example above, it's recommended to always include the main project (the ".") in the list.

With that you can start the Hugo server with that workspace enabled:

```sh
HUGO_MODULE_WORKSPACE=hugo.work hugo server --ignoreVendorPaths "**"
```

The `--ignoreVendorPaths` flag is added above to ignore any of the vendored dependencies inside `_vendor`. If you don't use vendoring, you don't need that flag. But now the server is set up watching the files and directories in the workspace and you can see your local edits reloaded.
