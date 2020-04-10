---
title: Use Hugo Modules
linktitle: Use Hugo Modules
description: How to use Hugo Modules to build and manage your site.
date: 2019-07-24
categories: [hugo modules]
keywords: [install, themes, source, organization, directories,usage,modules]
menu:
  docs:
    parent: "modules"
    weight: 20
weight: 20
sections_weight: 20
draft: false
aliases: [/themes/usage/,/themes/installing/,/installing-and-using-themes/]
toc: true
---

## Prerequisite

{{< gomodules-info >}}



## Initialize a New Module

Use `hugo mod init` to initialize a new Hugo Module. If it fails to guess the module path, you must provide it as an argument, e.g.:

```bash
hugo mod init github.com/gohugoio/myShortcodes
```


Also see the [CLI Doc](/commands/hugo_mod_init/).

## Update Modules

Modules will be downloaded and added when you add them as imports to your configuration, see [Module Imports](/hugo-modules/configuration/#module-config-imports).

To update or manage versions, you can use `hugo mod get`.

Some examples:

### Update All Modules

```bash
hugo mod get -u
```

### Update All Modules Recursively

{{< new-in "0.65.0" >}}

```bash
hugo mod get -u ./...
```

### Update One Module

```bash
hugo mod get -u github.com/gohugoio/myShortcodes
```
### Get a Specific Version

```bash
hugo mod get github.com/gohugoio/myShortcodes@v1.0.7
```

Also see the [CLI Doc](/commands/hugo_mod_get/).

## Make and test changes in a module

One way to do local development of a module imported in a project is to add a replace directive to a local directory with the source in `go.mod`:

```bash
replace github.com/bep/hugotestmods/mypartials => /Users/bep/hugotestmods/mypartials
```

If you have the `hugo server` running, the configuration will be reloaded and `/Users/bep/hugotestmods/mypartials` put on the watch list. 


## Print Dependency Graph


Use `hugo mod graph` from the relevant module directory and it will print the dependency graph, including vendoring, module replacement or disabled status.

E.g.:

```
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

## Vendor Your Modules

`hugo mod vendor` will write all the module depencies to a `_vendor` folder, which will then be used for all subsequent builds.

Note that:

* You can run `hugo mod vendor` on any level in the module tree.
* Vendoring will not store modules stored in your `themes` folder.
* Most commands accept a `--ignoreVendor` flag, which will then run as if the none of the `_vendor` folders in the module tree existed.

Also see the [CLI Doc](/commands/hugo_mod_vendor/).


## Tidy go.mod, go.sum

Run `hugo mod tidy` to remove unused entries in `go.mod` and `go.sum`.

Also see the [CLI Doc](/commands/hugo_mod_clean/).

## Clean Module Cache

Run `hugo mod clean` to delete the entire modules cache.

Note that you can also configure the `modules` cache with a `maxAge`, see [File Caches](/hugo-modules/configuration/#configure-file-caches).



Also see the [CLI Doc](/commands/hugo_mod_clean/).
