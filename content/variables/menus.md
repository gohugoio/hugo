---
title: Menu Variables
linktitle: Menu Variables
description: A menu entry in a menu template has specific variables and functions to make menu management easier.
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

The [menu template][] has the following properties:

.URL
: string

.Name
: string

.Title
: string

This is a link title, meant to be used in `title`-Attributes of the menu's `<a>`-tags.
By default it returns `.Page.LinkTitle`, as long as the menu entry was created
through the page's front matter and not through the site config.
Setting it explicitly in the site config or the page's front matter overrides this behaviour.

.Page
: [Page Object](/variables/page/)

The `.Page` variable holds a reference to the page.
It's only set when the menu entry is created from the page's front matter,
not when it's created from the site config.


.Menu
: string

.Identifier
: string

.Pre
: template.HTML

.Post
: template.HTML

.Weight
: int

.Parent
: string

.Children
: Menu

[menu template]: /templates/menu-templates/
