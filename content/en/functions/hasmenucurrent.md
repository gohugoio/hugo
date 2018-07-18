---
title: .HasMenuCurrent
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [menus]
signature: ["PAGE.HasMenuCurrent MENU MENUENTRY"]
workson: [menus]
hugoversion:
relatedfuncs: ["IsMenuCurrent"]
deprecated: false
toc: false
draft: false
aliases: []
---

`.HasMenuCurrent` is a method in `Page` object returning a _boolean_ value. It
returns `true` if the PAGE is the same object as the `.Page` in one of the
**children menu entries** under MENUENTRY in a given MENU.

You can find its example use in [menu templates](/templates/menu-templates/).
