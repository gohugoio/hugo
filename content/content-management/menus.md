---
title: Menus
linktitle: Menus
description: Hugo's simple yet powerful menu system allows you to add content to menus with a high level of control and low level of effort.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templating]
tags: [lists,sections,menus]
draft: false
weight: 120
aliases: [/extras/menus/]
toc: true
---

Hugo's simple yet powerful menu system allows you to add content to menus with a high degree of control and a very low amount of effort on the part of content authors.

{{% note "Lazy Blogger"%}}
If all you want is a simple menu for your sections, see the ["Section Menu for Lazy Bloggers" in Menu Templates](/templates/menu-templates/#section-menu-for-lazy-blogger).
{{% /note %}}

Hugo's robust menu features allow you to do the following:

* Place content in one or many menus
* Handle nested menus with unlimited depth
* Create menu entries without being attached to any content
* Distinguish active element (and active branch)

## What is a Menu in Hugo?

A **menu** is a named array of menu entries accessible by name via the [`.Site.Menus` site variable][sitevars]. For example, you can access your site's `main` menu via `.Site.Menus.main`.

{{% note "Menus on Multilingual Sites" %}}
If you make use of the [multilingual feature](/content-management/multilingual/), you can define language-independent menus.
{{% /note %}}

A menu entry has the following properties (i.e., variables) available to it:

{{< readfile file="content/readfiles/menuvars.md" markdown="true" >}}

## Adding content to menus

Hugo allows you to add content to a menu via the content's [front matter](/content-management/front-matter/).

### Simple

If all you need to do is add an entry to a menu, the simple form works well.

#### A Single Menu

```yaml
---
menu: "main"
---
```

#### Multiple Menus

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

You can also add entries to menus that aren’t attached to a piece of content. This takes place in your Hugo project's [`config` file][config].

Here’s an example snippet pulled from a `config.toml`:

{{% code file="config.toml" %}}
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
{{% /code %}}

Here's the equivalent snippet in a `config.yaml`:

{{% code file="config.yml" %}}
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
{{% /code %}}

{{% note %}}
The URLs must be relative to the context root. If the `baseURL` is `http://example.com/mysite/`, then the URLs in the menu must not include the context root `mysite`. Using an absolute URL will overide the baseURL. If the value used for `URL` in the above example is `http://subdomain.example.com/`, the output will be `http://subdomain.example.com`.
{{% /note %}}

## Nesting

All nesting of content is done via the `parent` field.

The parent of an entry should be the identifier of another entry. The identifier should be unique (within a menu).

The following order is used to determine an Identifier:

`.Name > .LinkTitle > .Title`

This means that `.Title` will be used unless `.LinkTitle` is present, etc. In practice, `.Name` and `.Identifier` are only used to structure relationships and therefore never displayed.

In this example, the top level of the menu is defined in your [site `config` file][config]). All content entries are attached to one of these entries via the `.Parent` field.

## Rendering Menus

See [Menu Templates](/templates/menu-templates/) for information on how to render your site menus within your templates.

[config]: /getting-started/configuration/
[multilingual]: /content-management/multilingual/
[sitevars]: /variables/
