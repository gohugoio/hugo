---
title: Menus
linktitle: Menus
description: Hugo has a simple yet powerful menu system.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-31
categories: [content management]
keywords: [menus]
draft: false
menu:
  docs:
    parent: "content-management"
    weight: 120
weight: 120	#rem
aliases: [/extras/menus/]
toc: true
---

{{% note "Lazy Blogger"%}}
If all you want is a simple menu for your sections, see the ["Section Menu for Lazy Bloggers" in Menu Templates](/templates/menu-templates/#section-menu-for-lazy-bloggers).
{{% /note %}}

You can do this:

* Place content in one or many menus
* Handle nested menus with unlimited depth
* Create menu entries without being attached to any content
* Distinguish active element (and active branch)

## What is a Menu in Hugo?

A **menu** is a named array of menu entries accessible by name via the [`.Site.Menus` site variable][sitevars]. For example, you can access your site's `main` menu via `.Site.Menus.main`.

{{% note "Menus on Multilingual Sites" %}}
If you make use of the [multilingual feature](/content-management/multilingual/), you can define language-independent menus.
{{% /note %}}

See the [Menu Entry Properties][me-props] for all the variables and functions related to a menu entry.

## Add content to menus

Hugo allows you to add content to a menu via the content's [front matter](/content-management/front-matter/).

### Simple

If all you need to do is add an entry to a menu, the simple form works well.

#### A Single Menu

```
---
menu: "main"
---
```

#### Multiple Menus

```
---
menu: ["main", "footer"]
---
```

#### Advanced


```
---
menu:
  docs:
    parent: 'extras'
    weight: 20
---
```

## Add Non-content Entries to a Menu

You can also add entries to menus that aren’t attached to a piece of content. This takes place in your Hugo project's [`config` file][config].

Here’s an example snippet pulled from a configuration file:

{{< code-toggle file="config" >}}
[[menu.main]]
    name = "about hugo"
    pre = "<i class='fa fa-heart'></i>"
    weight = -110
    identifier = "about"
    url = "/about/"
[[menu.main]]
    name = "getting started"
    pre = "<i class='fa fa-road'></i>"
    post = "<span class='alert'>New!</span>"
    weight = -100
    url = "/getting-started/"
{{< /code-toggle >}}

{{% note %}}
The URLs must be relative to the context root. If the `baseURL` is `https://example.com/mysite/`, then the URLs in the menu must not include the context root `mysite`. Using an absolute URL will override the baseURL. If the value used for `URL` in the above example is `https://subdomain.example.com/`, the output will be `https://subdomain.example.com`.
{{% /note %}}

## Nesting

All nesting of content is done via the `parent` field.

The parent of an entry should be the identifier of another entry. The identifier should be unique (within a menu).

The following order is used to determine an Identifier:

`.Name > .LinkTitle > .Title`

This means that `.Title` will be used unless `.LinkTitle` is present, etc. In practice, `.Name` and `.Identifier` are only used to structure relationships and therefore never displayed.

In this example, the top level of the menu is defined in your [site `config` file][config]. All content entries are attached to one of these entries via the `.Parent` field.

## Render Menus

See [Menu Templates](/templates/menu-templates/) for information on how to render your site menus within your templates.

[config]: /getting-started/configuration/
[multilingual]: /content-management/multilingual/
[sitevars]: /variables/
[me-props]: /variables/menus/
