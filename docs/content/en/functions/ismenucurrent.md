---
title: .IsMenuCurrent
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
signature: ["PAGE.IsMenuCurrent MENU MENUENTRY"]
workson: [menus]
hugoversion:
relatedfuncs: ["HasMenuCurrent"]
deprecated: false
draft: false
aliases: []
needsexample: true
---

`.IsMenuCurrent` is a method in `Page` object returning a _boolean_ value. It
returns `true` if the PAGE is the same object as the `.Page` in MENUENTRY in a
given MENU.

You can find its example use in [menu templates](/templates/menu-templates/).
