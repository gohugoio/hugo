---
title: Menus
linktitle: Menus
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templating]
tags: [lists,sections,menus]
draft: false
weight: 120
aliases: [/extras/menus/]
toc: true
needsreview: true
---

Hugo has a simple yet powerful menu system that permits content to be placed in menus with a good degree of control without a lot of work.

{{% note "Lazy Blogger"%}}
If all you want is a simple menu for your sections, see the ["Section Menu for Lazy Bloggers" in Menu Templates](/templates/menu-templates/#section-menu-for-lazy-blogger).
{{% /note %}}

Some of the features of Hugo Menus:

* Place content in one or many menus
* Handle nested menus with unlimited depth
* Create menu entries without being attached to any content
* Distinguish active element (and active branch)

## What is a Menu?

A menu is a named array of menu entries accessible on the site under `.Site.Menus` by name. For example, if I have a menu called `main`, I would access it via `.Site.Menus.main`.

If you make use of the [multilingual feature](content-management/multilingual-mode/) you can define language-independent menus.

A menu entry has the following properties:

`URL`
: string

`Name`
: string

`Menu`
: string

`Identifier`
: string

`Pre`
: template.HTML

`Post`
: template.HTML

`Weight`
: int

`Parent`
: string

`Children`
: Menu

And the following functions:

`HasChildren`
: boolean

Additionally, there are some relevant functions available on the page:

`IsMenuCurrent`
: (menu string, menuEntry *MenuEntry ) boolean

`HasMenuCurrent`
: (menu string, menuEntry *MenuEntry) bool

## Adding content to menus

Hugo supports a couple of different methods of adding a piece of content to the front matter.

### Simple

If all you need to do is add an entry to a menu, the simple form works well.

#### A Single Menu:

```yaml
---
menu: "main"
---
```

#### Multiple Menus:

```yaml
---
menu: ["main", "footer"]
---
```

### Advanced

Take the advanced approach if more control is required. All of the menu entry properties listed above are available.

```yaml
---
menu:
  main:
    parent: 'extras'
    weight: 20
---
```

## Adding Non-content Entries to a Menu

You can also add entries to menus that aren’t attached to a piece of content. This takes place in the sitewide [config file](/overview/configuration/).

Here’s an example `config.toml`:

```toml
[[menu.main]]
    name = "about hugo"
    pre = "<i class='fa fa-heart'></i>"
    weight = -110
    identifier = "about"
    url = "/about/"
[[menu.main]]
    name = "getting started"
    pre = "<i class='fa fa-road'></i>"
    weight = -100
    url = "/getting-started/"
```

And the equivalent example `config.yaml`:

```yaml
---
menu:
  main:
      - Name: "about hugo"
        Pre: "<i class='fa fa-heart'></i>"
        Weight: -110
        Identifier: "about"
        URL: "/about/"
      - Name: "getting started"
        Pre: "<i class='fa fa-road'></i>"
        Weight: -100
        URL: "/getting-started/"
---
```


**NOTE:** The URLs must be relative to the context root. If the `baseURL` is `http://example.com/mysite/`, then the URLs in the menu must not include the context root `mysite`. Using an absolute URL will overide the baseURL. If the `URL` is `http://subdomain.example.com/`, the output will be `http://subdomain.example.com`.

## Nesting

All nesting of content is done via the `parent` field.

The parent of an entry should be the identifier of another entry.
Identifier should be unique (within a menu).

The following order is used to determine an Identifier:

* Name >
    * LinkTitle >
        * Title

This means that `title` will be used unless `linktitle` is present, etc. In practice, Name and Identifier are never displayed and only used to structure relationships.

In this example, the top level of the menu is defined in the config file
and all content entries are attached to one of these entries via the
`parent` field.

## Rendering Menus

See [Menu Templates](/templates/menu-templates/) for information on how to render your site menus.
