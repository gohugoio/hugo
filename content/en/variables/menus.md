---
title: Menu Variables
description: Use these variables and methods in your menu templates.
categories: [variables and params]
keywords: [menus]
menu:
  docs:
    parent: variables
    weight: 50
weight: 50
aliases: [/variables/menu/]
---

## Variables

After [defining the menu entries], access their properties with these variables.

.Children
: (`menu`) A collection of child menu entries, if any, under the current menu entry.

.Identifier
: (`string`) The `identifier` property of the menu entry. If you define the menu entry [automatically], the page's `.Section`.

.KeyName
: (`string`) The `identifier` property of the menu entry, else the `name` property.

.Menu
: (`string`) The identifier of the menu that contains the menu entry.

.Name
: (`string`) The `name` property of the menu entry.

- If you define the menu entry [automatically], the page's `.LinkTitle`, else the page's `.Title`.
- If you define the menu [in front matter] or [in site configuration], falls back to the page's `.LinkTitle`, then to the page's `.Title`.

.Page
: (`page`) A reference to the page associated with the menu entry.

<!-- This provides no value when rendering menu. Omitting to avoid confusion.
.PageRef
: (`string`) The `pageRef` property of the menu entry.
-->

.Params
: (`map`) The `params` property of the menu entry.

.Parent
: (`string`)  The `parent` property of the menu entry.

.Post
: (`template.HTML`) The `post` property of the menu entry.

.Pre
: (`template.HTML`) The `pre` property of the menu entry.

.Title
: (`string`) The `title` property of the menu entry.

- If you define the menu entry [automatically], the page's `.LinkTitle`, else the page's `.Title`.
- If you define the menu [in front matter] or [in site configuration], falls back to the page's `.LinkTitle`, then to the page's `.Title`.

.URL
: (`string`) The `.RelPermalink` of the page associated with the menu entry. For menu entries pointing to external resources, the `url` property of the menu entry.

.Weight
: (`int`) The `weight` property of the menu entry.

- If you define the menu entry [automatically], the page's `.Weight`.
- If you define the menu [in front matter] or [in site configuration], falls back to the page's `.Weight`.

## Methods

.HasChildren
: (`bool`) Returns `true` if `.Children` is non-nil.

.IsEqual
: (`bool`) Returns `true` if the compared menu entries represent the same menu entry.

.IsSameResource
: (`bool`) Returns `true` if the compared menu entries point to the same resource.

.Page.HasMenuCurrent
: (`bool`) Use this method to determine ancestors of the active menu entry. See [details](/functions/hasmenucurrent/).

.Page.IsMenuCurrent
: (`bool`) Use this method to determine the active menu entry. See [details](/functions/ismenucurrent/).

[automatically]: /content-management/menus/#define-automatically
[defining the menu entries]: /content-management/menus/#overview
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration
