---
title: Go Plugins
linktitle: plugins
description: Hugo can use go plugin on rendering.
date: 2021-05-07
publishdate: 2021-05-07
lastmod: 2021-05-07
categories: [templates]
keywords: [plugin,module]
menu:
  docs:
    parent: "templates"
    weight: 180
weight: 180
sections_weight: 180
draft: false
aliases: []
toc: false
---

If your current project needs additional templating functions, you can use go plugin.

## Write your own plugin

Install [go](https://golang.org/doc/install) then follow [documentation to write and compile your own plugin](https://pkg.go.dev/plugin).

Once compiled, put the `*.so` file in `plugins/` folder.

## Use it in content files

### Variables

Variables in go plugin can be fetched thanks to the [`plugin.Get` function](../functions/plugin.Get).

### Functions

Functions can be called by using the [`plugin.Call` functions](../functions/plugin.Call).
