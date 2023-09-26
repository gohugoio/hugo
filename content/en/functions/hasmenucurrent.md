---
title: .HasMenuCurrent
description:
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace:
relatedFuncs:
  - .HasMenuCurrent
  - .IsMenuCurrent
signature:
  - PAGE.HasMenuCurrent MENU MENUENTRY
---

`.HasMenuCurrent` is a method in `Page` object returning a _boolean_ value. It
returns `true` if the PAGE is the same object as the `.Page` in one of the
**children menu entries** under MENUENTRY in a given MENU.

If MENUENTRY's `.Page` is a [section](/content-management/sections/) then, from Hugo `0.86.0`, this method also returns true for any descendant of that section..

You can find its example use in [menu templates](/templates/menu-templates/).
