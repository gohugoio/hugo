---
title: Menu Variables
linktitle: Menu Variables
description: A menu entry in a menu template has specific variables and functions to make menu management easier.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
categories: [variables and params]
tags: [menus]
draft: false
menu:
  main:
    parent: "Variables"
    weight: 50
weight: 50
sections_weight: 50
aliases: [/variables/menu/]
toc: false
---

A menu entry in a [menu template][] has the following properties:

`.URL`
: string

`.Name`
: string

`.Menu`
: string

`.Identifier`
: string

`.Pre`
: template.HTML

`.Post`
: template.HTML

`.Weight`
: int

`.Parent`
: string

`.Children`
: Menu

[menu template]: /templates/menu-templates/