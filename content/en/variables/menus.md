---
title: Menu Entry Properties
linktitle: Menu Entry Properties
description: A menu entry in a menu-template has specific variables and functions to make menu management easier.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
categories: [variables and params]
keywords: [menus]
draft: false
menu:
  docs:
    title: "variables defined by a menu entry"
    parent: "variables"
    weight: 50
weight: 50
sections_weight: 50
aliases: [/variables/menu/]
toc: false
---

A **menu entry** has the following properties available that can be used in a
[menu template][menu-template].

## Menu Entry Variables

.Menu
: _string_ <br />
Name of the **menu** that contains this **menu entry**.

.URL
: _string_ <br />
URL that the menu entry points to. The `url` key, if set for the menu entry,
sets this value. If that key is not set, and if the menu entry is set in a page
front-matter, this value defaults to the page's `.RelPermalink`.

.Page
: _\*Page_ <br />
Reference to the [page object][page-object] associated with the menu entry. This
will be non-nil if the menu entry is set via a page's front-matter and not via
the site config.

.Name
: _string_ <br />
Name of the menu entry. The `name` key, if set for the menu entry, sets
this value. If that key is not set, and if the menu entry is set in a page
front-matter, this value defaults to the page's `.LinkTitle`.

.Identifier
: _string_ <br />
Value of the `identifier` key if set for the menu entry. This value must be
unique for each menu entry. **It is necessary to set a unique identifier
manually if two or more menu entries have the same `.Name`.**

.Pre
: _template.HTML_ <br />
Value of the `pre` key if set for the menu entry. This value typically contains
a string representing HTML.

.Post
: _template.HTML_ <br />
Value of the `post` key if set for the menu entry. This value typically contains
a string representing HTML.

.Weight
: _int_ <br />
Value of the `weight` key if set for the menu entry. If that key is not set,
and if the menu entry is set in a page front-matter, this value defaults to the
page's `.Weight`.

.Parent
: _string_ <br />
Name (or Identifier if present) of this menu entry's parent **menu entry**. The
`parent` key, if set for the menu entry, sets this value. If this key is set,
this menu entry nests under that parent entry, else it nests directly under the
`.Menu`.

.Children
: _Menu_ <br />
This value is auto-populated by Hugo. It is a collection of children menu
entries, if any, under the current menu entry.

## Menu Entry Functions

Menus also have the following functions available:

.HasChildren
: _boolean_ <br />
Returns `true` if `.Children` is non-nil.

.KeyName
: _string_ <br />
Returns the `.Identifier` if present, else returns the `.Name`.

.IsEqual
: _boolean_ <br />
Returns `true` if the two compared menu entries represent the same menu entry.

.IsSameResource
: _boolean_ <br />
Returns `true` if the two compared menu entries have the same `.URL`.

.Title
: _string_ <br />
Link title, meant to be used in the `title` attribute of a menu entry's
`<a>`-tags.  Returns the menu entry's `title` key if set. Else, if the menu
entry was created through a page's front-matter, it returns the page's
`.LinkTitle`. Else, it just returns an empty string.

## Other Menu-related Functions

Additionally, here are some relevant methods available to menus on a page:

.IsMenuCurrent
: _(menu string, menuEntry *MenuEntry ) boolean_ <br />
See [`.IsMenuCurrent` method](/functions/ismenucurrent/).

.HasMenuCurrent
: _(menu string, menuEntry *MenuEntry) boolean_ <br />
See [`.HasMenuCurrent` method](/functions/hasmenucurrent/).


[menu-template]: /templates/menu-templates/
[page-object]: /variables/page/
